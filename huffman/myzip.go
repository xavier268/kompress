package huffman

import (
	"errors"
	"io"
)

// myZipWriter is able to zip a byte stream into a byte stream,
// using some preprocessing and an adaptative huffman tree compression.
type myZipWriter struct {
	writer SymbolWriteCloser
}

// NewMyZipWriter constructor.
func NewMyZipWriter(iow io.Writer) io.WriteCloser {

	// Write bits to bytes
	var b2bw BitWriteCloser = NewBitToByteWriter(iow)

	// transform symbols to bits
	eof := Symbol(256)
	weights := [257]int{}
	sch := SCAlways
	var wrt SymbolWriteCloser = NewDWriter(b2bw, eof, weights[:], sch)

	return &myZipWriter{wrt}

}

// Close attemps to close & flush everything it can...
// It is required to call close when finished compressing data.
func (mz *myZipWriter) Close() error {
	return mz.writer.Close()
}

// Write the data down the compression engine.
func (mz *myZipWriter) Write(data []byte) (n int, err error) {
	n = 0
	for _, d := range data {
		if err := mz.writer.WriteSymbol(Symbol(d)); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// myZipReader reads from compressed bytes to original bytes.
type myZipReader struct {
	reader SymbolReader
}

// NewMyZipReader constructor.
func NewMyZipReader(ior io.Reader) io.Reader {

	var br BitReader = NewBitFromByteReader(ior)

	eof := Symbol(256)
	weights := [257]int{}
	sch := SCAlways
	var dr SymbolReader = NewDReader(br, eof, weights[:], sch)

	return &myZipReader{dr}

}

func (mz *myZipReader) Read(bb []byte) (int, error) {
	for i := range bb {
		s, err := mz.reader.ReadSymbol()
		if err != nil {
			return i, err
		}
		if int(s) > 255 {
			return i, errors.New("invalid symbol - not a byte")
		}
		bb[i] = byte(s)
	}
	return len(bb), nil
}
