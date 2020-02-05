package kompress

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// Compiler check
var _ Compresser = new(Kdyndelta)

func TestKdindeltaBasic(t *testing.T) {

	k := NewKdyndelta(1)
	source := "abcabbcabcaabbcababbbababababbababababdc"
	in := strings.NewReader(source)
	out := bytes.NewBuffer(nil)

	k.Compress(in, out)
	fmt.Println("From\t", []byte(source), "\nTo \t", out.Bytes())

	if len(k.buf) != 3 || len(source) != len(out.Bytes()) {
		t.Fatal("inconsistent Kdyndelta compression")
	}

}

func TestKdyndeltaCompressDecompress(t *testing.T) {
	k := NewKdyndelta(4)
	kk := NewKdyndelta(4)
	data := getTestData()

	for i := 0; i < len(data); i++ {
		k.Reset()
		kk.Reset()

		out1 := bytes.NewBuffer(nil)
		k.Compress(bytes.NewReader(data[i]), out1)

		out2 := bytes.NewBuffer(nil)
		kk.Decompress(bytes.NewReader(out1.Bytes()), out2)

		if len(out2.Bytes()) != len(out1.Bytes()) || bytes.Compare(data[i], out2.Bytes()) != 0 {
			fmt.Print(data[i], " ==> ", out1.Bytes())
			fmt.Println(" ==> ", out2.Bytes())

			t.Fatal("invalid Kdyndelta compress/decompress operation")
		}
	}
}
