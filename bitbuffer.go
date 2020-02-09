package kompress

import (
	"errors"
)

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

// WriteBit writes a bit if b is not nul.
// Write 0 if b == 0.
func (bb *BitBuffer) WriteBit(b byte) {

	bb.size++
	// increase storage capacity if needed
	if bb.offset+bb.size >= 8*len(bb.bytes) {
		bb.bytes = append(bb.bytes, 0)
	}
	n := ((bb.offset + bb.size + 16) % 8)
	k := (bb.offset + bb.size) / 8
	//fmt.Println(k, n)
	if b != 0 {
		bb.bytes[k] |= byte(1 << n)
	} else {
		bb.bytes[k] &= (byte(1<<n) ^ byte(255))
	}
}

// ReadBit from buffer, in fifo order.
func (bb *BitBuffer) ReadBit() byte {

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

// ReadByte reads a single byte from buffer.
// Error if cannot read at least 8 bits.
func (bb *BitBuffer) ReadByte() (byte, error) {
	var b byte
	if bb.Size() < 8 {
		return 0, errors.New("bit buffer underflow")
	}
	for i := 128; i > 0; i >>= 1 {
		// DEBUG
		// fmt.Printf("%08b\n", i)

		if bb.ReadBit() == 1 {
			b |= byte(i)
		}
	}
	return b, nil
}

// WriteByte write the corresponding bits in the buffer.
func (bb *BitBuffer) WriteByte(b byte) error {
	for i := 128; i > 0; i >>= 1 {
		bb.WriteBit(b & byte(i))
	}
	return nil
}
