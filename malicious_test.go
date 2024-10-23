package cbor

import (
	"bytes"
	"io"
	"testing"
)

var malicious = []byte{0x9B, 0x00, 0x00, 0x42, 0xFA, 0x42, 0xFA, 0x42, 0xFA, 0x42}

func Test_Malicious(t *testing.T) {
    t.Run("ReadAny", func(t *testing.T) {
		in := bytes.NewBuffer(malicious)
        v, err := ReadAny(in)
        if v != nil {
            t.Fatal("unexpected value")
        }
        if err != io.EOF {
            t.Fatal("expected EOF")
        }
    })
    t.Run("ReadOver", func(t *testing.T) {
		in := bytes.NewBuffer(malicious)
        err := ReadOver(in)
        if err != io.EOF {
            t.Fatal("expected EOF")
        }
    })
    t.Run("ReadRaw", func(t *testing.T) {
		in := bytes.NewBuffer(malicious)
        out := bytes.NewBuffer(nil)
        err := ReadRaw(in, out)
        if err != io.EOF {
            t.Fatal("expected EOF")
        }
    })
}
