package huffman

import (
	"fmt"
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
	fmt.Println("Size : ", bb.Size())

}
