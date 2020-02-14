package huffman

import (
	"fmt"
	"io"
	"testing"
)

func TestWriteRead2ExpectPanic(t *testing.T) {
	freq := []int{0, 2, 8, 1, 3, 5, 1}
	syms := []Symbol{0, 2, 2, 1, 3, 4, 5, 6}
	eof := Symbol(6) // here, never written

	// Catch expected panic, because eof will not be kept.
	defer expectPanic(t)

	writeReadSymbol(t, syms, eof, freq)
}

// writeReadSymbol will write all symbols,
// log compression size, then read them back,
// and check identity.
func writeReadSymbol(t *testing.T, syms []Symbol, eof Symbol, freq []int) {

	buf := NewBitBuffer()

	if buf.Size() != 0 {
		t.Fatal("non empty buffer ?")
	}

	// write all symbols
	w := newWriter(buf, 6, freq)
	for i, s := range syms {
		err := w.WriteSymbol(s)
		if err != nil {
			w.Dump()
			t.Log("Last symbol written", syms[:i+1])
			t.Log("error :", err)
			t.Fatal(err)
		}
	}
	// Compression ratio
	fmt.Println("Compressed from\t", len(syms)*8, " bits\tto", buf.Size(), " bits")
	if buf.Size() >= len(syms)*8 {
		t.Fatal("no actual compression !?")
	}

	// read all symbols
	r := newReader(buf, 6, freq)
	var got []Symbol
	for buf.Size() > 0 {
		s, err := r.ReadSymbol()
		if err == io.EOF {
			break
		}
		if err != nil {
			r.Dump()
			t.Fatal("Unexpected Read error : ", s, err)
		}
		got = append(got, s)

	}

	if len(got) != len(syms) {
		//r.Dump()
		t.Log("Got  : ", got)
		t.Log("Want : ", syms)
		panic("Length do not match !")
	}
	for s := range got {
		if got[s] != syms[s] {
			r.Dump()
			t.Log("Got  : ", got)
			t.Log("Want : ", syms)
			t.Fail()
		}
	}
}
