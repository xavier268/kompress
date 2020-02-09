package kompress

// Bit are 0 or 1
type Bit byte

// BitBuffer is a buffer to store bits, not bytes.
type BitBuffer struct {
	bytes  []byte // compressed content
	offset int    // for alignement purposes
	size   int    // actual useful content size
}

// NewBitBuffer creates a new empty buffer.
func NewBitBuffer() *BitBuffer {
	b := new(BitBuffer)
	b.bytes = make([]byte, 0, 8)
	return b
}

// Size is the actual useful content size.
func (bb *BitBuffer) Size() int {
	return bb.size
}

// Dump the useful content of buffer
func (bb *BitBuffer) Dump() {

}

// WriteBit writes the following bits in the buffer,
// Panic for value not 0 or 1
func (bb *BitBuffer) WriteBit(b Bit) {

	if b > 1 {
		panic("invalid bit value")
	}
	bb.size++
	// increase storage capacity if needed
	if bb.offset+bb.size >= 8*len(bb.bytes) {
		bb.bytes = append(bb.bytes, 0)
	}
	n := ((bb.offset + bb.size + 16) % 8)
	k := (bb.offset + bb.size) / 8
	//fmt.Println(k, n)
	if b == 1 {
		bb.bytes[k] |= byte(1 << n)
	} else {
		bb.bytes[k] &= (byte(1<<n) ^ byte(255))
	}
}

// ReadBit from buffer, in fifo order.
func (bb *BitBuffer) ReadBit() Bit {

	bb.size--
	bb.offset++

	n := ((bb.offset + 16) % 8)
	k := bb.offset / 8
	//fmt.Println(bb.offset, bb.size, " => ", k, n)
	b := bb.bytes[k] & byte(1<<n)

	// cleanup
	if bb.offset > 8 {
		bb.bytes = bb.bytes[1:]
		bb.offset = bb.offset - 8
	}

	if b == 0 {
		return 0
	}
	return 1

}
