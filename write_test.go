package cbor

import (
	"bytes"
	"github.com/x448/float16"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_WriteUnsigned(t *testing.T) {
	tests := []struct {
		withUint8  uint8
		withUint16 uint16
		withUint32 uint32
		withUint64 uint64
		want       string
	}{
		{
			withUint8:  0,
			withUint16: 0,
			withUint32: 0,
			withUint64: 0,
			want:       "00",
		},
		{
			withUint8:  1,
			withUint16: 1,
			withUint32: 1,
			withUint64: 1,
			want:       "01",
		},
		{
			withUint8:  23,
			withUint16: 23,
			withUint32: 23,
			withUint64: 23,
			want:       "17",
		},
		{
			withUint8:  24,
			withUint16: 24,
			withUint32: 24,
			withUint64: 24,
			want:       "1818",
		},
		{
			withUint8:  0xff,
			withUint16: 0xff,
			withUint32: 0xff,
			withUint64: 0xff,
			want:       "18ff",
		},
		{
			withUint16: 0x100,
			withUint32: 0x100,
			withUint64: 0x100,
			want:       "190100",
		},
		{
			withUint16: 0x0102,
			withUint32: 0x0102,
			withUint64: 0x0102,
			want:       "190102",
		},
		{
			withUint16: 0xffff,
			withUint32: 0xffff,
			withUint64: 0xffff,
			want:       "19ffff",
		},
		{
			withUint32: 0x10000,
			withUint64: 0x10000,
			want:       "1a00010000",
		},
		{
			withUint32: 0x01020304,
			withUint64: 0x01020304,
			want:       "1a01020304",
		},
		{
			withUint32: 0xffffffff,
			withUint64: 0xffffffff,
			want:       "1affffffff",
		},
		{
			withUint64: 0x100000000,
			want:       "1b0000000100000000",
		},
		{
			withUint64: 0x0102030405060708,
			want:       "1b0102030405060708",
		},
		{
			withUint64: 0xffffffffffffffff,
			want:       "1bffffffffffffffff",
		},
	}

	for _, tt := range tests {
		want := decodeHex(t, tt.want)
		wantN := len(want)
		t.Run(tt.want, func(t *testing.T) {
			if wantN <= 2 {
				runTest_WriteUnsigned(t, "uint8", tt.withUint8, want, wantN)
			}
			if wantN <= 3 {
				runTest_WriteUnsigned(t, "uint16", tt.withUint16, want, wantN)
			}
			if wantN <= 5 {
				runTest_WriteUnsigned(t, "uint32", tt.withUint32, want, wantN)
			}
			if wantN <= 9 {
				runTest_WriteUnsigned(t, "uint64", tt.withUint64, want, wantN)
			}
		})
	}
}

func runTest_WriteUnsigned[N uint8 | uint16 | uint32 | uint64](t *testing.T, name string, with N, want []byte, wantN int) {
	t.Run(name, func(t *testing.T) {
		out := bytes.NewBuffer(nil)
		n, err := WriteUnsigned(out, with)
		if diff := cmp.Diff(wantN, n); diff != "" {
			t.Fatal(diff)
		}
		if diff := cmp.Diff(want, out.Bytes()); diff != "" {
			t.Fatal(diff)
		}
		if err != nil {
			t.Fatal(err)
		}
	})
}

func Test_WriteSigned(t *testing.T) {
	tests := []struct {
		with      string
		withInt8  int8
		withInt16 int16
		withInt32 int32
		withInt64 int64
		want      string
	}{
		{
			withInt8:  0,
			withInt16: 0,
			withInt32: 0,
			withInt64: 0,
			want:      "00",
		},
		{
			withInt8:  1,
			withInt16: 1,
			withInt32: 1,
			withInt64: 1,
			want:      "01",
		},
		{
			withInt8:  23,
			withInt16: 23,
			withInt32: 23,
			withInt64: 23,
			want:      "17",
		},
		{
			withInt8:  24,
			withInt16: 24,
			withInt32: 24,
			withInt64: 24,
			want:      "1818",
		},
		{
			withInt8:  -1,
			withInt16: -1,
			withInt32: -1,
			withInt64: -1,
			want:      "20",
		},
		{
			withInt8:  -24,
			withInt16: -24,
			withInt32: -24,
			withInt64: -24,
			want:      "37",
		},
		{
			withInt8:  -25,
			withInt16: -25,
			withInt32: -25,
			withInt64: -25,
			want:      "3818",
		},
		{
			withInt8:  127,
			withInt16: 127,
			withInt32: 127,
			withInt64: 127,
			want:      "187f",
		},
		{
			withInt8:  -128,
			withInt16: -128,
			withInt32: -128,
			withInt64: -128,
			want:      "387f",
		},
	}

	for _, tt := range tests {
		want := decodeHex(t, tt.want)
		t.Run(tt.want, func(t *testing.T) {
			runTest_WriteSigned(t, "int8", tt.withInt8, want)
			runTest_WriteSigned(t, "int16", tt.withInt16, want)
			runTest_WriteSigned(t, "int32", tt.withInt32, want)
			runTest_WriteSigned(t, "int64", tt.withInt64, want)
		})
	}
}

func runTest_WriteSigned[N int8 | int16 | int32 | int64](t *testing.T, name string, with N, want []byte) {
	t.Run(name, func(t *testing.T) {
		out := bytes.NewBuffer(nil)
		_, err := WriteSigned(out, with)
		if diff := cmp.Diff(want, out.Bytes()); diff != "" {
			t.Fatal(diff)
		}
		if err != nil {
			t.Fatal(err)
		}
	})
}

func Test_WriteFloat(t *testing.T) {
	tests := []struct {
		skipFloat16 bool
		withFloat16 float16.Float16
		skipFloat32 bool
		withFloat32 float32
		withFloat64 float64
		want        string
	}{
		{
			withFloat16: 0,
			withFloat32: 0,
			withFloat64: 0,
			want:        "f90000",
		},
		{
			withFloat16: float16.Fromfloat32(1.0),
			withFloat32: 1.0,
			withFloat64: 1.0,
			want:        "f93c00",
		},
		{
			skipFloat16: true,
			skipFloat32: true,
			withFloat64: 1.1,
			want:        "fb3ff199999999999a",
		},
		{
			withFloat16: float16.Fromfloat32(1.5),
			withFloat32: 1.5,
			withFloat64: 1.5,
			want:        "f93e00",
		},
		{
			withFloat16: float16.Fromfloat32(65504.0),
			withFloat32: 65504.0,
			withFloat64: 65504.0,
			want:        "f97bff",
		},
		{
			skipFloat16: true,
			withFloat32: 100000.0,
			withFloat64: 100000.0,
			want:        "fa47c35000",
		},
		{
			skipFloat16: true,
			withFloat32: 3.4028234663852886e+38,
			withFloat64: 3.4028234663852886e+38,
			want:        "fa7f7fffff",
		},
		{
			skipFloat16: true,
			skipFloat32: true,
			withFloat64: 1.0e+300,
			want:        "fb7e37e43c8800759c",
		},
		//{ TODO ??
		//	withFloat16: float16.Fromfloat32(5.960464477539063e-8),
		//	withFloat32: 5.960464477539063e-8,
		//	withFloat64: 5.960464477539063e-8,
		//	want:        "f90001",
		//},
		{
			withFloat16: float16.Fromfloat32(0.00006103515625),
			withFloat32: 0.00006103515625,
			withFloat64: 0.00006103515625,
			want:        "f90400",
		},
		{
			withFloat16: float16.Fromfloat32(-4.0),
			withFloat32: -4.0,
			withFloat64: -4.0,
			want:        "f9c400",
		},
		{
			skipFloat16: true,
			skipFloat32: true,
			withFloat64: -4.1,
			want:        "fbc010666666666666",
		},
	}

	for _, tt := range tests {
		want := decodeHex(t, tt.want)

		t.Run(tt.want, func(t *testing.T) {
			if !tt.skipFloat16 {
				runTest_WriteFloat(t, "float16", tt.withFloat16, want)
			}
			if !tt.skipFloat32 {
				runTest_WriteFloat(t, "float32", tt.withFloat32, want)
			}
			runTest_WriteFloat(t, "float64", tt.withFloat64, want)
		})
	}
}

func runTest_WriteFloat[N float16.Float16 | float32 | float64](t *testing.T, name string, with N, want []byte) {
	t.Run(name, func(t *testing.T) {
		out := bytes.NewBuffer(nil)
		_, err := WriteFloat(out, with)
		if diff := cmp.Diff(want, out.Bytes()); diff != "" {
			t.Fatal(diff)
		}
		if err != nil {
			t.Fatal(err)
		}
	})
}

func Test_WriteBool(t *testing.T) {
	tests := []struct {
		with bool
		want string
	}{
		{
			with: false,
			want: "f4",
		},
		{
			with: true,
			want: "f5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			out := bytes.NewBuffer(nil)
			_, err := WriteBool(out, tt.with)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(decodeHex(t, tt.want), out.Bytes()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func Test_WriteBytes(t *testing.T) {
	tests := []struct {
		with []byte
		want string
	}{
		{
			with: []byte{},
			want: "40",
		},
		{
			with: []byte{1, 2, 3, 4},
			want: "4401020304",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			out := bytes.NewBuffer(nil)
			_, err := WriteBytes(out, tt.with)
			if diff := cmp.Diff(decodeHex(t, tt.want), out.Bytes()); diff != "" {
				t.Fatal(diff)
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func Test_WriteString(t *testing.T) {
	tests := []struct {
		with string
		want string
	}{
		{
			with: "",
			want: "60",
		},
		{
			with: "a",
			want: "6161",
		},
		{
			with: "IETF",
			want: "6449455446",
		},
		{
			with: "\"\\",
			want: "62225c",
		},
		{
			with: "\u00fc",
			want: "62c3bc",
		},
		{
			with: "\u6c34",
			want: "63e6b0b4",
		},
		{
			with: "\xf0\x90\x85\x91", // "\ud800\udd51",
			want: "64f0908591",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			out := bytes.NewBuffer(nil)
			_, err := WriteString(out, tt.with)
			if diff := cmp.Diff(decodeHex(t, tt.want), out.Bytes()); diff != "" {
				t.Fatal(diff)
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func Test_WriteArray(t *testing.T) {
	tests := []struct {
		with []int
		want string
	}{
		{
			with: []int{},
			want: "80",
		},
		{
			with: []int{1, 2, 3},
			want: "83010203",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			out := bytes.NewBuffer(nil)
			_, err := WriteArrayHeader(out, uint64(len(tt.with)))
			if err != nil {
				t.Fatal(err)
			}
			for _, i := range tt.with {
				_, err = WriteSigned(out, int8(i))
				if err != nil {
					t.Fatal(err)
				}
			}
			if diff := cmp.Diff(decodeHex(t, tt.want), out.Bytes()); diff != "" {
				t.Fatal(diff)
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func Test_WriteMap(t *testing.T) {
	tests := []struct {
		with map[int]int
		want string
	}{
		{
			with: map[int]int{},
			want: "a0",
		},
		{
			with: map[int]int{1: 2, 3: 4},
			want: "a201020304",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			out := bytes.NewBuffer(nil)
			_, err := WriteMapHeader(out, uint64(len(tt.with)))
			if err != nil {
				t.Fatal(err)
			}
			for k, v := range tt.with {
				_, err = WriteSigned(out, int8(k))
				if err != nil {
					t.Fatal(err)
				}
				_, err = WriteSigned(out, int8(v))
				if err != nil {
					t.Fatal(err)
				}
			}
			if diff := cmp.Diff(decodeHex(t, tt.want), out.Bytes()); diff != "" {
				t.Fatal(diff)
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func Test_WriteTag(t *testing.T) {
	tests := []struct {
		with uint64
		want string
	}{
		{
			with: 0,
			want: "c0",
		},
		{
			with: 1,
			want: "c1",
		},
		{
			with: 23,
			want: "d7",
		},
		{
			with: 24,
			want: "d818",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			out := bytes.NewBuffer(nil)
			_, err := WriteTag(out, tt.with)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(decodeHex(t, tt.want), out.Bytes()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
