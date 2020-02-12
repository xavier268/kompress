package kompress

import "fmt"

type node struct {
	value  byte // value, if a leaf
	signal byte // special leaf for eof marking or other signalling
	leaf   bool // are we a leaf ?
	left   *node
	right  *node
	parent *node
	weight int // cumulative weight of the subtree below.
}

// Signals
const (
	VALUE byte = 0 //  The value represent a byte value
	EOF   byte = 1 // Signal logical EOF. Value has no meaning.
)

// points to the leafs corresponding to each byte.
type alphabet [256]*node

// newTree creates a new Huffman Tree from the byte frequency table provided.
func newTree(freq [256]int) (*node, *alphabet) {

	// All nodes to process, includes signals with 0 frequency.
	nodes := make([]*node, 257, 257)

	// alphabet fo the leaves
	alpha := new(alphabet)

	// create a node per byte
	for i := 0; i < 256; i++ {
		n := new(node)
		n.value = byte(i)
		n.leaf = true
		n.weight = freq[i]
		nodes[i] = n
		alpha[i] = n
	}

	// Add the eof signal nodes
	eof := new(node)
	eof.leaf = true
	eof.signal = EOF
	nodes[256] = eof

	// Loop until we only have one node left (the root)
	for len(nodes) > 1 {

		// find the smallest weight node
		w1, w2 := nodes[0].weight, nodes[1].weight
		s1, s2 := 0, 1
		// find best
		for i, n := range nodes {
			if i != s2 && n.weight < w1 {
				s1 = i
				w1 = n.weight
			}
		}
		// find second best
		for i, n := range nodes {
			if i != s1 && n.weight < w2 {
				s2 = i
				w2 = n.weight
			}
		}
		// debug
		if s1 == s2 {
			panic("internal logic error")
		}

		// create parent node
		par := new(node)
		par.weight = w1 + w2
		par.leaf = false
		par.left = nodes[s1]
		par.right = nodes[s2]
		nodes[s1].parent = par
		nodes[s2].parent = par

		// update nodes, making sure we don't grow the slice !
		nodes[s2] = nodes[len(nodes)-1] // overwite s1 with last node
		nodes[s1] = par                 // overwite s2
		nodes = nodes[:len(nodes)-1]    // truncate
	}

	return nodes[0], alpha
}

func (n *node) String() string {
	if n.leaf {
		return fmt.Sprintf("Leaf %p \tvalue = %d, signal = %d, parent = %p", n, n.value, n.signal, n.parent)
	}

	return fmt.Sprintf("Core %p \tleft = %p, right = %p, parent = %p", n, n.left, n.right, n.parent)

}

// Dump alphabet
func (a *alphabet) Dump() {
	for _, n := range a {
		fmt.Println(n.String())
	}
}
