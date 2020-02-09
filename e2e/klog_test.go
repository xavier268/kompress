// Package ext runs test from outside the packages
package e2e

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/xavier268/kompress"
)

func TestBasicRead(t *testing.T) {
	source := "1234567891234567892456789"

	klog := kompress.NewKlogReader(strings.NewReader(source))

	res, err := ioutil.ReadAll(klog)
	if err != nil {
		t.Fatal(err.Error())
	}

	if bytes.Compare(res, []byte(source)) != 0 {
		t.Fatal("invalid log output")
	}

}

func TestBasicWrite(t *testing.T) {

	source := "1234567891234567892456789"

	res := bytes.NewBuffer(nil)
	klog := kompress.NewKlogWriter(res)

	n, err := klog.Write([]byte(source))
	if err != nil {
		t.Fatal(err.Error())
	}
	if n != len(source) {
		t.Fatal("length error")
	}
	if bytes.Compare(res.Bytes(), []byte(source)) != 0 {
		t.Fatal("invalid log output")
	}
}
