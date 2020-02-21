package kompress

import (
	"fmt"
	"io"
	"testing"
)

func TestLZWManualExplore1(t *testing.T) {
	var nbIn, nbOut, seqMax int
	var source []Symbol

	nbIn, nbOut = 3, 30
	seqMax = 3
	source = []Symbol{1, 2, 1, 2, 1, 2, 1, 2}
	testLZWReadWrite(t, nbIn, nbOut, seqMax, source)

	seqMax = 4
	testLZWReadWrite(t, nbIn, nbOut, seqMax, source)

	seqMax = 5
	testLZWReadWrite(t, nbIn, nbOut, seqMax, source)

	seqMax = 5
	source = []Symbol{1, 2, 0, 2, 1, 2, 2, 0, 0, 0, 2, 1, 1, 0, 1, 2, 2, 0, 1, 0, 1}
	testLZWReadWrite(t, nbIn, nbOut, seqMax, source)

	//t.Skip()

	nbIn, nbOut = 256, 100*255
	seqMax = 4
	for _, x := range []byte(text1) {
		source = append(source, Symbol(x))
	}
	testLZWReadWrite(t, nbIn, nbOut, seqMax, source)

	nbIn, nbOut = 256, 100*255
	seqMax = 4
	for _, x := range []byte(text2) {
		source = append(source, Symbol(x))
	}
	testLZWReadWrite(t, nbIn, nbOut, seqMax, source)

}

// ======  utilities back & forth test =============
func testLZWReadWrite(t *testing.T, nbIn, nbOut,
	seqMax int, source []Symbol) {

	bb := NewSymbolBuffer()
	l := newlzwwriter(bb, nbIn, nbOut, seqMax)

	for _, s := range source {
		//fmt.Println("Processing=", s)
		l.WriteSymbol(s)
		//l.dump()
	}
	l.Close()
	//fmt.Println("Closing write ...")
	l.dump()
	if len(source) < 15 {
		fmt.Println("From \t", len(source), source)
	} else {
		fmt.Println("From \t", len(source), source[:14], "...")
	}
	if len(bb.buf) < 15 {
		fmt.Println("To   \t", len(bb.buf), bb.buf)
	} else {
		fmt.Println("To   \t", len(bb.buf), bb.buf[:14], "...")
	}

	for _, s := range bb.buf {
		if int(s) >= nbOut {
			t.Fatal("using symbols that are out of the alphabet (too large)")
		}
	}

	// read back buffer

	r := newlzwreader(bb, nbIn, nbOut, seqMax)
	res := []Symbol{}
	for len(res) <= len(source)+2 {
		s1, err := r.ReadSymbol()
		//r.dump()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		res = append(res, s1)
	}

	if len(res) < 15 {
		fmt.Println("Back    ", len(res), res)
	} else {
		fmt.Println("Back    ", len(res), res[:14])
	}
	fmt.Println()

	// r.dump()
	if len(res) != len(source) {
		r.dump()
		fmt.Println("Source : ", source)
		fmt.Println("Back   : ", res)
		t.Fatal("source and back do not match !")
	}

	for i := range source {
		if source[i] != res[i] {
			r.dump()
			fmt.Println("Source : ", source)
			fmt.Println("Back   : ", res)
			t.Fatal("source and back do not match !")
		}
	}

}
