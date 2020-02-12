package huffman

import "io"

// Bit are internally represented as bytes,
// where a 0 bit is a zero byte, all other values are 1s.
type Bit byte

// Symbol are to be transcoded to/from bits.
type Symbol uint

// BitWriteCloser allows wrinting bits.
type BitWriteCloser interface {
	WriteBit(Bit) error
	io.Closer
}

// SymbolWriteCloser writes symbols.
type SymbolWriteCloser interface {
	WriteSymbol(Symbol) error
	io.Closer
}

// SymbolReader reads Symbol.
type SymbolReader interface {
	ReadSymbol() (Symbol, error)
}

// BitReader reads bits, one by one.
type BitReader interface {
	ReadBit() (Bit, error)
}
