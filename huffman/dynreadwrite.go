package huffman

// dwriter and dreader are dynamic reader/writer,
// that will reconstruct thier tree from time to time,
// to match the actual frequencies.

//============ dwriter ====================

type dwriter struct {
	*hwriter
}

// NewDWriter provides a SymbolWriteCloser that dynamically adapts its decoding tree.
func NewDWriter(bw BitWriteCloser, eof Symbol, weights []int) SymbolWriteCloser {
	return newDWriter(bw, eof, weights)
}

// Same, for debugging purposes, exposing all internals.
func newDWriter(bw BitWriteCloser, eof Symbol, weights []int) *dwriter {
	d := new(dwriter)
	d.hwriter = newWriter(bw, eof, weights)
	return d
}

func (dw *dwriter) WriteSymbol(s Symbol) error {
	err := dw.hwriter.WriteSymbol(s)
	if dw.shouldUpdate() {
		dw.hwriter = newWriter(dw.bwriter, dw.eof, dw.actfreq)
	}

	return err
}

//========== dreader ================

type dreader struct {
	*hreader
}

// NewDReader provides a SymbolReader that dynamically adapts its decoding tree.
func NewDReader(br BitReader, eof Symbol, weights []int) SymbolReader {
	return newDReader(br, eof, weights)
}
func newDReader(br BitReader, eof Symbol, weights []int) *dreader {
	dr := new(dreader)
	dr.hreader = newReader(br, eof, weights)
	return dr
}

func (dr *dreader) ReadSymbol() (Symbol, error) {
	s, err := dr.hreader.ReadSymbol()
	if dr.shouldUpdate() {
		dr.hreader = newReader(dr.breader, dr.eof, dr.actfreq)
	}
	return s, err
}

//=========== scheduler =============

// Common scheduling function to decide when to update.
// Used for reading and writing.
// Should be deterministic (no time-based decisions !)
func (e *engine) shouldUpdate() bool {
	return true
}
