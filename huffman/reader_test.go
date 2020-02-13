package huffman

import (
	"fmt"
	"testing"
)

func TestReader1(t *testing.T) {

	bits := []Bit{1, 0, 0, 0, 1, 2, 0, 0, 0}
	bb := NewBitBuffer(bits...)

	freq := []int{0, 0, 4, 2, 2, 8, 1, 9}

	r := newReader(bb, 0, freq)
	r.Dump()

	for bb.Size() > 0 && r.err == nil {
		fmt.Println("Reading from : ", bb.bits)
		s, err := r.ReadSymbol()
		fmt.Println("Read symbol : ", s, ", error : ", err)
	}
}

func TestReader2(t *testing.T) {

	bits := []Bit{1, 0, 0, 0, 1, 2, 0, 0}
	bb := NewBitBuffer(bits...)

	freq := []int{0, 0, 4, 2, 2, 8, 1, 9}

	r := newReader(bb, 0, freq)
	r.Dump()

	for bb.Size() > 0 && r.err == nil {
		fmt.Println("Reading from : ", bb.bits)
		s, err := r.ReadSymbol()
		fmt.Println("Read symbol : ", s, ", error : ", err)
	}
}

func TestReader3LogicEOF(t *testing.T) {

	bits := []Bit{1, 0, 0, 0, 1, 2, 0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	bb := NewBitBuffer(bits...)

	freq := []int{0, 0, 4, 2, 2, 8, 1, 9}

	r := newReader(bb, 0, freq)
	r.Dump()

	for bb.Size() > 0 && r.err == nil {

		fmt.Println("Reading from : ", bb.bits)
		s, err := r.ReadSymbol()
		fmt.Println("Read symbol : ", s, ", error : ", err)
	}
}
