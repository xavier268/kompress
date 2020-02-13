package huffman

import (
	"fmt"
	"io"
)

// LogWriter is a Writer that can write bits or Symbols
// to stdout for debugging.
type LogWriter struct{}

// compiler checks
var _ BitWriteCloser = new(LogWriter)
var _ SymbolWriteCloser = new(LogWriter)

// Close file
func (l *LogWriter) Close() error {
	fmt.Println("\n--- LogWriter was closed ---")
	return nil
}

// WriteBit to stdout
func (l *LogWriter) WriteBit(b Bit) error {
	if b == 0 {
		fmt.Print("0")
	} else {
		fmt.Print("1")
	}
	return nil
}

// WriteBits to stdout
func (l *LogWriter) WriteBits(bb ...Bit) error {
	for _, b := range bb {
		l.WriteBit(b)
	}
	return nil
}

// WriteSymbol to stdout
func (l *LogWriter) WriteSymbol(s Symbol) error {
	fmt.Printf("symbol(%d)", s)
	return nil
}

// ========================  BitBuffer ============================

// A BitBuffer can be read or written bits.
type BitBuffer struct {
	bits []Bit
}

// NewBitBuffer with optionnal preset bits.
func NewBitBuffer(bits ...Bit) *BitBuffer {
	bb := new(BitBuffer)
	bb.bits = make([]Bit, len(bits))
	for i, b := range bits {
		bb.bits[i] = b
	}
	return bb
}

// Compiler checks
var _ BitReader = new(BitBuffer)
var _ BitWriteCloser = new(BitBuffer)

// ReadBit from buffer
func (bb *BitBuffer) ReadBit() (Bit, error) {
	if len(bb.bits) == 0 {
		return 0, io.EOF
	}
	b := bb.bits[0]
	bb.bits = bb.bits[1:]
	if b == 0 {
		return 0, nil
	}
	return 1, nil
}

// Close does nothing
func (bb *BitBuffer) Close() error {
	return nil
}

// WriteBit to buffer
func (bb *BitBuffer) WriteBit(b Bit) error {
	bb.bits = append(bb.bits, b)
	return nil
}

// WriteBits to buffer
func (bb *BitBuffer) WriteBits(b ...Bit) error {
	bb.bits = append(bb.bits, b...)
	return nil
}

// Size of the buffer.
func (bb *BitBuffer) Size() int {

	return len(bb.bits)
}
