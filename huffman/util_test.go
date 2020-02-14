package huffman

import (
	"bytes"
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

func TestBitFromByteReader(t *testing.T) {

	buf := []byte{0b1100_0010}
	ior := bytes.NewBuffer(buf)
	r := NewBitFromByteReader(ior)

	b, err := r.ReadBit()
	if b != 1 || err != nil {
		t.Fatal("wrong bit or error ", b, err)
	}
	b, err = r.ReadBit()
	if b != 1 || err != nil {
		t.Fatal("wrong bit or error ", b, err)
	}
	b, err = r.ReadBit()
	if b != 0 || err != nil {
		t.Fatal("wrong bit or error ", b, err)
	}
	b, err = r.ReadBit()
	if b != 0 || err != nil {
		t.Fatal("wrong bit or error ", b, err)
	}

	b, err = r.ReadBit()
	if b != 0 || err != nil {
		t.Fatal("wrong bit or error ", b, err)
	}
	b, err = r.ReadBit()
	if b != 0 || err != nil {
		t.Fatal("wrong bit or error ", b, err)
	}
	b, err = r.ReadBit()
	if b != 1 || err != nil {
		t.Fatal("wrong bit or error ", b, err)
	}
	b, err = r.ReadBit()
	if b != 0 || err != nil {
		t.Fatal("wrong bit or error ", b, err)
	}

	// expect  error ...
	b, err = r.ReadBit()
	if b != 0 || err != io.EOF {
		t.Fatal("wrong bit", b)
	}
}

func TestBitToByteWriter(t *testing.T) {

	var source []Bit
	var is, should []byte

	source = []Bit{1, 0, 0, 1,
		1, 1, 0, 1,
		0, 1, 1, 0,
		1, 1}
	should = []byte{0b_1001_1101, 0b_0110_1100}

	buf := bytes.NewBuffer(is)

	w := NewBitToByteWriter(buf)
	err := w.WriteBits(source...)

	// Before close ...
	if bytes.Compare(should[:1], buf.Bytes()) != 0 || err != nil {
		fmt.Println("\nsource:", source)
		fmt.Println("got before close:", buf.Bytes())
		fmt.Println("want:", should)
		fmt.Println("error:", err)
		t.Fatal(err)
	}

	// close
	w.Close()
	if bytes.Compare(should, buf.Bytes()) != 0 || err != nil {
		fmt.Println("\nsource:", source)
		fmt.Println("got:", buf.Bytes())
		fmt.Println("want:", should)
		fmt.Println("error:", err)
		t.Fatal(err)
	}

}
