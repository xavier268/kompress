package huffman

import (
	"fmt"
	"testing"
)

func TestDico(t *testing.T) {

	d := newlzwdico(4, 100, 4)
	source := []Symbol{1, 2, 3, 1, 2, 3, 1, 2, 3, 2, 0, 1, 0, 0, 0, 1, 1, 1}

	fmt.Println("Learning all subsequences in ", source)
	d.learn(source...)
	d.dump()
	if len(d.codes) != 37 {
		t.Fatal("unexpected dictionnary length")
	}

}
