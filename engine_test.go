package kompress

import "testing"

func TestBasicEngine(t *testing.T) {

	freq := []int{1, 5, 1, 2, 2, 12, 3}

	e := newEngine(freq, 0)
	// check reasonable balance in tree

	wr := e.root.weight
	w0 := e.root.child0.weight
	w1 := e.root.child1.weight

	if wr != w0+w1 {
		t.Fatal("root weights do not match : ", w0, w1, wr)
	}

	if (w0*w1 != 12*14) || (w0+w1 != 12+14) {
		t.Fatal("weights do not match : ", w0, w1)
	}

	// t.Skip()
	e.Dump()
}
