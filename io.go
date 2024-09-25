package cbor

import "io"

type peekReader struct {
	r  io.Reader // wrapped reader
	p  byte      // peeked byte
	pv bool      // peeded valid
}

func (r *peekReader) Read(out []byte) (int, error) {
	ol := len(out)
	if ol == 0 {
		return 0, nil
	}

	if r.pv {
		out[0] = r.p
		r.pv = false
		if ol == 1 {
			return 1, nil
		}

		n, err := r.r.Read(out[1:ol])
		return n + 1, err
	}

	return r.r.Read(out)
}

func (r *peekReader) PeekByte() (byte, error) {
	if r.pv {
		return r.p, nil
	}

	var pb = []byte{0}
	n, err := r.r.Read(pb)

	if n != 1 {
		return 0, err
	}

	r.p = pb[0]
	r.pv = true
	return r.p, nil
}
