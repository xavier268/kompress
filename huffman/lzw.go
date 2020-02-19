package huffman

import "fmt"

// LZW algorith implementation.
//
// Principles for encoding :
// ========================
//
// 1. Dictionnary contains sequences <-> key; are multi symbol
// 2. Find longest dictionnary string that match current input
// 3. Replace output with key for that sequence
// 4. Add to dictionnary the sequence formed with the matching sequence
//    followed by the last symbol, until max table size is reached.
//
// Encoding changes the encoding alphabet, from Symbol1 to Symbol2
// Symbol1 is the original Symbol,
// Symbol2 has up to dicomax values,

// maxLen is the max length of the sequence to be encoded.
const maxLen = 20

// ============ structure for the sequence distionnary is a tree ========

type lzwNode struct {
	childs map[Symbol]*lzwNode
	// encoding Symbol "out" of the sequence at this point
	// only valid if a leaf node (ie, no child)
	value Symbol
}

func newLzwNode() *lzwNode {
	n := new(lzwNode)
	n.childs = make(map[Symbol]*lzwNode, 0)
	return n
}

// is the node a leaf ?
func (n *lzwNode) isLeaf() bool {
	return len(n.childs) == 0
}

// dump a (sub) node, using the provided symbol subsequence.
func (n *lzwNode) dump(seq []Symbol) {

	if len(seq) != 0 {
		fmt.Println(n.value, "\t:\t", seq)
	}

	for s, nn := range n.childs {
		nn.dump(append(seq, Symbol(s)))
	}
}

// ============  lzw structure itself ====================================

type lzw struct {
	// root of the dictionnary tree for sequences.
	// It maps sequences to a new "out" Symbol.
	root *lzwNode
	// reverse maping, from Symbol "Out" to the node containing that value
	rev map[Symbol]*lzwNode
	// points to the current sub tree, given sequence seen so far.
	current *lzwNode
	// current seq length so far
	seqLen int
	// max sequence length
	seqMax int
	// alphabet length in and out
	nbIn, nbOut int
}

// newLZW constructs a LZW engine,
// from one alphabet size to another.
func newLzw(alphaLenIn int, alphaLenOut int, maxSeq int) *lzw {
	if alphaLenOut <= alphaLenIn {
		panic(fmt.Sprint("incompatible input parameters :",
			alphaLenIn, alphaLenOut, maxSeq))
	}
	l := new(lzw)
	l.nbIn = alphaLenIn
	l.nbOut = alphaLenOut
	l.seqMax = maxSeq
	l.root = newLzwNode()
	l.current = l.root
	l.rev = make(map[Symbol]*lzwNode)

	// init the tree, creating child nodes for each Symbol "in"
	for s := 0; s < l.nbIn; s++ {
		n := newLzwNode()
		n.value = Symbol(s)
		l.root.childs[Symbol(s)] = n
		l.rev[Symbol(s)] = n
	}
	return l
}

// forward incoming symbol s, emitting symbol as needed.
func (lz *lzw) forward(s Symbol, emit func(Symbol)) {

	// update sequence length
	lz.seqLen++

	n, ok := lz.current.childs[s]
	switch {
	case ok && lz.seqLen <= lz.seqMax:
		// we can accepts more syms in sequence
		lz.current = n
		return
	case ok && lz.seqLen > lz.seqMax:
		// have to stop here, say we found it !
		lz.current = lz.root
		lz.seqLen = 0
		emit(n.value)
		return
	case !ok:
		// here, the full seq ending with s does not exist

		// save what we had  before adding s, that will be returned
		emit(lz.current.value)
		if lz.seqLen <= lz.seqMax && len(lz.rev) < lz.nbOut {
			// we store the new sequence, adding s
			// provided length is still ok
			nn := newLzwNode()
			nn.value = Symbol(len(lz.rev))
			lz.rev[nn.value] = nn
			lz.current.childs[s] = nn
		}
		// update pointers and sequence length
		lz.current = lz.root.childs[s]
		lz.seqLen = 1
		return
	default:
		panic("invalid state")

	}

}
