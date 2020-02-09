package e2e

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/xavier268/kompress"
)

func TestKDeltaBasicWriter(t *testing.T) {

	source := "abcabcabcdabcdeaaacbaa"
	res := bytes.NewBuffer(nil)

	k := kompress.NewKdeltaWriter(kompress.NewKlogWriter(res), 2)
	n, err := k.Write([]byte(source))
	k.Close()

	if n != len(source) || err != nil {
		t.Fatal(err)
	}
	if len(res.Bytes()) != len(source) {
		fmt.Println("\nSize went from ", len(source), " to ", len(res.Bytes()))
		t.Fatal("not the same byte length")
	}
}
