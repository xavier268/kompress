package huffman

// Predictive encoding.
// This transformation does not change the alphabet, not the length.
// The transformation tries to predict any value below len,
// encoding the observed delta(xor).
// Predictions are made based on the last buffered context.

// ratio # seq/nb bytes
const hashmax = 1000

type predict struct {
	// number of symbol in alphabet
	len int
	// internal context buffer
	buf []Symbol
	// buffer index (rotating over length)
	bufi int
	// map hash context to predictions
	preds map[int]Symbol
}

type predictwriter struct {
	*predict
	writer SymbolWriteCloser
}

// NewPredictWriter constructor
func NewPredictWriter(sw SymbolWriteCloser, alphabetLen int, bufSize byte) SymbolWriteCloser {
	w := new(predictwriter)
	w.writer = sw
	w.predict = newPredict(alphabetLen, bufSize)
	return w
}

func (pw *predictwriter) Close() error {
	return pw.writer.Close()
}

func (pw *predictwriter) WriteSymbol(s Symbol) error {
	return pw.WriteSymbol(pw.forward(s))
}

type predictreader struct {
	*predict
	reader SymbolReader
}

// NewPredictReader constructor
func NewPredictReader(r SymbolReader, alphabetLen int, bufSize byte) SymbolReader {
	p := new(predictreader)
	p.reader = r
	p.predict = newPredict(alphabetLen, bufSize)
	return p
}

func (pr *predictreader) ReadSymbol() (Symbol, error) {
	delta, err := pr.reader.ReadSymbol()
	if err != nil {
		return 0, err
	}
	return pr.predict.backward(delta), nil
}

// Will build history of bufSize bytes. Suited for periods of bufSize+1.
func newPredict(alphabetLen int, bufSize byte) *predict {
	p := new(predict)
	p.len = alphabetLen
	p.buf = make([]Symbol, bufSize, bufSize)
	p.preds = make(map[int]Symbol, hashmax*int(bufSize))
	return p
}

// Compute context hash value
func (p *predict) hash() int {
	h := 0
	for _, b := range p.buf {
		h += (7*int(b) + 5) % hashmax
	}
	return h
}

// Learn the new value as the
// prediction for the previous context.
// then update the context.
// Return the delta value to store.
func (p *predict) forward(s Symbol) (delta Symbol) {
	// get prvious prediction ... 0 if prediction worked
	h := p.hash()
	delta = Symbol((-int(p.preds[h]) + int(s) + p.len) % p.len)
	// learn the new value
	p.preds[h] = s
	// update context buffer, increment buffer pointer
	p.buf[p.bufi] = s
	p.bufi = (p.bufi + 1) % (len(p.buf))
	return delta
}

// Reads the delta, evalute prediction,
// recompute original value to return,
// use original value to update context and learn map.
// Return the original value.
func (p *predict) backward(delta Symbol) (s Symbol) {
	// compute original
	h := p.hash()
	s = Symbol((int(p.preds[h]) + int(delta) + p.len) % p.len)
	// learn the new value
	p.preds[h] = s
	// update context buffer
	p.buf[p.bufi] = s
	p.bufi = (p.bufi + 1) % (len(p.buf))
	return s
}
