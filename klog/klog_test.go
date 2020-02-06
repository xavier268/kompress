package klog

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

// A Reader implements the io.Reader interface.
var _ io.Reader = new(Reader)

func TestBasic(t *testing.T) {
	source := "1234567891234567892456789"

	in := NewReader(strings.NewReader(source))
	klog := NewReader(in)

	res, err := ioutil.ReadAll(klog)
	if err != nil {
		t.Fatal(err.Error())
	}

	if bytes.Compare(res, []byte(source)) != 0 {
		t.Fatal("invalid log output")
	}

}
