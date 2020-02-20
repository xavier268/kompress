package huffman

import (
	"fmt"
	"io"
	"testing"
)

func TestLZWManualExplore1(t *testing.T) {

	bb := NewSymbolBuffer()

	l := newlzwwriter(bb, 3, 20, 3)

	source := []Symbol{1, 2, 1, 2, 1, 2, 1, 2}

	for _, s := range source {
		fmt.Println("Processing=", s)
		l.WriteSymbol(s)
		fmt.Println("Nb of nodes : ", len(l.rev))
		l.dump()
	}
	l.Close()
	l.dump()

	fmt.Println("From \t", source)
	fmt.Println("To   \t", bb.buf)

	// read back buffer
	var err error
	var s1 Symbol
	r := newlzwreader(bb, 3, 20, 3)

	res := []Symbol{}
	for {
		r.dump()
		s1, err = r.ReadSymbol()
		if err != nil {
			break
		}
		res = append(res, s1)
		fmt.Println("Read : ", res)

	}
	fmt.Println("Back :", res)
	if err != io.EOF {
		t.Fatal("unexpected error : ", err)
	}

	if len(source) != len(res) {
		fmt.Println("Source :", source)
		fmt.Println("Res    :", res)
		t.Fatal("Retuned data does not match source data length")
	}

	for i := range source {
		if source[i] != res[i] {
			fmt.Println("Source :", source)
			fmt.Println("Res    :", res)
			t.Fatal("source and res do not match")
		}
	}
}
