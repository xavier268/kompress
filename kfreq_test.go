package kompress

import (
	"fmt"
	"testing"
)

func TestKfreqBasic(t *testing.T) {

	source := [...]byte{0, 3, 3, 3, 2, 1, 6}

	f := new(kfreq)
	if f.esc != 0 {
		fmt.Print(f.esc)
		t.Fatal("wrong initial esc value")
	}
	for _, b := range source {
		f.update(b)
	}

	// define expecte result
	var r [256]int
	r[0], r[1], r[2], r[3], r[6] = 1, 1, 1, 3, 1

	if r != f.freq {
		t.Fatal("bad freq table")
	}
	if f.esc != 4 {
		fmt.Print(f.esc)
		t.Fatal("wrong esc value")
	}

	// force successive esc value change
	f.update(4)
	if f.esc != 5 {
		fmt.Print(f.esc)
		t.Fatal("wrong  esc value")
	}
	f.update(5)
	if f.esc != 7 {
		fmt.Print(f.esc)
		t.Fatal("wrong  esc value")
	}

}
