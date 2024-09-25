package cbor

import (
	"github.com/x448/float16"
	"io"
	"math"
)

// readMajorType reads the major type and any header arguments from [in].
func readMajorType(in io.Reader) (MajorType, Arg, uint64, error) {
	b := sharedBuffer[:1]
	n, err := in.Read(b)
	if n != 1 {
		return 0, 0, 0, err
	}

	majorType, arg := decodePrefix(b[0])
	arg, l, v, err := decodeArg(arg)
	if err != nil {
		return 0, 0, 0, err
	}

	if l > 0 {
		b = sharedBuffer[1 : 1+l]
		n, err = in.Read(b)
		if n != int(l) {
			return 0, 0, 0, err
		}
		v = shiftBytesInto[uint64](b)
	}

	return majorType, arg, v, nil
}

// decodePrefix returns the major type, argument, and the length, in bytes, of
// the remaining part of the header if any.
func decodePrefix(p byte) (MajorType, Arg) {
	majorType := MajorType(p & majorTypeMask)
	arg := Arg(p & argMask)
	return majorType, arg
}

func decodeArg(arg Arg) (Arg, uint8, uint64, error) {
	switch {
	case arg < Arg8:
		return 0, 0, uint64(arg), nil
	case arg == Arg8:
		return arg, 1, 0, nil
	case arg == Arg16:
		return arg, 2, 0, nil
	case arg == Arg32:
		return arg, 4, 0, nil
	case arg == Arg64:
		return arg, 8, 0, nil
	case arg == ArgIndefinite:
		return arg, 0, 0, nil
	default:
		return 0, 0, 0, ErrNotWellFormed
	}
}

// ReadTag reads the next object from [in] as a tag and returns the value.
func ReadTag(in io.Reader) (uint64, error) {
	majorType, _, value, err := readMajorType(in)
	if err != nil {
		return 0, err
	}

	if majorType != MajorTypeTagged {
		return 0, ErrUnsupportedMajorType
	}

	return value, nil
}

func ReadUnsigned[T uint8 | uint16 | uint32 | uint64](in io.Reader) (T, error) {
	majorType, arg, value, err := readMajorType(in)
	if err != nil {
		return 0, err
	}

	return readUnsigned[T](majorType, arg, value)
}

func readUnsigned[T uint8 | uint16 | uint32 | uint64](majorType MajorType, arg Arg, value uint64) (T, error) {
	if majorType == MajorTypeUInt ||
		(majorType == MajorTypeSimpleFloat && arg == 0 && value < uint64(SimpleFalse)) ||
		(majorType == MajorTypeSimpleFloat && arg == SimpleUint8) {

		switch any(T(0)).(type) {
		case uint8:
			if (value & 0xffffffffffffff00) != 0 {
				return 0, ErrOverflow
			}
		case uint16:
			if (value & 0xffffffffffff0000) != 0 {
				return 0, ErrOverflow
			}
		case uint32:
			if (value & 0xffffffff00000000) != 0 {
				return 0, ErrOverflow
			}
		case uint64:
			// noop
		default:
			panic("unreachable")
		}

		return T(value), nil
	}

	return 0, ErrUnsupportedMajorType
}

func ReadSigned[T int8 | int16 | int32 | int64](in io.Reader) (T, error) {
	majorType, arg, value, err := readMajorType(in)
	if err != nil {
		return 0, err
	}

	return readSigned[T](majorType, arg, value)
}

func readSigned[T int8 | int16 | int32 | int64](majorType MajorType, arg Arg, value uint64) (T, error) {
	if majorType == MajorTypeUInt ||
		majorType == MajorTypeNInt ||
		(majorType == MajorTypeSimpleFloat && arg == 0 && value < uint64(SimpleFalse)) ||
		(majorType == MajorTypeSimpleFloat && arg == SimpleUint8) {

		switch any(T(0)).(type) {
		case int8:
			if (value & 0xffffffffffffff80) != 0 {
				return 0, ErrOverflow
			}
		case int16:
			if (value & 0xffffffffffff8000) != 0 {
				return 0, ErrOverflow
			}
		case int32:
			if (value & 0xffffffff80000000) != 0 {
				return 0, ErrOverflow
			}
		case int64:
			if (value & 0x8000000000000000) != 0 {
				return 0, ErrOverflow
			}
		default:
			panic("unreachable")
		}

		if majorType == MajorTypeNInt {
			return -1 - T(value), nil
		}

		return T(value), nil
	}

	return 0, ErrUnsupportedMajorType
}

func ReadFloat[T float32 | float64](in io.Reader) (T, error) {
	majorType, arg, value, err := readMajorType(in)
	if err != nil {
		return 0, err
	}

	return readFloat[T](majorType, arg, value)
}

func readFloat[T float32 | float64](majorType MajorType, arg Arg, value uint64) (T, error) {
	if majorType == MajorTypeUInt ||
		(majorType == MajorTypeSimpleFloat && arg == SimpleUint8) {
		return T(value), nil
	}

	if majorType == MajorTypeNInt {
		return -1 - T(value), nil
	}

	if majorType == MajorTypeSimpleFloat {
		if arg == SimpleFloat16 {
			return T(float16.Frombits(uint16(value)).Float32()), nil
		}

		if arg == SimpleFloat32 {
			return T(math.Float32frombits(uint32(value))), nil
		}

		if arg == SimpleFloat64 {
			v := math.Float64frombits(value)
			switch any(T(0)).(type) {
			case float32:
				v32 := T(v)
				if v != float64(v32) {
					return 0, ErrOverflow
				}
				return v32, nil
			case float64:
				return T(v), nil
			default:
				panic("unreachable")
			}
		}

		return 0, ErrUnsupportedValue
	}

	return 0, ErrUnsupportedMajorType
}

func ReadBool(in io.Reader) (bool, error) {
	majorType, _, value, err := readMajorType(in)
	if err != nil {
		return false, err
	}

	return readBool(majorType, value)
}

func readBool(majorType MajorType, value uint64) (bool, error) {
	if majorType == MajorTypeSimpleFloat {
		if byte(value) == SimpleFalse {
			return false, nil
		}

		if byte(value) == SimpleTrue {
			return true, nil
		}

		return false, ErrUnsupportedValue
	}

	return false, ErrUnsupportedMajorType
}

func ReadBytes(
	in io.Reader,
	readLength func(indefinite bool, length uint64) error,
	out io.Writer,
) error {
	majorType, arg, value, err := readMajorType(in)
	if err != nil {
		return err
	}

	return readBytes(in, majorType, arg, value, readLength, out)
}

func readBytes(
	in io.Reader,
	majorType MajorType,
	arg Arg,
	value uint64,
	readLength func(indefinite bool, length uint64) error,
	out io.Writer,
) error {
	if majorType != MajorTypeBstr && majorType != MajorTypeTstr {
		return ErrUnsupportedMajorType
	}

	indefinite := arg == ArgIndefinite

	err := readLength(indefinite, value)
	if err != nil {
		return err
	}

	if indefinite {
		for {
			majorType, arg, value, err := readMajorType(in)
			if err != nil {
				return err
			}

			if majorType == MajorTypeSimpleFloat && arg == SimpleBreak {
				break
			}

			if majorType != MajorTypeBstr && majorType != MajorTypeTstr {
				return ErrUnsupportedMajorType
			}

			if arg == ArgIndefinite {
				return ErrNestedIndefinite
			}

			err = readByteChunks(in, value, out)
			if err != nil {
				return err
			}
		}
	} else {
		err = readByteChunks(in, value, out)
		if err != nil {
			return err
		}
	}

	return nil
}

func readByteChunks(
	in io.Reader,
	length uint64,
	out io.Writer,
) error {
	var l int
	var b []byte
	var n int
	var err error
	for length > 0 {
		l = int(min(lenSharedBuffer, length))
		b = sharedBuffer[:l]
		n, err = in.Read(b)
		if n != l {
			return err
		}
		n, err = out.Write(b[:n])
		if n != l {
			return err
		}

		length -= uint64(n)
	}

	return nil
}

func ReadArray(
	in io.Reader,
	readLength func(indefinite bool, length uint64) error,
	readItem func(in io.Reader) error,
) error {
	majorType, arg, value, err := readMajorType(in)
	if err != nil {
		return err
	}

	return readArray(in, majorType, arg, value, readLength, readItem)
}

func readArray(
	in io.Reader,
	majorType MajorType,
	arg Arg,
	value uint64,
	readLength func(indefinite bool, length uint64) error,
	readItem func(in io.Reader) error,
) error {
	if majorType != MajorTypeArray {
		return ErrUnsupportedMajorType
	}

	indefinite := arg == ArgIndefinite

	err := readLength(indefinite, value)
	if err != nil {
		return err
	}

	if indefinite {
		pin := &peekReader{r: in}
		for {
			r, err := pin.PeekByte()
			if err != nil {
				return err
			}
			if r == valueBreak {
				break
			}

			err = readItem(pin)
			if err != nil {
				return err
			}
		}
	} else {
		for i := uint64(0); i < value; i++ {
			err = readItem(in)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ReadMap(
	in io.Reader,
	readLength func(indefinite bool, length uint64) error,
	readKeyValue func(in io.Reader) error,
) error {
	majorType, arg, value, err := readMajorType(in)
	if err != nil {
		return err
	}

	return readMap(in, majorType, arg, value, readLength, readKeyValue)
}

func readMap(
	in io.Reader,
	majorType MajorType,
	arg Arg,
	value uint64,
	readLength func(indefinite bool, length uint64) error,
	readKeyValue func(in io.Reader) error,
) error {
	if majorType != MajorTypeMap {
		return ErrUnsupportedMajorType
	}

	indefinite := arg == ArgIndefinite

	err := readLength(indefinite, value)
	if err != nil {
		return err
	}

	if indefinite {
		pin := &peekReader{r: in}
		for {
			r, err := pin.PeekByte()
			if err != nil {
				return err
			}
			if r == valueBreak {
				break
			}

			err = readKeyValue(pin)
			if err != nil {
				return err
			}
		}
	} else {
		for i := uint64(0); i < value; i++ {
			err = readKeyValue(in)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
