package cbor

import "io"

// ReadOver skips the next object in [in].
func ReadOver(in io.Reader) error {
	majorType, arg, value, err := readMajorType(in)
	if err != nil {
		return err
	}

	return readOver(in, majorType, arg, value)
}

// readOver skips the next object in [in], starting from after the header.
func readOver(in io.Reader, majorType MajorType, arg Arg, value uint64) error {
	switch majorType {
	case MajorTypeUInt,
		MajorTypeNInt,
		MajorTypeSimpleFloat:
		return nil

	case MajorTypeBstr,
		MajorTypeTstr:
		if arg == ArgIndefinite {
			for {
				majorType, arg, value, err := readMajorType(in)
				if err != nil {
					return err
				}

				if majorType == MajorTypeSimpleFloat && arg == SimpleBreak {
					break
				}

				for value > 0 {
					l := int(min(lenSharedBuffer, value))
					n, err := in.Read(sharedBuffer[:l])
					if n != l {
						return err
					}
					value -= uint64(n)
				}
			}
			return nil
		} else {
			for value > 0 {
				l := int(min(lenSharedBuffer, value))
				n, err := in.Read(sharedBuffer[0:l])
				if n != l {
					return err
				}
				value -= uint64(n)
			}
			return nil
		}

	case MajorTypeArray:
		if arg == ArgIndefinite {
			for {
				majorType, arg, value, err := readMajorType(in)
				if err != nil {
					return err
				}

				if majorType == MajorTypeSimpleFloat && arg == SimpleBreak {
					break
				}

				if err = readOver(in, majorType, arg, value); err != nil {
					return err
				}
			}
			return nil
		} else {
			for i := uint64(0); i < value; i++ {
				if err := ReadOver(in); err != nil {
					return err
				}
			}
			return nil
		}

	case MajorTypeMap:
		if arg == ArgIndefinite {
			for {
				majorType, arg, value, err := readMajorType(in)
				if err != nil {
					return err
				}

				if majorType == MajorTypeSimpleFloat && arg == SimpleBreak {
					break
				}

				if err = readOver(in, majorType, arg, value); err != nil {
					return err
				}
				if err = ReadOver(in); err != nil {
					return err
				}
			}
			return nil
		} else {
			for i := uint64(0); i < value; i++ {
				if err := ReadOver(in); err != nil {
					return err
				}
				if err := ReadOver(in); err != nil {
					return err
				}
			}
			return nil
		}

	case MajorTypeTagged:
		return ReadOver(in)

	default:
		return ErrUnsupportedMajorType
	}
}
