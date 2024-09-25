package cbor

import (
	"bytes"
	"encoding/hex"
	"io"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/x448/float16"
)

func Test_ReadRaw(t *testing.T) {
	for _, tt := range tests_ExampleEncoded {
		t.Run(tt.encoded, func(t *testing.T) {
			encoded := decodeHex(t, tt.encoded)

			in := bytes.NewReader(encoded)
			out := bytes.NewBuffer(make([]byte, 0, 128))

			err := ReadRaw(in, out)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(encoded, out.Bytes()); diff != "" {
				t.Fatal(diff)
			}

			n := in.Len()
			if n != 0 {
				rem := make([]byte, n)
				_, _ = in.Read(rem)
				t.Fatalf("trailing data - %s", hex.EncodeToString(rem))
			}
		})
	}
}

func Benchmark_ReadRaw(b *testing.B) {
	for _, tt := range tests_ExampleEncoded {
		encoded := decodeHex(b, tt.encoded)

		in := bytes.NewReader(encoded)
		out := bytes.NewBuffer(make([]byte, 0, 128))

		b.Run(tt.encoded, func(b *testing.B) {
			for range b.N {
				_, _ = in.Seek(io.SeekStart, 0)
				out.Reset()
				err := ReadRaw(in, out)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func Test_ReadOver(t *testing.T) {
	for _, tt := range tests_ExampleEncoded {
		t.Run(tt.encoded, func(t *testing.T) {
			encoded := decodeHex(t, tt.encoded)

			in := bytes.NewReader(encoded)

			err := ReadOver(in)
			if err != nil {
				t.Fatal(err)
			}

			n := in.Len()
			if n != 0 {
				rem := make([]byte, n)
				_, _ = in.Read(rem)
				t.Fatalf("trailing data - %s", hex.EncodeToString(rem))
			}
		})
	}
}

func Benchmark_ReadOver(b *testing.B) {
	for _, tt := range tests_ExampleEncoded {
		encoded := decodeHex(b, tt.encoded)

		in := bytes.NewReader(encoded)

		b.Run(tt.encoded, func(b *testing.B) {
			for range b.N {
				_, _ = in.Seek(0, io.SeekStart)
				_ = ReadOver(in)
			}
		})
	}
}

func Test_ReadNumbers(t *testing.T) {
	tests := []struct {
		encoded string

		wantInt8       int8
		wantInt8Error  error
		wantInt16      int16
		wantInt16Error error
		wantInt32      int32
		wantInt32Error error
		wantInt64      int64
		wantInt64Error error

		wantUint8       uint8
		wantUint8Error  error
		wantUint16      uint16
		wantUint16Error error
		wantUint32      uint32
		wantUint32Error error
		wantUint64      uint64
		wantUint64Error error

		wantFloat32      float32
		wantFloat32Error error
		wantFloat64      float64
		wantFloat64Error error
	}{
		{
			encoded:     "00",
			wantInt8:    0x00,
			wantInt16:   0x00,
			wantInt32:   0x00,
			wantInt64:   0x00,
			wantUint8:   0x00,
			wantUint16:  0x00,
			wantUint32:  0x00,
			wantUint64:  0x00,
			wantFloat32: 0x00,
			wantFloat64: 0x00,
		},
		{
			encoded:     "01",
			wantInt8:    0x01,
			wantInt16:   0x01,
			wantInt32:   0x01,
			wantInt64:   0x01,
			wantUint8:   0x01,
			wantUint16:  0x01,
			wantUint32:  0x01,
			wantUint64:  0x01,
			wantFloat32: 0x01,
			wantFloat64: 0x01,
		},
		{
			encoded:     "17",
			wantInt8:    0x17,
			wantInt16:   0x17,
			wantInt32:   0x17,
			wantInt64:   0x17,
			wantUint8:   0x17,
			wantUint16:  0x17,
			wantUint32:  0x17,
			wantUint64:  0x17,
			wantFloat32: 0x17,
			wantFloat64: 0x17,
		},
		{
			encoded:     "1818",
			wantInt8:    0x18,
			wantInt16:   0x18,
			wantInt32:   0x18,
			wantInt64:   0x18,
			wantUint8:   0x18,
			wantUint16:  0x18,
			wantUint32:  0x18,
			wantUint64:  0x18,
			wantFloat32: 0x18,
			wantFloat64: 0x18,
		},
		{
			encoded:       "18FF",
			wantInt8Error: ErrOverflow,
			wantInt16:     0xFF,
			wantInt32:     0xFF,
			wantInt64:     0xFF,
			wantUint8:     0xFF,
			wantUint16:    0xFF,
			wantUint32:    0xFF,
			wantUint64:    0xFF,
			wantFloat32:   0xFF,
			wantFloat64:   0xFF,
		},
		{
			encoded:        "1901FF",
			wantInt8Error:  ErrOverflow,
			wantInt16:      0x01FF,
			wantInt32:      0x01FF,
			wantInt64:      0x01FF,
			wantUint8Error: ErrOverflow,
			wantUint16:     0x01FF,
			wantUint32:     0x01FF,
			wantUint64:     0x01FF,
			wantFloat32:    0x01FF,
			wantFloat64:    0x01FF,
		},
		{
			encoded:        "19FFFF",
			wantInt8Error:  ErrOverflow,
			wantInt16Error: ErrOverflow,
			wantInt32:      0xFFFF,
			wantInt64:      0xFFFF,
			wantUint8Error: ErrOverflow,
			wantUint16:     0xFFFF,
			wantUint32:     0xFFFF,
			wantUint64:     0xFFFF,
			wantFloat32:    0xFFFF,
			wantFloat64:    0xFFFF,
		},
		{
			encoded:         "1A0001FFFF",
			wantInt8Error:   ErrOverflow,
			wantInt16Error:  ErrOverflow,
			wantInt32:       0x01FFFF,
			wantInt64:       0x01FFFF,
			wantUint8Error:  ErrOverflow,
			wantUint16Error: ErrOverflow,
			wantUint32:      0x01FFFF,
			wantUint64:      0x01FFFF,
			wantFloat32:     0x01FFFF,
			wantFloat64:     0x01FFFF,
		},
		{
			encoded:         "1AFFFFFFFF",
			wantInt8Error:   ErrOverflow,
			wantInt16Error:  ErrOverflow,
			wantInt32Error:  ErrOverflow,
			wantInt64:       0xFFFFFFFF,
			wantUint8Error:  ErrOverflow,
			wantUint16Error: ErrOverflow,
			wantUint32:      0xFFFFFFFF,
			wantUint64:      0xFFFFFFFF,
			wantFloat32:     0xFFFFFFFF,
			wantFloat64:     0xFFFFFFFF,
		},
		{
			encoded:         "1B00000001FFFFFFFF",
			wantInt8Error:   ErrOverflow,
			wantInt16Error:  ErrOverflow,
			wantInt32Error:  ErrOverflow,
			wantInt64:       0x01FFFFFFFF,
			wantUint8Error:  ErrOverflow,
			wantUint16Error: ErrOverflow,
			wantUint32Error: ErrOverflow,
			wantUint64:      0x01FFFFFFFF,
			wantFloat32:     0x01FFFFFFFF,
			wantFloat64:     0x01FFFFFFFF,
		},
		{
			encoded:         "1BFFFFFFFFFFFFFFFF",
			wantInt8Error:   ErrOverflow,
			wantInt16Error:  ErrOverflow,
			wantInt32Error:  ErrOverflow,
			wantInt64Error:  ErrOverflow,
			wantUint8Error:  ErrOverflow,
			wantUint16Error: ErrOverflow,
			wantUint32Error: ErrOverflow,
			wantUint64:      0xFFFFFFFFFFFFFFFF,
			wantFloat32:     0xFFFFFFFFFFFFFFFF,
			wantFloat64:     0xFFFFFFFFFFFFFFFF,
		},
		// TODO negatives
		{
			encoded:     "F800",
			wantInt8:    0x00,
			wantInt16:   0x00,
			wantInt32:   0x00,
			wantInt64:   0x00,
			wantUint8:   0x00,
			wantUint16:  0x00,
			wantUint32:  0x00,
			wantUint64:  0x00,
			wantFloat32: 0x00,
			wantFloat64: 0x00,
		},
		{
			encoded:       "F8FF",
			wantInt8Error: ErrOverflow,
			wantInt16:     0xFF,
			wantInt32:     0xFF,
			wantInt64:     0xFF,
			wantUint8:     0xFF,
			wantUint16:    0xFF,
			wantUint32:    0xFF,
			wantUint64:    0xFF,
			wantFloat32:   0xFF,
			wantFloat64:   0xFF,
		},
		{
			encoded:         "F90000",
			wantInt8Error:   ErrUnsupportedMajorType,
			wantInt16Error:  ErrUnsupportedMajorType,
			wantInt32Error:  ErrUnsupportedMajorType,
			wantInt64Error:  ErrUnsupportedMajorType,
			wantUint8Error:  ErrUnsupportedMajorType,
			wantUint16Error: ErrUnsupportedMajorType,
			wantUint32Error: ErrUnsupportedMajorType,
			wantUint64Error: ErrUnsupportedMajorType,
			wantFloat32:     0,
			wantFloat64:     0,
		},
		{
			encoded:         "F97BFF", // max
			wantInt8Error:   ErrUnsupportedMajorType,
			wantInt16Error:  ErrUnsupportedMajorType,
			wantInt32Error:  ErrUnsupportedMajorType,
			wantInt64Error:  ErrUnsupportedMajorType,
			wantUint8Error:  ErrUnsupportedMajorType,
			wantUint16Error: ErrUnsupportedMajorType,
			wantUint32Error: ErrUnsupportedMajorType,
			wantUint64Error: ErrUnsupportedMajorType,
			wantFloat32:     float16.Frombits(0x7bFF).Float32(),
			wantFloat64:     float64(float16.Frombits(0x7bFF).Float32()),
		},
		{
			encoded:         "F97C00", // +ve inf
			wantInt8Error:   ErrUnsupportedMajorType,
			wantInt16Error:  ErrUnsupportedMajorType,
			wantInt32Error:  ErrUnsupportedMajorType,
			wantInt64Error:  ErrUnsupportedMajorType,
			wantUint8Error:  ErrUnsupportedMajorType,
			wantUint16Error: ErrUnsupportedMajorType,
			wantUint32Error: ErrUnsupportedMajorType,
			wantUint64Error: ErrUnsupportedMajorType,
			wantFloat32:     float32(math.Inf(1)),
			wantFloat64:     math.Inf(1),
		},
		{
			encoded:         "F9FC00", // -ve inf
			wantInt8Error:   ErrUnsupportedMajorType,
			wantInt16Error:  ErrUnsupportedMajorType,
			wantInt32Error:  ErrUnsupportedMajorType,
			wantInt64Error:  ErrUnsupportedMajorType,
			wantUint8Error:  ErrUnsupportedMajorType,
			wantUint16Error: ErrUnsupportedMajorType,
			wantUint32Error: ErrUnsupportedMajorType,
			wantUint64Error: ErrUnsupportedMajorType,
			wantFloat32:     float32(math.Inf(-1)),
			wantFloat64:     math.Inf(-1),
		},
		{
			encoded:         "F90000", // 0
			wantInt8Error:   ErrUnsupportedMajorType,
			wantInt16Error:  ErrUnsupportedMajorType,
			wantInt32Error:  ErrUnsupportedMajorType,
			wantInt64Error:  ErrUnsupportedMajorType,
			wantUint8Error:  ErrUnsupportedMajorType,
			wantUint16Error: ErrUnsupportedMajorType,
			wantUint32Error: ErrUnsupportedMajorType,
			wantUint64Error: ErrUnsupportedMajorType,
			wantFloat32:     0,
			wantFloat64:     0,
		},
		{
			encoded:         "F98000", // -0
			wantInt8Error:   ErrUnsupportedMajorType,
			wantInt16Error:  ErrUnsupportedMajorType,
			wantInt32Error:  ErrUnsupportedMajorType,
			wantInt64Error:  ErrUnsupportedMajorType,
			wantUint8Error:  ErrUnsupportedMajorType,
			wantUint16Error: ErrUnsupportedMajorType,
			wantUint32Error: ErrUnsupportedMajorType,
			wantUint64Error: ErrUnsupportedMajorType,
			wantFloat32:     -0,
			wantFloat64:     -0,
		},
		{
			encoded:         "FA7F7FFFFF", // max
			wantInt8Error:   ErrUnsupportedMajorType,
			wantInt16Error:  ErrUnsupportedMajorType,
			wantInt32Error:  ErrUnsupportedMajorType,
			wantInt64Error:  ErrUnsupportedMajorType,
			wantUint8Error:  ErrUnsupportedMajorType,
			wantUint16Error: ErrUnsupportedMajorType,
			wantUint32Error: ErrUnsupportedMajorType,
			wantUint64Error: ErrUnsupportedMajorType,
			wantFloat32:     math.Float32frombits(0x7F7FFFFF),
			wantFloat64:     float64(math.Float32frombits(0x7F7FFFFF)),
		},
		{
			encoded:         "FA7F800000", // +ve inf
			wantInt8Error:   ErrUnsupportedMajorType,
			wantInt16Error:  ErrUnsupportedMajorType,
			wantInt32Error:  ErrUnsupportedMajorType,
			wantInt64Error:  ErrUnsupportedMajorType,
			wantUint8Error:  ErrUnsupportedMajorType,
			wantUint16Error: ErrUnsupportedMajorType,
			wantUint32Error: ErrUnsupportedMajorType,
			wantUint64Error: ErrUnsupportedMajorType,
			wantFloat32:     float32(math.Inf(1)),
			wantFloat64:     math.Inf(1),
		},
		{
			encoded:         "FAFF800000", // -ve inf
			wantInt8Error:   ErrUnsupportedMajorType,
			wantInt16Error:  ErrUnsupportedMajorType,
			wantInt32Error:  ErrUnsupportedMajorType,
			wantInt64Error:  ErrUnsupportedMajorType,
			wantUint8Error:  ErrUnsupportedMajorType,
			wantUint16Error: ErrUnsupportedMajorType,
			wantUint32Error: ErrUnsupportedMajorType,
			wantUint64Error: ErrUnsupportedMajorType,
			wantFloat32:     float32(math.Inf(-1)),
			wantFloat64:     math.Inf(-1),
		},
		{
			encoded:         "FA00000000", // 0
			wantInt8Error:   ErrUnsupportedMajorType,
			wantInt16Error:  ErrUnsupportedMajorType,
			wantInt32Error:  ErrUnsupportedMajorType,
			wantInt64Error:  ErrUnsupportedMajorType,
			wantUint8Error:  ErrUnsupportedMajorType,
			wantUint16Error: ErrUnsupportedMajorType,
			wantUint32Error: ErrUnsupportedMajorType,
			wantUint64Error: ErrUnsupportedMajorType,
			wantFloat32:     0,
			wantFloat64:     0,
		},
		{
			encoded:         "FA80000000", // -0
			wantInt8Error:   ErrUnsupportedMajorType,
			wantInt16Error:  ErrUnsupportedMajorType,
			wantInt32Error:  ErrUnsupportedMajorType,
			wantInt64Error:  ErrUnsupportedMajorType,
			wantUint8Error:  ErrUnsupportedMajorType,
			wantUint16Error: ErrUnsupportedMajorType,
			wantUint32Error: ErrUnsupportedMajorType,
			wantUint64Error: ErrUnsupportedMajorType,
			wantFloat32:     -0,
			wantFloat64:     -0,
		},
		{
			encoded:          "FB7FEFFFFFFFFFFFFF", // max
			wantInt8Error:    ErrUnsupportedMajorType,
			wantInt16Error:   ErrUnsupportedMajorType,
			wantInt32Error:   ErrUnsupportedMajorType,
			wantInt64Error:   ErrUnsupportedMajorType,
			wantUint8Error:   ErrUnsupportedMajorType,
			wantUint16Error:  ErrUnsupportedMajorType,
			wantUint32Error:  ErrUnsupportedMajorType,
			wantUint64Error:  ErrUnsupportedMajorType,
			wantFloat32Error: ErrOverflow,
			wantFloat64:      math.Float64frombits(0x7FEFFFFFFFFFFFFF),
		},
		{
			encoded:         "FB7FF0000000000000", // +ve inf
			wantInt8Error:   ErrUnsupportedMajorType,
			wantInt16Error:  ErrUnsupportedMajorType,
			wantInt32Error:  ErrUnsupportedMajorType,
			wantInt64Error:  ErrUnsupportedMajorType,
			wantUint8Error:  ErrUnsupportedMajorType,
			wantUint16Error: ErrUnsupportedMajorType,
			wantUint32Error: ErrUnsupportedMajorType,
			wantUint64Error: ErrUnsupportedMajorType,
			wantFloat32:     float32(math.Inf(1)),
			wantFloat64:     math.Inf(1),
		},
		{
			encoded:         "FBFFF0000000000000", // -ve inf
			wantInt8Error:   ErrUnsupportedMajorType,
			wantInt16Error:  ErrUnsupportedMajorType,
			wantInt32Error:  ErrUnsupportedMajorType,
			wantInt64Error:  ErrUnsupportedMajorType,
			wantUint8Error:  ErrUnsupportedMajorType,
			wantUint16Error: ErrUnsupportedMajorType,
			wantUint32Error: ErrUnsupportedMajorType,
			wantUint64Error: ErrUnsupportedMajorType,
			wantFloat32:     float32(math.Inf(-1)),
			wantFloat64:     math.Inf(-1),
		},
		{
			encoded:         "FB0000000000000000", // 0
			wantInt8Error:   ErrUnsupportedMajorType,
			wantInt16Error:  ErrUnsupportedMajorType,
			wantInt32Error:  ErrUnsupportedMajorType,
			wantInt64Error:  ErrUnsupportedMajorType,
			wantUint8Error:  ErrUnsupportedMajorType,
			wantUint16Error: ErrUnsupportedMajorType,
			wantUint32Error: ErrUnsupportedMajorType,
			wantUint64Error: ErrUnsupportedMajorType,
			wantFloat32:     0,
			wantFloat64:     0,
		},
		{
			encoded:         "FB8000000000000000", // -0
			wantInt8Error:   ErrUnsupportedMajorType,
			wantInt16Error:  ErrUnsupportedMajorType,
			wantInt32Error:  ErrUnsupportedMajorType,
			wantInt64Error:  ErrUnsupportedMajorType,
			wantUint8Error:  ErrUnsupportedMajorType,
			wantUint16Error: ErrUnsupportedMajorType,
			wantUint32Error: ErrUnsupportedMajorType,
			wantUint64Error: ErrUnsupportedMajorType,
			wantFloat32:     0,
			wantFloat64:     -0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.encoded, func(t *testing.T) {
			encoded := decodeHex(t, tt.encoded)

			runTest_ReadSigned(t, "int8", encoded, tt.wantInt8, tt.wantInt8Error)
			runTest_ReadSigned(t, "int16", encoded, tt.wantInt16, tt.wantInt16Error)
			runTest_ReadSigned(t, "int32", encoded, tt.wantInt32, tt.wantInt32Error)
			runTest_ReadSigned(t, "int32", encoded, tt.wantInt64, tt.wantInt64Error)

			runTest_ReadUnsigned(t, "uint8", encoded, tt.wantUint8, tt.wantUint8Error)
			runTest_ReadUnsigned(t, "uint16", encoded, tt.wantUint16, tt.wantUint16Error)
			runTest_ReadUnsigned(t, "uint32", encoded, tt.wantUint32, tt.wantUint32Error)
			runTest_ReadUnsigned(t, "uint64", encoded, tt.wantUint64, tt.wantUint64Error)

			runTest_ReadFloat(t, "float32", encoded, tt.wantFloat32, tt.wantFloat32Error)
			runTest_ReadFloat(t, "float64", encoded, tt.wantFloat64, tt.wantFloat64Error)
		})
	}
}

func runTest_ReadSigned[T int8 | int16 | int32 | int64](t *testing.T, name string, encoded []byte, want T, wantErr error) {
	t.Run(name, func(t *testing.T) {
		got, err := ReadSigned[T](bytes.NewReader(encoded))
		if err != wantErr {
			t.Fatalf("want %v, got %v", wantErr, err)
		}
		if err != nil {
			return
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatal(diff)
		}
	})
}

func runTest_ReadUnsigned[T uint8 | uint16 | uint32 | uint64](t *testing.T, name string, encoded []byte, want T, wantErr error) {
	t.Run(name, func(t *testing.T) {
		got, err := ReadUnsigned[T](bytes.NewReader(encoded))
		if err != wantErr {
			t.Fatalf("want %v, got %v", wantErr, err)
		}
		if err != nil {
			return
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatal(diff)
		}
	})
}

func runTest_ReadFloat[T float32 | float64](t *testing.T, name string, encoded []byte, want T, wantErr error) {
	t.Run(name, func(t *testing.T) {
		got, err := ReadFloat[T](bytes.NewReader(encoded))
		if err != wantErr {
			t.Fatalf("want %v, got %v", wantErr, err)
		}
		if err != nil {
			return
		}
		if math.IsNaN(float64(want)) && math.IsNaN(float64(got)) {
			// pass
		} else if diff := cmp.Diff(want, got); diff != "" {
			t.Fatal(diff)
		}
	})
}

func Test_ReadBool(t *testing.T) {
	for _, tt := range tests_ExampleEncoded {
		t.Run(tt.encoded, func(t *testing.T) {
			encoded := decodeHex(t, tt.encoded)

			in := bytes.NewReader(encoded)

			got, err := ReadBool(in)

			var want bool
			var wantErr error
			switch {
			case encoded[0] == 0xf4:
				want = false
			case encoded[0] == 0xf5:
				want = true
			case encoded[0]&majorTypeMask == MajorTypeSimpleFloat:
				wantErr = ErrUnsupportedValue
			default:
				wantErr = ErrUnsupportedMajorType
			}

			if got != want {
				t.Fatalf("want = %t, got = %t", want, got)
			}
			if err != wantErr {
				t.Fatalf("wantErr = %v, err = %v", wantErr, err)
			}
		})
	}
}

func Test_ReadArray(t *testing.T) {
	tests := []struct {
		encoded string
		want    []int32
	}{
		{
			encoded: "80",
			want:    []int32{},
		},
		{
			encoded: "83010203",
			want:    []int32{1, 2, 3},
		},
		{
			encoded: "98190102030405060708090a0b0c0d0e0f101112131415161718181819",
			want:    []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
		},
		{
			encoded: "9f0102030405060708090a0b0c0d0e0f101112131415161718181819ff",
			want:    []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
		},
	}

	for _, tt := range tests {
		t.Run(tt.encoded, func(t *testing.T) {
			in := bytes.NewBuffer(decodeHex(t, tt.encoded))
			var out []int32
			err := ReadArray(in,
				func(indefinite bool, length uint64) error {
					out = make([]int32, 0, length)
					return nil
				},
				func(in io.Reader) error {
					v, err := ReadSigned[int32](in)
					if err != nil {
						return err
					}
					out = append(out, v)
					return nil
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, out); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func Test_ReadMap(t *testing.T) {
	tests := []struct {
		encoded string
		want    map[int32]int32
	}{
		{
			encoded: "a0",
			want:    map[int32]int32{},
		},
		{
			encoded: "a201020304",
			want:    map[int32]int32{1: 2, 3: 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.encoded, func(t *testing.T) {
			in := bytes.NewBuffer(decodeHex(t, tt.encoded))
			var out map[int32]int32
			err := ReadMap(in,
				func(indefinite bool, length uint64) error {
					out = make(map[int32]int32, length)
					return nil
				},
				func(in io.Reader) error {
					k, err := ReadSigned[int32](in)
					if err != nil {
						return err
					}
					v, err := ReadSigned[int32](in)
					if err != nil {
						return err
					}
					out[k] = v
					return nil
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, out); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func Test_ReadBytes(t *testing.T) {
	tests := []struct {
		encoded string
		want    []byte
	}{
		{
			encoded: "40",
			want:    []byte{},
		},
		{
			encoded: "4401020304",
			want:    []byte{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.encoded, func(t *testing.T) {
			in := bytes.NewBuffer(decodeHex(t, tt.encoded))
			out := bytes.NewBuffer([]byte{})
			err := ReadBytes(
				in,
				func(indefinite bool, length uint64) error {
					out.Grow(int(length))
					return nil
				},
				out,
			)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, out.Bytes()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func Test_ReadTag(t *testing.T) {
	tests := []struct {
		encoded string
		want    uint64
	}{
		{
			encoded: "c0",
			want:    0,
		},
		{
			encoded: "c1",
			want:    1,
		},
		{
			encoded: "d7",
			want:    23,
		},
		{
			encoded: "d818",
			want:    24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.encoded, func(t *testing.T) {
			in := bytes.NewBuffer(decodeHex(t, tt.encoded))
			v, err := ReadTag(in)
			if err != nil {
				t.Fatal(v)
			}
			if diff := cmp.Diff(tt.want, v); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
