// Package ext runs test from outside the packages
package ext

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/xavier268/kompress/klog"
)

func TestBasicRead(t *testing.T) {
	source := "1234567891234567892456789"

	klog := klog.NewKlogReader(strings.NewReader(source))

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
	klog := klog.NewKlogWriter(res)

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
