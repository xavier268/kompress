package huffman

import (
	"fmt"
	"io"
	"testing"
)

func TestDynWriteManualExploration(t *testing.T) {

	// uncomment for manual exploration
	t.Skip()

	bb := NewBitBuffer()
	freq := []int{2, 2, 2, 2, 2, 2}
	syms := []Symbol{0, 1, 2, 0, 1, 0, 1, 0, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}

	dw := newDWriter(bb, 7, freq, SCAlways)
	fmt.Println("Initial dynwriter : ")
	dw.Dump()
	for _, s := range syms {
		fmt.Println("Writing to dynWriter : ", s)
		err := dw.WriteSymbol(s)
		dw.Dump()
		if err != nil {
			panic(err)
		}
	}

}

func TestDynReadWrite1(t *testing.T) {

	freq := []int{2, 2, 2, 2, 2, 2}
	syms := []Symbol{0, 1, 2, 0, 1, 0, 1, 0, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}
	eof := Symbol(7)
	dynreadwrite(t, syms, eof, freq)
}
func TestDynReadWrite2(t *testing.T) {
	freq := []int{2, 2, 222, 2, 2, 2}
	syms := []Symbol{2}
	eof := Symbol(0)
	dynreadwrite(t, syms, eof, freq)
}
func TestDynReadWritePanicExpected(t *testing.T) {
	defer expectPanic(t)
	freq := []int{2, 2, 2, 2, 2, 2}
	syms := []Symbol{3, 5}
	eof := Symbol(5)
	dynreadwrite(t, syms, eof, freq)
}

func dynreadwrite(t *testing.T, syms []Symbol, eof Symbol, freq []int) {

	bb := NewBitBuffer()
	dw := newDWriter(bb, eof, freq, SCAlways)

	// write all
	for _, s := range syms {
		if err := dw.WriteSymbol(s); err != nil {
			t.Fatal("Error while writing : ", err)
		}
	}
	// close and flush ...
	dw.Close()

	t.Log("DynReadWrite compressed from ", len(syms)*8, " bits, to ", bb.Size(), " bits")

	dr := newDReader(bb, eof, freq, SCAlways)

	for _, s := range syms {
		ss, err := dr.ReadSymbol()
		if err != nil {
			t.Log("Read error :", ss, s, err)
			panic("Unexpected read error")
		}
		if s != ss {
			t.Fatal("Got ", ss, " instead of ", s)
		}
	}

	_, err := dr.ReadSymbol()
	if err != io.EOF {
		t.Fatal("Expected io.EOF, got ", err)
	}
}
