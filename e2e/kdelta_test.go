package e2e

import (
	"bytes"
	"testing"

	"github.com/xavier268/kompress"
)

func TestKDeltaBasicWriter(t *testing.T) {

	source := "abcabcabcdabcdeaaacbaa"
	res := bytes.NewBuffer(nil)

	k := kompress.NewKdeltaWriter(kompress.NewKrlenWriter(kompress.NewKlogWriter(res)), 2)
	n, err := k.Write([]byte(source))
	k.Close()

	if n != len(source) || err != nil {
		t.Fatal(err)
	}
	if len(res.Bytes()) != len(source) {
		t.Fail()
	}
}
