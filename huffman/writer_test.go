package huffman

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
		t.Fatal(err)
	}

}
