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
		res = append(res, byte(bb.ReadBit()))
	}

	if bytes.Compare(res, tb) != 0 {
		fmt.Println(tb)
		fmt.Println(res)
		t.Fail()

	}
}

func TestReadByteBitBuffer(t *testing.T) {

	bb := NewBitBuffer()
	tb := []byte{0, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1}

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

}
