package huffman

import (
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
	if err = bb.WriteBit(1); err != nil {
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
