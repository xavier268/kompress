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
	value  Symbol
	parent *lzwNode
	// what symbol "in" lead us her ?
	pchild Symbol
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

// sequence from any node, upwards until root
func (n *lzwNode) sequence() []Symbol {
	var sq []Symbol

	for nn := n; nn != nil; nn = nn.parent {
		if nn.parent != nil {
			sq = append([]Symbol{nn.pchild}, sq...)
		}
	}

	return sq

}

// dump a (sub) node, using the provided symbol subsequence.
func (n *lzwNode) dump(seq []Symbol) {

	if len(seq) != 0 {
		fmt.Println(n.value, "\t:\t", seq, "should equal \t", n.sequence())
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
		n.parent = l.root
		n.pchild = Symbol(s)
		l.root.childs[Symbol(s)] = n
		l.rev[Symbol(s)] = n
	}
	return l
}

func (lz *lzw) dump() {
	fmt.Println("Number of nodes : ", len(lz.rev))
	for s, n := range lz.rev {
		fmt.Println(s, "\t==>\t", n.sequence())
	}
}

//============ lzwwriter =================================

type lzwwriter struct {
	*lzw
	w SymbolWriteCloser
}

// NewLZWWriter constructor.
func NewLZWWriter(sw SymbolWriteCloser,
	alphaLenIn int, alphaLenOut int, maxSeq int) SymbolWriteCloser {
	return newlzwwriter(sw, alphaLenIn, alphaLenOut, maxSeq)
}

func newlzwwriter(sw SymbolWriteCloser,
	alphaLenIn int, alphaLenOut int, maxSeq int) *lzwwriter {

	l := new(lzwwriter)
	l.lzw = newLzw(alphaLenIn, alphaLenOut, maxSeq)
	l.w = sw
	return l
}

// WriteSymbol incoming symbol s, emitting symbol as needed.
func (lz *lzwwriter) WriteSymbol(s Symbol) error {

	// update sequence length
	lz.seqLen++

	n, ok := lz.current.childs[s]
	switch {
	case ok && lz.seqLen <= lz.seqMax:
		// we can accepts more syms in sequence
		lz.current = n
		return nil
	case ok && lz.seqLen > lz.seqMax:
		// have to stop here, say we found it !
		lz.current = lz.root
		lz.seqLen = 0
		return lz.w.WriteSymbol(n.value)

	case !ok:
		// here, the full seq ending with s does not exist

		// what we had  before adding s is returned
		err := lz.w.WriteSymbol(lz.current.value)
		if lz.seqLen <= lz.seqMax && len(lz.rev) < lz.nbOut {
			// we store the new sequence, adding s
			// provided length is still ok
			nn := newLzwNode()
			nn.value = Symbol(len(lz.rev))
			nn.pchild = s
			lz.rev[nn.value] = nn
			lz.current.childs[s] = nn
			nn.parent = lz.current
		}
		// update pointers and sequence length
		lz.current = lz.root.childs[s]
		lz.seqLen = 1
		return err
	default:
		panic("invalid state")

	}

}

// Close and flush pending sequence.
func (lz *lzwwriter) Close() error {
	if lz.seqLen != 0 {
		err := lz.w.WriteSymbol(lz.current.value)
		if err != nil {
			return err
		}

	}

	return lz.w.Close()
}

//============ lzwreader =================================

type lzwreader struct {
	*lzw
	sr SymbolReader
	// sequence buffer
	seq []Symbol
}

func newlzwreader(sr SymbolReader, alphaLenIn int, alphaLenOut int, maxSeq int) *lzwreader {
	r := new(lzwreader)
	r.lzw = newLzw(alphaLenIn, alphaLenOut, maxSeq)
	r.sr = sr
	return r
}

// ReadSymbol back from compressed reader.
func (lz *lzwreader) ReadSymbol() (s1 Symbol, err error) {

	var ok bool
	var s2 Symbol
	var n *lzwNode

	// a sequence partially available ...
	if len(lz.seq) > 0 {
		s1 = lz.seq[0]
		lz.seq = lz.seq[1:]
		return s1, nil
	}

	// no seq available, get leaf node and sequence
	s2, err = lz.sr.ReadSymbol()
	if err != nil {
		return 0, err
	}
	n, ok = lz.rev[s2]
	if !ok {
		panic("invalid state")
	}
	lz.seq = n.sequence()
	s1 = lz.seq[0]
	lz.seq = lz.seq[1:]

	// try to extend former leaf node, pointed by current,
	// with the first new sequence symbol
	if lz.seqLen <= lz.seqMax && len(lz.rev) < lz.nbOut {
		if lz.current.childs[s1] == nil { // don't overwrite !
			nn := new(lzwNode)
			nn.value = Symbol(len(lz.rev))
			nn.pchild = s1
			lz.rev[nn.value] = nn
			lz.current.childs[s1] = nn
			nn.parent = lz.current
		}
	}

	// remember what former node was ...
	lz.seqLen = len(lz.seq)
	lz.current = n

	return s1, nil

}
