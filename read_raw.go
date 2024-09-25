package cbor

import (
	"io"
)

func ReadRaw(
	in io.Reader,
	out io.Writer,
) error {
	n, err := in.Read(sharedBuffer[:1])
	if n != 1 {
		return err
	}

	majorType, arg := decodePrefix(sharedBuffer[0])
	arg, l, v, err := decodeArg(arg)
	if err != nil {
		return err
	}

	ve := 1 + l
	if l > 0 {
		n, err = in.Read(sharedBuffer[1:ve])
		if n != int(l) {
			return err
		}
		v = shiftBytesInto[uint64](sharedBuffer[1:ve])
	}

	n, err = out.Write(sharedBuffer[0:ve])
	if n != int(ve) {
		return err
	}

	switch majorType {
	case MajorTypeUInt, MajorTypeNInt, MajorTypeSimpleFloat:
		return nil

	case MajorTypeBstr, MajorTypeTstr:
		if arg == ArgIndefinite {
			pin := &peekReader{r: in}
			for {
				b, err := pin.PeekByte()
				if err != nil {
					return err
				}

				if b == valueBreak {
					out.Write([]byte{b})
					break
				}

				err = ReadRaw(pin, out)
				if err != nil {
					return err
				}
			}
			return nil
		} else {
			return readBytes(in, majorType, arg, v,
				func(indefinite bool, length uint64) error { return nil },
				out,
			)
		}

	case MajorTypeArray:
		if arg == ArgIndefinite {
			pin := &peekReader{r: in}
			for {
				b, err := pin.PeekByte()
				if err != nil {
					return err
				}

				if b == valueBreak {
					out.Write([]byte{b})
					break
				}

				err = ReadRaw(pin, out)
				if err != nil {
					return err
				}
			}
			return nil
		} else {
			for range v {
				err = ReadRaw(in, out)
				if err != nil {
					return err
				}
			}
			return nil
		}

	case MajorTypeMap:
		if arg == ArgIndefinite {
			pin := &peekReader{r: in}
			for {
				b, err := pin.PeekByte()
				if err != nil {
					return err
				}

				if b == valueBreak {
					out.Write([]byte{b})
					break
				}

				for range 2 {
					err = ReadRaw(pin, out)
					if err != nil {
						return err
					}
				}
			}
			return nil
		} else {
			for range v {
				for range 2 {
					err = ReadRaw(in, out)
					if err != nil {
						return err
					}
				}
			}
			return nil
		}

	default: // MajorTypeTagged
		return ReadRaw(in, out)
	}
}
