package kompress

import (
	"errors"
	"fmt"
	"io"
)

// A Compresser is able to compress or decompress data.
type Compresser interface {
	Compress(io.Reader, io.Writer) error
	Decompress(io.Reader, io.Writer) error
	Reset()
}

// ErrorInvalidCompressionFormat is a compression grammar error
var ErrorInvalidCompressionFormat = errors.New("invalid compression format")

// Dumper is an io.Writer that dumps its input to the screen.
type Dumper struct{}

func (w *Dumper) Write(b []byte) (n int, err error) {

	if len(b) == 0 {
		return 0, io.EOF
	}

	fmt.Printf("%v\n", b)
	return len(b), nil

}
