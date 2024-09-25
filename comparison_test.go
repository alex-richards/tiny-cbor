//go:build cbor_comparison

package cbor

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	fxcbor "github.com/fxamacker/cbor/v2"
	"github.com/google/go-cmp/cmp"
)

func Test_ReadAny_Comparison(t *testing.T) {
	for _, tt := range tests_ExampleEncoded {
		t.Run(tt.encoded, func(t *testing.T) {
			encoded := decodeHex(t, tt.encoded)
			in := bytes.NewReader(encoded)

			var err error
			var out any
			var outFx any

			{
				in.Seek(0, io.SeekStart)
				out, err = ReadAny(in)
				if err != nil {
					t.Fatal(err)
				}
			}

			{
				in.Seek(0, io.SeekStart)
				err = fxcbor.Unmarshal(encoded, &outFx)
				if err != nil {
					t.Fatal(err)
				}
			}

			if diff := cmp.Diff(out, outFx); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
func Benchmark_ReadAny_Comparison(b *testing.B) {
	for _, tt := range tests_ExampleEncoded {
		encoded := decodeHex(b, tt.encoded)
		in := bytes.NewReader(encoded)

		var out any
		var err error

		b.Run(fmt.Sprintf("%s", tt.encoded), func(b *testing.B) {
			for range b.N {
				in.Seek(0, io.SeekStart)
				out, err = ReadAny(in)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(fmt.Sprintf("fx %s", tt.encoded), func(b *testing.B) {
			for range b.N {
				in.Seek(0, io.SeekStart)
				err = fxcbor.Unmarshal(encoded, &out)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
