package kompress

import (
	"bytes"
	"io"
)

// Kcombine is a Compresser made of successive compressors.
type Kcombine struct {
	// The array of compresser to be applied.
	cc []Compresser
}

// Reset all underlying compressers.
func (k *Kcombine) Reset() {
	for _, kk := range k.cc {
		kk.Reset()
	}
}

// Compress applying successively the underlying Compressers.
// TODO - error management
// TODO - check what happens with large files and small buffers ?
func (k *Kcombine) Compress(in io.Reader, out io.Writer) error {

	bin := in

	for _, kk := range k.cc {
		buf := bytes.NewBuffer(nil)
		kk.Compress(bin, buf)
		bin = buf
	}
	io.Copy(out, bin)
	return nil
}

// Decompress using Compressers in reverse order.
func (k *Kcombine) Decompress(in io.Reader, out io.Writer) error {
	panic("Not implemented")
}
