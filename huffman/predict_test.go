package huffman

import (
	"fmt"
	"testing"
)

func TestPred(t *testing.T) {
	suite := []Symbol{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}
	testPredictn(t, 7, suite)
	testPredictn(t, 6, suite)
	testPredictn(t, 5, suite)
	testPredictn(t, 4, suite)
	testPredictn(t, 3, suite)
	testPredictn(t, 2, suite)
	testPredictn(t, 1, suite)

	suite = []Symbol{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}
	testPredictn(t, 7, suite)
	testPredictn(t, 6, suite)
	testPredictn(t, 5, suite)
	testPredictn(t, 4, suite)
	testPredictn(t, 3, suite)
	testPredictn(t, 2, suite)
	testPredictn(t, 1, suite)

	suite = []Symbol{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2}
	testPredictn(t, 7, suite)
	testPredictn(t, 6, suite)
	testPredictn(t, 5, suite)
	testPredictn(t, 4, suite)
	testPredictn(t, 3, suite)
	testPredictn(t, 2, suite)
	testPredictn(t, 1, suite)

	suite = []Symbol{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	testPredictn(t, 7, suite)
	testPredictn(t, 6, suite)
	testPredictn(t, 5, suite)
	testPredictn(t, 4, suite)
	testPredictn(t, 3, suite)
	testPredictn(t, 2, suite)
	testPredictn(t, 1, suite)

}

func testPredictn(t *testing.T, period byte, source []Symbol) {

	fmt.Println("Testing period=", period)
	deltas := NewSymbolBuffer()
	source2 := NewSymbolBuffer()

	p := newPredict(10, period)
	for _, s := range source {
		d := p.forward(s)
		deltas.WriteSymbol(d)
	}

	fmt.Println("source :", source)
	fmt.Println("deltas :", deltas.buf)

	p = newPredict(10, period)
	for _, ss := range deltas.buf {
		s := p.backward(ss)
		source2.WriteSymbol(s)
	}

	// check identity ...
	if len(source) != len(source2.buf) {
		fmt.Println("source :", source)
		fmt.Println("deltas :", deltas)
		fmt.Println("back  :", source2)
		t.Fail()
	}
	for i := range source {
		if source[i] != source2.buf[i] {
			fmt.Println("source :", source)
			fmt.Println("deltas :", deltas)
			fmt.Println("back  :", source2)
			t.Fail()
		}
	}

}
