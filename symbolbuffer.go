package kompress

import "io"

// SymbolBuffer is a FIFO for Symbols
type SymbolBuffer struct {
	buf []Symbol
}

var _ SymbolWriteCloser = new(SymbolBuffer)

// NewSymbolBuffer constructor
func NewSymbolBuffer() *SymbolBuffer {
	sb := new(SymbolBuffer)
	sb.buf = make([]Symbol, 0)
	return sb
}

// Close does nothing.
func (sb *SymbolBuffer) Close() error {
	return nil
}

// WriteSymbol in buffer
func (sb *SymbolBuffer) WriteSymbol(s Symbol) error {
	sb.buf = append(sb.buf, s)
	return nil
}

// ReadSymbol returns eof if empty.
func (sb *SymbolBuffer) ReadSymbol() (Symbol, error) {
	if len(sb.buf) == 0 {
		return 0, io.EOF
	}
	s := sb.buf[0]
	sb.buf = sb.buf[1:]
	return s, nil
}
