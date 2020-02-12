package huffman

import "fmt"

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

// Engine to handle huffman trees.
type Engine struct {
	// Symbol alphabet size, number of leaves
	len int
	// root node
	root *node
	// nodes gather all nodes,
	// fisrt, leaf nodes, one by Symbol,
	// then the rest of the nodes, including the root, as the last one.
	nodes []node
}

func (e *Engine) String() string {
	r := fmt.Sprintf("Number of symbol :\t%d\n", e.len)
	for n := range e.nodes {
		r += fmt.Sprintln(e.nodes[n].String())
	}
	r += fmt.Sprintln("root:\n", e.root.String())
	return r
}

// New creates a new engine, with the initial weights provided.
// There should be exactly or more weights provided than the alphabet size.
// All weight value should be positive or zero.
func New(weights []int) *Engine {
	e := new(Engine)
	e.len = len(weights)
	e.nodes = make([]node, 2*e.len-1, 2*e.len-1)
	for i := range e.nodes {
		e.nodes[i].id = i
		if i < e.len {
			e.nodes[i].weight = uint(weights[i])
		}
	}
	e.makeTree()
	return e
}

// makeTree computes the huffman tree.
// No allocation is made, nodes are reused.
// Leaf weights are unchanged.
func (e *Engine) makeTree() {
	e.root = nil
	// reset all parent to zero
	for _, n := range e.nodes {
		n.parent = nil
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
		fmt.Println("i0, i1 = ", i0, " , ", i1)

		e.nodes[alloc].weight = e.nodes[i0].weight + e.nodes[i1].weight
		e.nodes[alloc].child0 = &e.nodes[i0]
		e.nodes[alloc].child1 = &e.nodes[i1]
		e.nodes[i0].parent, e.nodes[i1].parent = &e.nodes[alloc], &e.nodes[alloc]

		fmt.Println("Using    : ", e.nodes[i0].String())
		fmt.Println("Using    : ", e.nodes[i1].String())
		fmt.Println("Using (p): ", e.nodes[alloc].String())

		// root is last node created
		e.root = &e.nodes[alloc]
	}
}
