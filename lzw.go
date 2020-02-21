package kompress

import "fmt"

// LZW algorith implementation.
//
// Principles for encoding :
// ========================
//
// 1. Dictionnary contains sequences <-> key; are multi symbol
// 2. Find longest dictionnary string that match current input
// 3. Replace output with key for that sequence
// 4. ALearn the sequence formed with the matching sequence
//    followed by the last symbol, until max table size is reached.
//
// Encoding changes the encoding alphabet, from Symbol1 to Symbol2
// Symbol1 is the original Symbol,
// Symbol2 has up to dicomax values,

//============ lzwwriter =================================

type lzwwriter struct {
	*lzwdico
	// sequence prepared for writing ...
	seq []Symbol
	w   SymbolWriteCloser
}

// NewLZWWriter constructor.
func NewLZWWriter(sw SymbolWriteCloser,
	alphaLenIn int, alphaLenOut int, maxSeq int) SymbolWriteCloser {
	return newlzwwriter(sw, alphaLenIn, alphaLenOut, maxSeq)
}

func newlzwwriter(sw SymbolWriteCloser,
	alphaLenIn int, alphaLenOut int, maxSeq int) *lzwwriter {

	l := new(lzwwriter)
	l.lzwdico = newlzwdico(alphaLenIn, alphaLenOut, maxSeq)
	l.w = sw
	return l
}

// WriteSymbol incoming symbol s, emitting symbol as needed.
func (lz *lzwwriter) WriteSymbol(s Symbol) error {

	if int(s) > lz.nbIn {
		fmt.Println("input  token :", s)
		panic("invalid input  token (too large)")
	}

	lz.seq = append(lz.seq, s)

	// check if full seq is known ...
	code, ok := lz.getCode(lz.seq)
	switch {
	case ok && len(lz.seq) < lz.seqMax:
		// seq is known, but not yet max, try to get more ...
		return nil
	case ok && len(lz.seq) >= lz.seqMax:
		// seq is known, maximum reached, do not try further
		err := lz.w.WriteSymbol(code)
		if err != nil {
			return err
		}
		// learn, then clear what was learnt
		lz.learn(lz.seq...)
		lz.seq = []Symbol{}
		return nil

	case !ok && len(lz.seq) > 0:

		// here, previous sequence cannot be decoded ..
		// let's decode the first sub sequence that would work.
		// this can happen if the table is full.

		// adjust flush the same way.

		// start with the beginning of sequence, until we cannot decode,
		// then emit, and start again ...

		// find the longuest subsequence
		i := 1
		for _, ok = lz.getCode(lz.seq[:i]); ok && i < len(lz.seq); i++ {
		}
		i-- // i is now the longuest valid sub sequence
		code, ok := lz.getCode(lz.seq[:i])
		if !ok {
			lz.dump()
			fmt.Println("i=", i, "full seq : ", lz.seq, "subseq =", lz.seq[:i], ", s = ", s)
			panic("invalid state")
		}
		err := lz.w.WriteSymbol(code)
		if err != nil {
			return err
		}
		// learn and clear what was learnt
		lz.learn(lz.seq[:i]...)
		lz.seq = lz.seq[i:]
		return nil
	}
	panic("invalid state")
}

// Close (and flush ...)
func (lz *lzwwriter) Close() error {
	if len(lz.seq) != 0 {
		code, ok := lz.getCode(lz.seq)
		if ok {
			err := lz.w.WriteSymbol(code)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("Flushing sequence ", lz.seq)
			lz.dump()
			panic("invalid state")
		}
	}
	return lz.w.Close()
}

//============ lzwreader =================================

type lzwreader struct {
	*lzwdico
	// available symbols to be read
	avail []Symbol
	r     SymbolReader
}

// NewLZWReader constructor.
func NewLZWReader(sr SymbolReader,
	alphaLenIn int, alphaLenOut int, maxSeq int) SymbolReader {
	return newlzwreader(sr, alphaLenIn, alphaLenOut, maxSeq)
}

func newlzwreader(sr SymbolReader,
	alphaLenIn int, alphaLenOut int, maxSeq int) *lzwreader {

	l := new(lzwreader)
	l.lzwdico = newlzwdico(alphaLenIn, alphaLenOut, maxSeq)
	l.r = sr
	return l
}

func (lr *lzwreader) ReadSymbol() (Symbol, error) {

	if len(lr.avail) == 0 {
		lr.avail = []Symbol{} // memory control ...

		// here, nothing is ready to be read,
		// lets read more in the buffers
		s2, err := lr.r.ReadSymbol()
		if err != nil {
			return 0, err
		}
		sin, ok := lr.getSeq(s2)
		if !ok || len(sin) == 0 {
			lr.lzwdico.dump()
			fmt.Println("Could not find symbol out : ", s2)
			panic("out symbol unknown in dictionnary - invalid state")
		}
		// learn and update available symbols
		lr.learn(sin...)
		lr.avail = sin

	} // if len avail == 0

	if len(lr.avail) == 0 {
		panic("avail should not be empty - invalid state")
	}

	s := lr.avail[0]
	lr.avail = lr.avail[1:]

	return s, nil
}
