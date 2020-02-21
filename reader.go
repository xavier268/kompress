package kompress

import (
	"errors"
	"io"
)

type hreader struct {
	*engine
	breader BitReader
	err     error
}

// compiler checks
var _ SymbolReader = new(hreader)

// NewReader can read symbols from bits.
func NewReader(br BitReader, eof Symbol, weights []int) SymbolReader {
	return newReader(br, eof, weights)
}

// NewReader can read symbols from bits.
func newReader(br BitReader, eof Symbol, weights []int) *hreader {
	hr := new(hreader)
	hr.engine = newEngine(weights, eof)
	hr.breader = br
	return hr
}

// ReadSymbol using the provided bit reader.
// It updates the actual frequency tables for the  symbol,
// but does not recompute the tree unless told to do so.
func (hr *hreader) ReadSymbol() (Symbol, error) {
	if hr.err != nil {
		return 0, hr.err
	}
	n := hr.root
	// loop until non leaf ...
	for n.id >= hr.len {
		if n == nil {
			return 0, errors.New("invalid bit sequence")
		}
		b, err := hr.breader.ReadBit()
		if err != nil {
			hr.err = err
			return 0, err
		}
		if b == 0 {
			n = n.child0
		} else {
			n = n.child1
		}
	}
	if Symbol(n.id) == hr.eof {
		hr.err = io.EOF
		return 0, hr.err
	}

	// increment actual freq
	hr.engine.actfreq[n.id]++

	return Symbol(n.id), nil
}
