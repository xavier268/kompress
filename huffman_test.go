package kompress

import "testing"

func TestConstructHuffman(t *testing.T) {

	var freq [256]int
	freq[33] = 4
	freq[35] = 7
	freq[20] = 12

	_, a := newTree(freq)

	a.Dump()
	t.Fail()
}
