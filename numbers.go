package cbor

func shiftBytesInto[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64](b []byte) (o T) {
	n := len(b)
	for i := range n {
		o |= T(b[i]) << (8 * (n - (i + 1)))
	}
	return
}

func shiftBytesFrom[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64](v T, b []byte) {
	n := len(b)
	for i := range n {
		b[i] = byte(v >> (8 * ((n - i) - 1)))
	}
}
