package kompress

import (
	"errors"
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
