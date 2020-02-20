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
		d.learn([]Symbol{Symbol(i)})
	}
	return d
}

// learn the new sequence, if it can, returning the code symbol out
func (d *lzwdico) learn(sin []Symbol) {
	if len(sin) == 0 ||
		len(sin) >= d.seqMax ||
		len(d.codes) >= d.nbOut {
		// ignore ...
		return
	}
	// test if already known, don't change it and ignore ...
	if _, ok := d.codes[fmt.Sprint(sin)]; ok {
		return
	}
	// actual update
	code := Symbol(len(d.codes))
	d.codes[fmt.Sprint(sin)] = code
	d.seqs[code] = sin
	return
}

// lean all the sub sequences
func (d *lzwdico) learnall(sin []Symbol) {

	for i := 0; i <= len(sin); i++ {
		for j := i + 1; j < len(sin); j++ {
			//fmt.Println("learning : ", sin[i:j])
			d.learn(sin[i:j])
		}
	}

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
	for code, seq := range d.seqs {
		fmt.Println("Out : ", code, "\t In : ", seq)
	}
}
