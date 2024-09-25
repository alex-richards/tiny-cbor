//go:build !cbor_no_readany

package cbor

import (
	"bytes"
	"io"
)

// ReadAny returns the next object form [in] regardless of type.
// Outputs can be any of int64, uint64, bool, []byte, string, []any,
// map[any]any, float32, float64, nil.
func ReadAny(in io.Reader) (any, error) {
	majorType, arg, value, err := readMajorType(in)
	if err != nil {
		return 0, err
	}

	return readAny(in, majorType, arg, value)
}

func readAny(in io.Reader, majorType MajorType, arg Arg, value uint64) (any, error) {
	var err error
	switch majorType {
	case MajorTypeUInt:
		switch {
		case arg <= Arg8:
			return readUnsigned[uint8](majorType, arg, value)
		case arg == Arg16:
			return readUnsigned[uint16](majorType, arg, value)
		case arg == Arg32:
			return readUnsigned[uint32](majorType, arg, value)
		case arg == Arg64:
			return readUnsigned[uint64](majorType, arg, value)
		default:
			return nil, ErrNotWellFormed
		}

	case MajorTypeNInt:
		switch {
		case arg <= Arg8:
			return readSigned[int8](majorType, arg, value)
		case arg == Arg16:
			return readSigned[int16](majorType, arg, value)
		case arg == Arg32:
			return readSigned[int32](majorType, arg, value)
		case arg == Arg64:
			return readSigned[int64](majorType, arg, value)
		default:
			return nil, ErrNotWellFormed
		}

	case MajorTypeBstr:
		b := bytes.NewBuffer(nil)
		err = readBytes(in, majorType, arg, value,
			func(indefinite bool, length uint64) error {
				b.Grow(int(length))
				return nil
			},
			b,
		)
		if err != nil {
			return nil, err
		}
		return b.Bytes(), nil

	case MajorTypeTstr:
		b := bytes.NewBuffer(nil)
		err = readBytes(in, majorType, arg, value,
			func(indefinite bool, length uint64) error {
				b.Grow(int(length))
				return nil
			},
			b,
		)
		if err != nil {
			return nil, err
		}
		return string(b.Bytes()), nil

	case MajorTypeArray:
		a := make([]any, value)
		if arg == ArgIndefinite {
			for {
				majorType, arg, value, err := readMajorType(in)
				if err != nil {
					return nil, err
				}

				if majorType == MajorTypeSimpleFloat && arg == SimpleBreak {
					break
				}

				v, err := readAny(in, majorType, arg, value)
				if err != nil {
					return nil, err
				}
				a = append(a, v)
			}
		} else {
			for i := uint64(0); i < value; i++ {
				v, err := ReadAny(in)
				if err != nil {
					return nil, err
				}
				a[i] = v
			}
		}
		return a, nil

	case MajorTypeMap:
		m := make(map[any]any, value)
		if arg == ArgIndefinite {
			for {
				majorType, arg, value, err := readMajorType(in)
				if err != nil {
					return nil, err
				}

				if majorType == MajorTypeSimpleFloat && arg == SimpleBreak {
					break
				}

				k, err := readAny(in, majorType, arg, value)
				if err != nil {
					return nil, err
				}
				v, err := ReadAny(in)
				if err != nil {
					return nil, err
				}
				m[k] = v
			}
		} else {
			for i := uint64(0); i < value; i++ {
				k, err := ReadAny(in)
				if err != nil {
					return nil, err
				}
				v, err := ReadAny(in)
				if err != nil {
					return nil, err
				}
				m[k] = v
			}
		}
		return m, nil

	case MajorTypeTagged:
		return ReadAny(in)

	default: // MajorTypeSimpleFloat:
		switch {
		case arg == 0 && value < uint64(SimpleFalse):
			return readUnsigned[uint8](majorType, arg, value)
		case arg == 0 && (value == uint64(SimpleFalse) || value == SimpleTrue):
			return readBool(majorType, value)
		case arg == 0 && (value == SimpleNull || value == SimpleUndefined):
			return nil, nil
		case arg == SimpleUint8:
			return readUnsigned[uint8](majorType, arg, value)
		case arg == SimpleFloat16:
			return readFloat[float32](majorType, arg, value)
		case arg == SimpleFloat32:
			return readFloat[float32](majorType, arg, value)
		case arg == SimpleFloat64:
			return readFloat[float64](majorType, arg, value)
		default: // SimpleBreak
			return nil, ErrNotWellFormed
		}
	}
}
