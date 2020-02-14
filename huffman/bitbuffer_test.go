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
	if err = bb.WriteBits(0, 5); err != nil {
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

	bb.Close()

	_, err = bb.ReadBit()
	if err != io.EOF {
		t.Fatal("exected io.EOF, but got ", err)
	}

}

func TestBitBuffer1Byte(t *testing.T) {
	bb := NewBitBuffer()

	// 1 byte
	bb.WriteByte(0xE3)
	b, err := bb.ReadByte()
	if b != 0xE3 || err != nil {
		t.Fatal("Unexpected bytes or error", b, err)
	}
}

func TestBitBuffer2Bytes(t *testing.T) {
	var b byte
	var err error

	bb := NewBitBuffer()

	// 2 bytes
	bb.WriteByte(0xE3)
	bb.WriteBit(1)
	bb.WriteByte(0xA7)

	b, err = bb.ReadByte()
	if b != 0xE3 || err != nil {
		t.Fatal("Unexpected bytes or error", b, err)
	}
	bi, err := bb.ReadBit()
	if bi != 1 || err != nil {
		t.Fatal("Unexpected bytes or error", b, err)
	}
	b, err = bb.ReadByte()
	if b != 0xA7 || err != nil {
		t.Fatal("Unexpected bytes or error", b, err)
	}

}
