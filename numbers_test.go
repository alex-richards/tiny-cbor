package cbor

import (
	"encoding/hex"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_shiftBytesInto(t *testing.T) {
	tests := []struct {
		got  string
		want uint64
	}{
		{
			got:  "00",
			want: 0x00,
		},
		{
			got:  "0102030405060708",
			want: 0x0102030405060708,
		},
		{
			got:  "01020304050607080910",
			want: 0x0304050607080910,
		},
	}

	for _, tt := range tests {
		t.Run(tt.got, func(t *testing.T) {
			got := shiftBytesInto[uint64](decodeHex(t, tt.got))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func Test_shiftBytesFrom(t *testing.T) {
	tests := []struct {
		got  uint64
		want string
	}{
		{
			got:  0x00,
			want: "0000000000000000",
		},
		{
			got:  0x0102030405060708,
			want: "0102030405060708",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := make([]byte, 8)
			shiftBytesFrom(tt.got, got)
			if diff := cmp.Diff(tt.want, hex.EncodeToString(got)); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
