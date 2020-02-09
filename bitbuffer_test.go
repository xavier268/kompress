package kompress

import (
	"bytes"
	"fmt"
	"testing"
)

func TestBitBufferBasic(t *testing.T) {

	bb := NewBitBuffer()
	tb := []byte{0, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1}

	for _, b := range tb {
		bb.WriteBit(b)
	}

	var res []byte
	for bb.Size() > 0 {
		bit, err := bb.ReadBit()
		if err != nil {
			t.Fatal(err)
		}
		res = append(res, bit)

	}

	if bytes.Compare(res, tb) != 0 {
		fmt.Println(tb)
		fmt.Println(res)
		t.Fail()

	}
}

func TestReadByteBitBuffer(t *testing.T) {

	bb := NewBitBuffer()
	tb := []byte{0, 1, 1, 0, 0, 1, 0, 1,
		1, 0, 0, 0, 0, 0, 1, 1,
		1}

	for _, b := range tb {
		bb.WriteBit(b)
	}

	b, err := bb.ReadByte()
	if err != nil {
		t.Fatal(err)
	}
	if b != byte(0b01100101) {
		fmt.Printf("%08b\n", b)
		t.Fatal("wrong first byte")
	}

	b, err = bb.ReadByte()
	if err != nil {
		t.Fatal(err)
	}
	if b != byte(0b10000011) {
		fmt.Printf("%08b\n", b)
		t.Fatal("wrong second byte")
	}

	b, err = bb.ReadByte()
	if err == nil {
		t.Fatal("underflow error was expected !")
	}
	if bb.Size() != 1 {
		fmt.Println("Buffer size :", bb.Size())
		t.Fatal("Size should not have changed")
	}

}

func TestWriteByteBitBuffer(t *testing.T) {

	bb := NewBitBuffer()

	bb.WriteBit(1)
	bb.WriteByte(0b10010010)
	bb.WriteBit(1)

	b, err := bb.ReadByte()
	if err != nil {
		t.Fatal(err)
	}
	if b != 0b11001001 {
		fmt.Printf("Read : %08b\n", b)
		t.Fatal("unexepected read byte value")
	}
	if bit, err := bb.ReadBit(); err != nil || bit != 0 {
		t.Fatal("Read bit error", err, bit)
	}
	if bit, err := bb.ReadBit(); err != nil || bit != 1 {
		t.Fatal("Read bit error", err, bit)
	}
	if bb.Size() != 0 {
		t.Fatal("buffer size sould be 0")
	}

	// read past limit
	if _, err := bb.ReadBit(); err == nil {
		t.Fatal("underflow was expected")
	}
	if bb.Size() != 0 {
		t.Fatal("buffer size sould be 0")
	}

}
