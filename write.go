package cbor

import (
	"github.com/x448/float16"
	"io"
	"math"
)

func writeMajorType(out io.Writer, majorType MajorType, value uint64) (int, error) {
	sharedBuffer[0] = byte(majorType)

	if value < uint64(Arg8) {
		sharedBuffer[0] |= byte(value)
		return out.Write(sharedBuffer[0:1])
	}

	n := 0
	switch {
	case value < 0x1_00:
		sharedBuffer[0] |= byte(Arg8)
		n = 1
	case value < 0x1_00_00:
		sharedBuffer[0] |= byte(Arg16)
		n = 2
	case value < 0x1_00_00_00_00:
		sharedBuffer[0] |= byte(Arg32)
		n = 4
	default:
		sharedBuffer[0] |= byte(Arg64)
		n = 8
	}

	shiftBytesFrom(value, sharedBuffer[1:1+n])
	return out.Write(sharedBuffer[0 : 1+n])
}

func WriteUnsigned[T uint8 | uint16 | uint32 | uint64](out io.Writer, value T) (int, error) {
	return writeMajorType(out, MajorTypeUInt, uint64(value))
}

func WriteSigned[T int8 | int16 | int32 | int64](out io.Writer, value T) (int, error) {
	if value >= 0 {
		return writeMajorType(out, MajorTypeUInt, uint64(value))
	}
	return writeMajorType(out, MajorTypeNInt, uint64(-value-1))
}

func WriteFloat[T float16.Float16 | float32 | float64](out io.Writer, value T) (int, error) {
	switch v := any(value).(type) {
	case float16.Float16:
		sharedBuffer[0] = MajorTypeSimpleFloat | SimpleFloat16
		shiftBytesFrom(uint16(v), sharedBuffer[1:3])
		return out.Write(sharedBuffer[0:3])

	case float32:
		if float16.PrecisionFromfloat32(v) == float16.PrecisionExact {
			return WriteFloat(out, float16.Fromfloat32(v))
		}

		sharedBuffer[0] = MajorTypeSimpleFloat | SimpleFloat32
		shiftBytesFrom(math.Float32bits(v), sharedBuffer[1:5])
		return out.Write(sharedBuffer[0:5])

	case float64:
		v32 := float32(v)
		// TODO NaN, inf...
		if v == float64(v32) {
			return WriteFloat(out, v32)
		}

		sharedBuffer[0] = MajorTypeSimpleFloat | SimpleFloat64
		shiftBytesFrom(math.Float64bits(v), sharedBuffer[1:9])
		return out.Write(sharedBuffer[0:9])

	default:
		panic("unreachable")
	}
}

func WriteBool(out io.Writer, value bool) (int, error) {
	if value {
		return writeMajorType(out, MajorTypeSimpleFloat, SimpleTrue)
	}
	return writeMajorType(out, MajorTypeSimpleFloat, uint64(SimpleFalse))
}

func WriteTag(out io.Writer, value uint64) (int, error) {
	return writeMajorType(out, MajorTypeTagged, value)
}

func WriteBytes(out io.Writer, value []byte) (int, error) {
	tn := 0
	n, err := writeMajorType(out, MajorTypeBstr, uint64(len(value)))
	tn += n
	if err != nil {
		return tn, err
	}
	n, err = out.Write(value)
	tn += n
	return tn, err
}

func WriteString(out io.Writer, value string) (int, error) {
	tn := 0
	n, err := writeMajorType(out, MajorTypeTstr, uint64(len(value)))
	tn += n
	if err != nil {
		return tn, err
	}
	n, err = out.Write(([]byte)(value))
	tn += n
	return tn, err
}

func WriteArrayHeader(out io.Writer, length uint64) (int, error) {
	return writeMajorType(out, MajorTypeArray, length)
}

func WriteMapHeader(out io.Writer, length uint64) (int, error) {
	return writeMajorType(out, MajorTypeMap, length)
}
