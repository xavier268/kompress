package test

import (
	"bytes"
	"testing"

	"github.com/xavier268/kompress/klog"
)

func TestKDeltaBasicWriter(t *testing.T) {

	source := "abcabcabcdabcdeaaacbaa"
	res := bytes.NewBuffer(nil)

	k := klog.NewKdeltaWriter(klog.NewKrlenWriter(klog.NewKlogWriter(res)), 2)
	n, err := k.Write([]byte(source))
	k.Close()

	if n != len(source) || err != nil {
		t.Fatal(err)
	}
	if len(res.Bytes()) != len(source) {
		t.Fail()
	}
}
