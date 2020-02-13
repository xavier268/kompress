package huffman

import (
	"fmt"
	"io"
	"testing"
)

func TestBitBuffer(t *testing.T) {

	var err error

	bb := new(BitBuffer)

	if err = bb.WriteBit(1); err != nil {
		t.Fatal(err)
	}
	if err = bb.WriteBit(0); err != nil {
		t.Fatal(err)
	}
	if err = bb.WriteBit(5); err != nil {
		t.Fatal(err)
	}
	if bb.Size() != 3 {
		t.Fatal("wrong size : ", bb.Size())
	}

	b, err := bb.ReadBit()
	if b != 1 || err != nil {
		t.Fatal("invalid : ", b, err)
	}
	b, err = bb.ReadBit()
	if b != 0 || err != nil {
		t.Fatal("invalid : ", b, err)
	}
	b, err = bb.ReadBit()
	if b != 1 || err != nil {
		t.Fatal("invalid : ", b, err)
	}
	if bb.Size() != 0 {
		t.Fatal("wrong size : ", bb.Size())
	}

	_, err = bb.ReadBit()
	if err != io.EOF {
		t.Fatal("exected io.EOF, but got ", err)
	}

}

func TestBitBufferBytes1(t *testing.T) {

	bb := NewBitBuffer()
	s := byte(0b_0101_1001)
	bb.WriteByte(s)
	// fmt.Println(bb)
	ss, err := bb.ReadByte()
	if ss != s || err != nil {
		fmt.Printf("From %08b   ==>   %08b\n", s, ss)
		t.Fatal(err, bb)
	}
}
func TestBitBufferBytes2Padding(t *testing.T) {

	bb := NewBitBuffer()
	s := byte(0b_0101_1001)
	bb.WriteByte(s)
	// fmt.Println(bb)
	bb.ReadBit()
	ss, err := bb.ReadByte()
	if ss != s*2 || err != io.EOF {
		fmt.Printf("From %08b   ==>   %08b\n", s, ss)
		t.Fatal(err, bb)
	}
}
