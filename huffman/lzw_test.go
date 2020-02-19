package huffman

import (
	"fmt"
	"testing"
)

func TestLZWManualExplore(t *testing.T) {

	l := newLzw(3, 20, 3)

	source := []Symbol{1, 2, 1, 2, 1, 2, 1, 2}
	for _, s := range source {

		fmt.Println("Processing=", s)
		l.forward(s, func(s2 Symbol) { fmt.Println("Emitting : ", s2) })
		fmt.Println("Nb of nodes : ", len(l.rev))
		l.root.dump(nil)

	}
}
