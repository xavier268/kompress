package huffman

import (
	"fmt"
	"testing"
)

func TestTree(t *testing.T) {

	e := New([]int{1, 20, 3, 7, 1, 5, 1})
	fmt.Println(e)

	s := Symbol(0)
	fmt.Println("encoding ", s)
	e.nodes[s].encode(nil)

	s++
	fmt.Println("encoding ", s)
	e.nodes[s].encode(nil)

	s++
	fmt.Println("encoding ", s)
	e.nodes[s].encode(nil)

	s++
	fmt.Println("encoding ", s)
	e.nodes[s].encode(nil)

	s++
	fmt.Println("encoding ", s)
	e.nodes[s].encode(nil)

	s++
	fmt.Println("encoding ", s)
	e.nodes[s].encode(nil)

	s++
	fmt.Println("encoding ", s)
	e.nodes[s].encode(nil)

	s++
	fmt.Println("encoding ", s)
	e.nodes[s].encode(nil)

	bits := []Bit{0}
	s, err := e.decode(e.root, bits...)
	fmt.Println(bits, " --> ", s, err)

	bits = []Bit{1}
	s, err = e.decode(e.root, bits...)
	fmt.Println(bits, " --> ", s, err)

	bits = []Bit{0, 1, 1, 0}
	s, err = e.decode(e.root, bits...)
	fmt.Println(bits, " --> ", s, err)

	bits = []Bit{0, 1, 1}
	s, err = e.decode(e.root, bits...)
	fmt.Println(bits, " --> ", s, err)

	bits = []Bit{0, 1}
	s, err = e.decode(e.root, bits...)
	fmt.Println(bits, " --> ", s, err)

}
