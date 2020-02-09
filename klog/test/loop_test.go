package test

import (
	"bytes"
	"io"
	"math/rand"
	"testing"

	"github.com/xavier268/kompress/klog"
)

func TestCompressDecompressKDelta(t *testing.T) {

	res := bytes.NewBuffer(nil)

	w := klog.NewKdeltaWriter(res, 4)
	r := klog.NewKdeltaReader(res, 4)

	compressDecompress(t, w, r)

}

func TestCompressDecompressKrlen(t *testing.T) {

	res := bytes.NewBuffer(nil)

	w := klog.NewKrlenWriter(res)
	r := klog.NewKrlenReader(res)

	compressDecompress(t, w, r)

}

func compressDecompress(t *testing.T, w io.WriteCloser, r io.Reader) {

	rnd := rand.New(rand.NewSource(42))

	buf := make([]byte, 2000)
	n, e := rnd.Read(buf)
	if n != len(buf) || e != nil {
		t.Fatal("could not read rnd bytes")
	}

	n, e = w.Write(buf[:1000])
	if n != 1000 || e != nil {
		t.Fatal("could not write all the bytes(1)")
	}
	n, e = w.Write(buf[1000:])
	if n != 1000 || e != nil {
		t.Fatal("could not write all the bytes(2)")
	}
	e = w.Close()
	if e != nil {
		t.Fatal("could not close writer")
	}

	buf2 := make([]byte, 2000)
	n, e = r.Read(buf2)
	if n != len(buf2) || e != nil {
		t.Fatal("could not read all bytes")
	}

	if bytes.Compare(buf, buf2) != 0 {
		t.Fatal("bytes do not match")
	}

}
