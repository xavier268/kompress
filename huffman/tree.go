package huffman

import "fmt"

// node makes an huffman tree
type node struct {
	weight uint
	parent *node
	// child for 0 or 1 respectively
	child0, child1 *node
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
	nodes []*node
}

/* func (n *node) String() string {
	return fmt.Sprintf("Node: %4p (weight : %d) \tpar: %4p\t childs: %4p\t%4p ",
		n, n.weight, n.parent, n.child0, n.child1)
} */

func (e *Engine) String() string {
	r := fmt.Sprintf("Number of symbol :\t%d\n", e.len)
	for _, n := range e.nodes {
		r += fmt.Sprintf("\n%p\t%+v", n, n)
	}
	r += fmt.Sprintf("\n\nroot : %p\t%+v", e.root, e.root)
	return r
}

// New creates a new engine, with the initial weights provided.
// There should be exactly or more weights provided than the alphabet size.
// All weight value should be positive or zero.
func New(weights []int) *Engine {
	e := new(Engine)
	e.len = len(weights)
	e.nodes = make([]*node, 2*e.len-1, 2*e.len-1)
	for i := range e.nodes {
		e.nodes[i] = new(node)
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
	// reset existing leaves
	for _, n := range e.nodes[:e.len] {
		n.parent = nil
	}
	// loop until all node capacity have been used
	// alloc points to the next free node.
	for alloc := e.len; alloc < 2*e.len-1; alloc++ {

		// find the lowest weight node among
		// the leaves and those that have been constructed,
		// but that do not have a parent yet
		var n0, n1 *node
		var w0, w1 uint
		for _, n := range e.nodes {
			if n0 == nil {
				n0 = n
				w0 = n.weight
				continue
			}
			if n.parent == nil && n.weight < w0 {
				n0 = n
				w0 = n.weight
			}
		}
		if n0 == nil {
			panic("internal logic error")
		}
		// same for second node, distinct ...
		for _, n := range e.nodes {
			if n1 == nil {
				n1 = n
				w1 = n.weight
				continue
			}
			if n.parent == nil && n.weight < w1 && n != n0 {
				n1 = n
				w1 = n.weight
			}
		}
		if n1 == nil {
			panic("internal logic error")
		}
		// create parent, reading beyond length, but within capacity ...
		p := e.nodes[alloc]
		p.weight = n0.weight + n1.weight
		p.child0 = n0
		p.child1 = n1
		n0.parent, n1.parent = p, p

		// root is last node created
		e.root = p
	}
}
