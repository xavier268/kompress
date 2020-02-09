package kompress

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

// A Reader implements the io.Reader interface.
var _ io.Reader = new(klogReader)

func TestBasicRead(t *testing.T) {
	source := "1234567891234567892456789"

	klog := NewKlogReader(strings.NewReader(source))

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
	klog := NewKlogWriter(res)

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
