package kompress

import (
	"testing"
)

func TestWriter(t *testing.T) {

	var err error
	freq := []int{1, 2, 0, 9, 1}
	bb := NewBitBuffer()

	w := newWriter(bb, 0, freq)
	err = w.WriteSymbol(2)

	if err != nil {
		w.Dump()
		t.Fatal("error on write :", err)
	}

	if w.engine.actfreq[2] != freq[2]+1 {
		t.Fatal("weights did not update on write !")
	}

	l := bb.Size()
	// write invalid symbol
	err = w.WriteSymbol(22)
	if err == nil {
		t.Fatal("expected error, but got nil")
	}
	if l != bb.Size() {
		t.Fatal("no bit should have been emitted for an invalid symbol")
	}

}
