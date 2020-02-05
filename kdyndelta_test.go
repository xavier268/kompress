package kompress

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// Compiler check
var _ Compresser = new(Kdyndelta)

func TestKDinDeltaBasic(t *testing.T) {

	k := NewKdyndelta(4)
	in := strings.NewReader("abcabbccabcaabbcc")
	k.Compress(in, new(Writer))

	if len(k.buf) != 4 {
		t.Fail()
	}

}

func TestKdynCompressDecompress(t *testing.T) {
	k := NewKdyndelta(4)
	kk := NewKdyndelta(4)
	data := getTestData()

	for i := 0; i < len(data) && i < 60; i++ {
		k.Reset()
		kk.Reset()

		out := bytes.NewBuffer(nil)
		k.Compress(bytes.NewReader(data[i]), out)
		fmt.Println(data[i], " ==> ", out.Bytes())

	}
}
