package kompress

import (
	"errors"
	"io"
	"os"
)

type kwriter struct {
	writer io.Writer
	err    error
}

type kreader struct {
	reader io.Reader
	err    error
}

func (k *kwriter) reset(w io.Writer) {
	k.err = nil
	if w == nil {
		k.writer = os.Stdout
	} else {
		k.writer = w
	}
}

func (k *kreader) reset(r io.Reader) {
	k.err = nil
	if r == nil {
		k.reader = os.Stdin
	} else {
		k.reader = r
	}
}

// Close flush all data, and attempt to close underlying writer,
// if it is possible.
func (k *kwriter) Close() error {
	k.err = errors.New("file was closed")
	if k.writer != os.Stdout {
		c, ok := k.writer.(io.Closer)
		if ok {
			return c.Close()
		}
	}
	return nil
}

// write the bytes in the underlying writer
func (k *kwriter) emit(bb ...byte) {

	if k.err != nil || len(bb) == 0 {
		return
	}

	var n int
	n, k.err = k.writer.Write(bb)
	if n != len(bb) {
		k.err = errors.New("writing is not possible")
	}
}

// readBytes reads a single byte, returning error if couldn't read.
// Errors are transmitted only if no bytes can be read.
func (k *kreader) readByte() (byte, error) {

	if k.err != nil {
		return 0, k.err
	}

	b := make([]byte, 1, 1)
	n, e := k.reader.Read(b)
	if n != 0 {
		return b[0], nil
	}
	return 0, e
}
