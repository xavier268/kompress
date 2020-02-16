package huffman

// repeat writer/reader is a preprocessing layer
// that will replace sequence of identical streams
// with an escape sequence, based upon a specific escape Symbol.
// The escape symbol is added to the initial alphabet, therefore
// increasing its length by exactly one.

// The encoding is as follows :
// x (for x < esc) 	=> x
// xx 				=> xx
// xxx				=> esc 0 x
// xxxx ... (n symbols, n<esc && n<=255) => esc (n-3) x
//
// This means that encoding is :
// esc n x  =>  xxxx.. repeated n+3 times)

// =============== repeatwriter ====================================

type repeatwriter struct {
	// here to write the processed symbols.
	swriter SymbolWriteCloser
	// escape symbol
	esc Symbol
	// last Symbol seen
	last Symbol
	// sequence lenth so far
	seq int
}

// interface checks
var _ SymbolWriteCloser = new(repeatwriter)

// NewRepeatWriter constructs a writer,
// where alphabetLength is the number of accepted input symbols.
func NewRepeatWriter(swr SymbolWriteCloser, alphabetLength int) SymbolWriteCloser {
	return newRepeatWriter(swr, alphabetLength)
}

// used for testing.
func newRepeatWriter(swr SymbolWriteCloser, alphabetLength int) *repeatwriter {
	r := new(repeatwriter)
	r.swriter = swr
	r.esc = Symbol(alphabetLength)
	return r
}

// Close flush pending sequence and close underlying writer.
func (rw *repeatwriter) Close() error {
	switch rw.seq {
	case 0: // ignore
	case 1:
		if err := rw.swriter.WriteSymbol(rw.last); err != nil {
			return err
		}
	case 2:
		if err := rw.swriter.WriteSymbol(rw.last); err != nil {
			return err
		}
		if err := rw.swriter.WriteSymbol(rw.last); err != nil {
			return err
		}
	default:
		if err := rw.swriter.WriteSymbol(rw.esc); err != nil {
			return err
		}
		if err := rw.swriter.WriteSymbol(Symbol(rw.seq - 3)); err != nil {
			return err
		}
		if err := rw.swriter.WriteSymbol(rw.last); err != nil {
			return err
		}
	}
	return rw.swriter.Close()
}

// WriteSymbol will write, but necessary flush, symbols.
func (rw *repeatwriter) WriteSymbol(s Symbol) error {
	if s >= rw.esc {
		panic("invalid input - not a valid symbol")
	}
	switch {
	case s == rw.last && rw.seq < 255+3 && rw.seq < int(rw.esc)+3:
		rw.seq++
		return nil
	case s == rw.last && (rw.seq == 255+3 || rw.seq == int(rw.esc)+3):
		if err := rw.swriter.WriteSymbol(rw.esc); err != nil {
			return err
		}
		if err := rw.swriter.WriteSymbol(Symbol(rw.seq - 3)); err != nil {
			return err
		}
		if err := rw.swriter.WriteSymbol(rw.last); err != nil {
			return err
		}
		rw.seq = 1
		return nil
	case s != rw.last:
		switch rw.seq {
		case 0:
			rw.seq++
			rw.last = s
			return nil
		case 1:
			if err := rw.swriter.WriteSymbol(rw.last); err != nil {
				return err
			}
			rw.seq = 1
			rw.last = s
			return nil
		case 2:
			if err := rw.swriter.WriteSymbol(rw.last); err != nil {
				return err
			}
			if err := rw.swriter.WriteSymbol(rw.last); err != nil {
				return err
			}
			rw.seq = 1
			rw.last = s
			return nil
		default:
			if err := rw.swriter.WriteSymbol(rw.esc); err != nil {
				return err
			}
			if err := rw.swriter.WriteSymbol(Symbol(rw.seq - 3)); err != nil {
				return err
			}
			if err := rw.swriter.WriteSymbol(rw.last); err != nil {
				return err
			}
			rw.last = s
			rw.seq = 1
			return nil
		}
	default:
		panic("invalid state)")
	}

}

//================ repeatreader ==============================

type repeatreader struct {
	rreader SymbolReader
	esc     Symbol
	seq     int    // remaining symbol to be read
	last    Symbol // symbol to be read
}

// NewRepeatReader constructor, alphetLength is without the REPEAT symbol.
func NewRepeatReader(sr SymbolReader, alphabetLength int) SymbolReader {
	return newRepeatReader(sr, alphabetLength)
}

// for testing purposes
func newRepeatReader(sr SymbolReader, alphabetLength int) *repeatreader {
	rr := new(repeatreader)
	rr.rreader = sr
	rr.esc = Symbol(alphabetLength)
	return rr
}

func (rr *repeatreader) ReadSymbol() (Symbol, error) {

	// Seq already avalable ...
	if rr.seq > 0 {
		rr.seq--
		return rr.last, nil
	}

	// now, seq == 0
	s, err := rr.rreader.ReadSymbol()
	if err != nil {
		return 0, err
	}

	// no escape symbol
	if s != rr.esc {
		return s, nil
	}

	// s == esc sequence
	len, err := rr.rreader.ReadSymbol()
	if err != nil {
		return 0, err
	}
	rr.last, err = rr.rreader.ReadSymbol()
	if err != nil {
		return 0, err
	}

	rr.seq = int(len) + 3 - 1
	return rr.last, nil

}
