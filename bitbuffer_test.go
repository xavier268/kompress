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
		bb.WriteBit(Bit(b))
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
