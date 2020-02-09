package kompress

import (
	"io"
	"os"
)

type kdelta struct {
	buf     []byte       // last bytes seen
	capa    int          // max capacity of the buffer
	exp     map[int]byte // expected bytes for a givent context hash
	hashmax int          // max size of the hash map
}

type kdeltaWriter struct {
	kwriter
	kdelta
}

type kdeltaReader struct {
	kreader
	kdelta
}

// NewKdeltaWriter constructs an io.WriteCloser that will output delta encoded bytes.
func NewKdeltaWriter(w io.Writer, capa int) io.WriteCloser {

	k := new(kdeltaWriter)
	k.capa = capa
	if w == nil {
		k.writer = os.Stdout
	} else {
		k.writer = w
	}
	k.hashmax = 1000 * k.capa
	k.exp = make(map[int]byte, 1000)
	k.buf = make([]byte, 0, k.capa)
	return k
}

// NewKdeltaReader constructs an io.Reader that will retrieve original bytes from delta encoded bytes.
func NewKdeltaReader(w io.Reader, capa int) io.Reader {

	k := new(kdeltaReader)
	k.capa = capa
	if w == nil {
		k.reader = os.Stdout
	} else {
		k.reader = w
	}
	k.hashmax = 1000 * k.capa
	k.exp = make(map[int]byte, 1000)
	k.buf = make([]byte, 0, k.capa)
	return k
}

// hash produce a hash from the current byte context
func (k *kdelta) hash() int {
	var r int
	for _, b := range k.buf {
		r = (r*7 + 11*int(b) + 3) % k.hashmax
	}
	return r
}

// Write will produce the delta encoded bytes.
func (k *kdeltaWriter) Write(bb []byte) (int, error) {

	if k.err != nil {
		return 0, k.err
	}

	var count int

	for _, b := range bb {

		// compute context, get expected byte
		h := k.hash()
		e := k.exp[h]

		// compute delta, write it
		d := e ^ b
		k.emit(d)
		if k.err != nil {
			return count, k.err
		}
		count++

		// update expected byte and context
		k.exp[h] = b
		k.buf = append(k.buf, b)
		if len(k.buf) > k.capa {
			k.buf = k.buf[1:]
		}
	}

	return count, k.err

}

// Read will read delta encoded bytes into the original bytes
func (k *kdeltaReader) Read(bb []byte) (int, error) {

	if k.err != nil {
		return 0, k.err
	}

	var b, d, e byte
	var err error
	var count, h int

	for i := 0; i < len(bb); i++ {

		d, err = k.readByte()
		if err != nil {
			k.err = err
			return count, err
		}
		count++

		// compute context, get expected byte
		h = k.hash()
		e = k.exp[h]

		// compute true value b, write it
		b = e ^ d
		bb[i] = b

		// update expected byte and context
		k.exp[h] = b
		k.buf = append(k.buf, b)
		if len(k.buf) > k.capa {
			k.buf = k.buf[1:]
		}
	}
	return count, err
}
