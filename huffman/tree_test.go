package huffman

import "testing"

func TestBasicTreeDump(t *testing.T) {
	freq := []int{1, 5, 0, 2, 2, 12, 3}
	e := newEngine(freq, 0)
	e.Dump()
}
