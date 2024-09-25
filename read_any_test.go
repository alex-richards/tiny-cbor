//go:build !cbor_no_readany

package cbor

import (
	"bytes"
	"encoding/hex"
	"io"
	"testing"
)

func Test_ReadAny(t *testing.T) {
	for _, tt := range tests_ExampleEncoded {
		t.Run(tt.encoded, func(t *testing.T) {
			decoded := decodeHex(t, tt.encoded)

			in := bytes.NewReader(decoded)

			v, err := ReadAny(in)
			if err != nil && err != ErrOverflow {
				t.Fatal(err)
			}

			n := in.Len()
			if n != 0 {
				rem := make([]byte, n)
				in.Read(rem)
				t.Fatalf("trailing data - %v - %s", v, hex.EncodeToString(rem))
			}
		})
	}
}

func Benchmark_ReadAny(b *testing.B) {
	for _, tt := range tests_ExampleEncoded {
		encoded := decodeHex(b, tt.encoded)

		in := bytes.NewReader(encoded)

		b.Run(tt.encoded, func(b *testing.B) {
			for range b.N {
				in.Seek(0, io.SeekStart)
				_, err := ReadAny(in)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
