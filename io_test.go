package cbor

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_peekReader(t *testing.T) {
	t.Run("Read", func(t *testing.T) {
		in := []byte{1, 2, 3, 4, 5, 6}
		out := make([]byte, len(in))

		reader := peekReader{r: bytes.NewReader(in)}

		_, err := reader.Read(out)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(in, out); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("Read Peak", func(t *testing.T) {
		in := []byte{1, 2, 3, 4, 5, 6}
		out := make([]byte, len(in))

		reader := peekReader{r: bytes.NewReader(in)}

		n, err := reader.Read(out[:3])
		if n != 3 {
			t.Log(out)
			t.Fatal(err)
		}

		p, err := reader.PeekByte()
		if p != 4 {
			t.Fatal("unexpected value")
		}

		n, err = reader.Read(out[3:])
		if n != 3 {
			t.Fatal(err)
		}

		if diff := cmp.Diff(in, out); diff != "" {
			t.Fatal(diff)
		}
	})
}
