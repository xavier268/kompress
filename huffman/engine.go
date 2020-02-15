package huffman

import (
	"fmt"
)

// node makes an huffman tree
type node struct {
	id     int
	weight uint
	parent *node
	// child for 0 or 1 respectively
	child0, child1 *node
}

func (n *node) String() string {
	r := fmt.Sprintf("Node %d, weight %d \t", n.id, n.weight)
	if n.parent == nil {
		r += "parent :  \t"
	} else {
		r += fmt.Sprintf("parent : %d\t", n.parent.id)
	}
	if n.child0 == nil {
		r += "child0 :  ,\t"
	} else {
		r += fmt.Sprintf("child0 : %d,\t", n.child0.id)
	}
	if n.child1 == nil {
		r += "child1 :  "
	} else {
		r += fmt.Sprintf("child1 : %d", n.child1.id)
	}
	return r
}

// engine to handle huffman trees.
type engine struct {
	// Symbol alphabet size, number of leaves
	len int
	// freq are the frequencies that were used to construct the tree
	// We need to keep them to reconstruct the initial tree (eg, dump)
	freq []int
	// actual frequencies, are updated everytime a symbol is written or read.
	actfreq []int
	// eof is an application level eof Symbol, part of the Alphabet
	eof Symbol
	// root node
	root *node
	// nodes gather all nodes,
	// fisrt, leaf nodes, one by Symbol,
	// then the rest of the nodes, including the root, as the last one.
	nodes []node
}

// newEngine creates a new engine, with the initial weights provided.
// There should be exactly or more weights provided than the alphabet size.
// All weight value should be positive or zero.
func newEngine(weights []int, eof Symbol) *engine {
	e := new(engine)
	e.freq = weights
	e.actfreq = append([]int(nil), weights...) // deep copy ...
	e.eof = eof
	e.len = len(weights)
	e.nodes = make([]node, 2*e.len-1, 2*e.len-1)
	for i := range e.nodes {
		e.nodes[i].id = i
		if i < e.len {
			e.nodes[i].weight = uint(e.freq[i])
		}
	}
	e.makeTree()
	return e
}

// makeTree computes the huffman tree.
// No allocation is made, nodes are reused.
// Leaf weights are unchanged.
func (e *engine) makeTree() {
	e.root = nil
	// reset all parent to zero, weights to initial values ..
	for i, n := range e.nodes {
		n.parent = nil
		if i < e.len {
			n.weight = uint(e.freq[i])
		}
	}
	// loop until all node capacity have been used
	// alloc points to the next free node.
	for alloc := e.len; alloc < 2*e.len-1; alloc++ {

		// find the lowest weight node among
		// the leaves and those that have been constructed,
		// but that do not have a parent yet
		i0, i1 := -1, -1
		var w0, w1 uint
		for i := 0; i < alloc; i++ {
			if e.nodes[i].parent == nil && i0 < 0 {
				i0 = i
				w0 = e.nodes[i].weight
			}
			if e.nodes[i].parent == nil && i0 >= 0 && e.nodes[i].weight < w0 {
				i0 = i
				w0 = e.nodes[i].weight
			}
		}
		if i0 < 0 {
			panic("internal logic error")
		}
		// same for second node, distinct ...
		for i := 0; i < alloc; i++ {
			if e.nodes[i].parent == nil && i1 < 0 && i != i0 {
				i1 = i
				w1 = e.nodes[i].weight
			}
			if e.nodes[i].parent == nil && e.nodes[i].weight < w1 && i != i0 {
				i1 = i
				w1 = e.nodes[i].weight
			}
		}
		if i1 < 0 {
			panic("internal logic error")
		}
		// create parent, reading beyond length, but within capacity ...
		// fmt.Println("i0, i1 = ", i0, " , ", i1)

		e.nodes[alloc].weight = e.nodes[i0].weight + e.nodes[i1].weight
		e.nodes[alloc].child0 = &e.nodes[i0]
		e.nodes[alloc].child1 = &e.nodes[i1]
		e.nodes[i0].parent, e.nodes[i1].parent = &e.nodes[alloc], &e.nodes[alloc]

		//fmt.Println("Using    : ", e.nodes[i0].String())
		//fmt.Println("Using    : ", e.nodes[i1].String())
		//fmt.Println("Using (p): ", e.nodes[alloc].String())

		// root is last node created
		e.root = &e.nodes[alloc]
	}
}

func (e *engine) Dump() {
	fmt.Println("Dumping huffman tree ( ", e.len, " symbols, eof is ", e.eof, " )")

	lw := new(LogWriter)
	w := NewWriter(lw, 0, e.freq)

	for i, f := range e.freq {
		s := Symbol(i)
		lw.WriteSymbol(s)
		fmt.Printf("/weight %d/actualfreq %d => ", f, e.actfreq[i])
		err := w.WriteSymbol(s)
		if err != nil {
			panic("unexpected write error ")
		}
		fmt.Println()
	}
}
