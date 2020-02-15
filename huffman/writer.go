package huffman

import "errors"

type hwriter struct {
	*engine
	bwriter BitWriteCloser
	err     error
}

// NewWriter constructs a hwriter, a SymbolWriteCloser,
// that can encode symbols into bits.
// eof is the Symbole that is written when file is closed.
func NewWriter(bw BitWriteCloser, eof Symbol, weights []int) SymbolWriteCloser {
	return newWriter(bw, eof, weights)
}

func newWriter(bw BitWriteCloser, eof Symbol, weights []int) *hwriter {
	h := new(hwriter)
	h.engine = newEngine(weights, eof)
	h.bwriter = bw
	return h
}

func (h *hwriter) Close() error {
	h.err = h.WriteSymbol(h.eof)
	if h.err != nil {

		return h.err
	}
	h.err = h.Close()
	return h.err
}

// WriteSymbol writes the symbol, incrementing the actual frequency
// everytime the symbol is encoded.
func (h *hwriter) WriteSymbol(s Symbol) error {
	if h.err != nil {
		return h.err
	}
	if int(s) >= h.len {
		h.err = errors.New("Symbol is not part of alphabet")
		return h.err
	}
	// increment actual frequency, except for EOF
	if s != h.eof {
		h.engine.actfreq[s]++
	}

	n := &h.nodes[s]
	h.err = h.writeSymbol(n, nil)

	return h.err
}

func (h *hwriter) writeSymbol(n, from *node) error {

	if n == nil {
		return nil
	}

	// recurrence on parent
	if n.parent != nil {
		err := h.writeSymbol(n.parent, n)
		if err != nil {
			return err
		}
	}

	// post process bits
	switch from {
	case nil: // do nothing, leaf
	case n.child0:
		err := h.bwriter.WriteBit(Bit(0))
		if err != nil {
			return err
		}
	case n.child1:
		err := h.bwriter.WriteBit(Bit(1))
		if err != nil {
			return err
		}
	}
	return nil
}
