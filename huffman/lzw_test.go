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
		//l.root.dump(nil)
		l.dump()
	}
	l.Close()
	//l.root.dump(nil)
	l.dump()

	fmt.Println("From \t", source)
	fmt.Println("To   \t", bb.buf)

	// read back buffer
	var err error
	var s1 Symbol
	r := newlzwreader(bb, 3, 20, 3)

	res := []Symbol{}
	for {
		s1, err = r.ReadSymbol()
		res = append(res, s1)
		if err != nil {
			break
		}
	}
	fmt.Println("Back :", res)
	if err != io.EOF {
		t.Fatal("unexpected error : ", err)
	}
}
