package huffman

import "fmt"

// lzwdico provides ways to store, index, retrieve sequence of symbols.

type lzwdico struct {
	// Alphabet out and in
	nbIn, nbOut int
	// max sequence length
	seqMax int
	// from a symbol in sequence to a symbole out code
	codes map[string]Symbol
	// form a symbol out code to a symbol in sequence
	seqs map[Symbol][]Symbol
	// buf is a buffer for everything we see
	// its length may not exceed seqMax
	buf []Symbol
}

func newlzwdico(nbIn, nbOut, seqMax int) *lzwdico {
	if seqMax < 3 || nbOut <= nbIn || seqMax >= nbOut {
		fmt.Println(nbIn, nbOut, seqMax)
		panic("invalid inputs")
	}
	d := new(lzwdico)
	d.nbIn, d.nbOut, d.seqMax = nbIn, nbOut, seqMax
	d.codes = make(map[string]Symbol)
	d.seqs = make(map[Symbol][]Symbol)
	// init the dico with codes for in symbols
	for i := 0; i < nbIn; i++ {
		d.learn(Symbol(i))
		d.buf = []Symbol{}
	}
	return d
}

// learn and add the Symbols to the buffer,
// udating the dictionnary on the way.
func (d *lzwdico) learn(ss ...Symbol) {

	// Do not update if all out symbols are used.
	if len(d.codes) >= d.nbOut {
		return
	}

	for _, s := range ss {

		// include in buffer
		d.buf = append(d.buf, s)
		if len(d.buf) > d.seqMax {
			d.buf = d.buf[1:]
		}

		// refresh subseq that are not yet known
		for i := 0; i < len(d.buf); i++ {
			// don't overwrite existing sequences !
			if _, ok := d.getCode(d.buf[i:]); !ok && len(d.codes) < d.nbOut {
				code := Symbol(len(d.codes))
				d.codes[fmt.Sprint(d.buf[i:])] = code
				d.seqs[code] = d.buf[i:]
			}
		}

	} // for s

}

// getCode return the symbol out for the sequence symbol in
func (d *lzwdico) getCode(sin []Symbol) (code Symbol, ok bool) {
	code, ok = d.codes[fmt.Sprint(sin)]
	return code, ok
}

// get the sequence of symbol in from the symbol out code.
func (d *lzwdico) getSeq(code Symbol) (sin []Symbol, ok bool) {
	sin, ok = d.seqs[code]
	return sin, ok
}

func (d *lzwdico) dump() {
	fmt.Println("There are ", len(d.codes), " codes in the dictionnary")
	fmt.Println("Internal buffer : ", d.buf)
	for code, seq := range d.seqs {
		fmt.Println("Out : ", code, "\t In : ", seq)
	}
}
