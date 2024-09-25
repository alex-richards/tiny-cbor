package cbor

import (
	"encoding/hex"
	"testing"
)

func decodeHex(tb testing.TB, encoded string) []byte {
	tb.Helper()

	decoded, err := hex.DecodeString(encoded)
	if err != nil {
		tb.Fatal(err)
	}

	return decoded
}
