package kompress

import (
	"io"
)

type kdelta struct {
	buf     []byte       // last bytes seen
	capa    int          // max capacity of the buffer
	exp     map[int]byte // expected bytes for a givent context hash
	hashmax int          // max size of the hash map
}

func (k *kdelta) reset(capa int) {
	k.capa = capa
	k.hashmax = 1000 * k.capa
	k.exp = make(map[int]byte, 1000)
	k.buf = make([]byte, 0, k.capa)
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
	k.kwriter.reset(w)
	k.kdelta.reset(capa)

	return k
}

// NewKdeltaReader constructs an io.Reader that will retrieve original bytes from delta encoded bytes.
func NewKdeltaReader(r io.Reader, capa int) io.Reader {

	k := new(kdeltaReader)
	k.kreader.reset(r)
	k.kdelta.reset(capa)

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

// expect gets the expected byte for the given context hash.
func (k *kdelta) expect(h int) byte {
	return k.exp[h]
}

// learn provides context learning info
// it also updates the context buffer.
func (k *kdelta) learn(h int, b byte) {
	k.exp[h] = b
	k.buf = append(k.buf, b)
	if len(k.buf) > k.capa {
		k.buf = k.buf[1:]
	}
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
		e := k.expect(h)

		// compute delta, write it
		d := e ^ b
		k.emit(d)
		if k.err != nil {
			return count, k.err
		}
		count++

		// update expected byte and context
		k.learn(h, b)
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
		e = k.expect(h)

		// compute true value b, write it
		b = e ^ d
		bb[i] = b

		// update expected byte and context
		k.learn(h, b)
	}
	return count, err
}
