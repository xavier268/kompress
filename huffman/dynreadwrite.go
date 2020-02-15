package huffman

// dwriter and dreader are dynamic reader/writer,
// that will reconstruct thier tree from time to time,
// to match the actual frequencies.

//============ dwriter ====================

type dwriter struct {
	*hwriter
}

// NewDWriter provides a SymbolWriteCloser that dynamically adapts its decoding tree.
func NewDWriter(bw BitWriteCloser, eof Symbol, weights []int, sch Scheduler) SymbolWriteCloser {
	return newDWriter(bw, eof, weights, sch)
}

// Same, for debugging purposes, exposing all internals.
func newDWriter(bw BitWriteCloser, eof Symbol, weights []int, sch Scheduler) *dwriter {
	d := new(dwriter)
	d.hwriter = newWriter(bw, eof, weights)
	d.scheduler = sch
	return d
}

func (dw *dwriter) WriteSymbol(s Symbol) error {
	err := dw.hwriter.WriteSymbol(s)
	if dw.scheduler != nil && dw.scheduler(dw.engine) {
		dw.hwriter = newWriter(dw.bwriter, dw.eof, dw.actfreq)
	}

	return err
}

//========== dreader ================

type dreader struct {
	*hreader
}

// NewDReader provides a SymbolReader that dynamically adapts its decoding tree.
func NewDReader(br BitReader, eof Symbol, weights []int, sch Scheduler) SymbolReader {
	return newDReader(br, eof, weights, sch)
}
func newDReader(br BitReader, eof Symbol, weights []int, sch Scheduler) *dreader {
	dr := new(dreader)
	dr.hreader = newReader(br, eof, weights)
	dr.scheduler = sch
	return dr
}

func (dr *dreader) ReadSymbol() (Symbol, error) {
	s, err := dr.hreader.ReadSymbol()
	if dr.scheduler != nil && dr.scheduler(dr.engine) {
		dr.hreader = newReader(dr.breader, dr.eof, dr.actfreq)
	}
	return s, err
}
