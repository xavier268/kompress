package e2e

import (
	"bytes"
	"io"
	"math/rand"
	"testing"

	"github.com/xavier268/kompress"
)

func TestCompressDecompressKDelta(t *testing.T) {

	res := bytes.NewBuffer(nil)
	w := kompress.NewKdeltaWriter(res, 4)
	r := kompress.NewKdeltaReader(res, 4)

	compressDecompress(t, getSourceRnd(t, 2000), w, r)

	data := getKrlenTestData()
	for _, d := range data {
		res = bytes.NewBuffer(nil)
		w = kompress.NewKdeltaWriter(res, 4)
		r = kompress.NewKdeltaReader(res, 4)
		compressDecompress(t, d, w, r)
	}

}

func TestCompressDecompressKrlen(t *testing.T) {

	res := bytes.NewBuffer(nil)

	w := kompress.NewKrlenWriter(res)
	r := kompress.NewKrlenReader(res)

	compressDecompress(t, getSourceRnd(t, 2000), w, r)

	data := getKrlenTestData()
	for _, d := range data {
		res = bytes.NewBuffer(nil)
		w = kompress.NewKrlenWriter(res)
		r = kompress.NewKrlenReader(res)
		compressDecompress(t, d, w, r)
	}

}

func getSourceRnd(t *testing.T, size int) []byte {
	rnd := rand.New(rand.NewSource(42))
	buf := make([]byte, size)
	n, e := rnd.Read(buf)
	if n != len(buf) || e != nil {
		t.Fatal("could not read rnd bytes")
	}
	return buf
}

func compressDecompress(t *testing.T, source []byte, w io.WriteCloser, r io.Reader) {

	n, e := w.Write(source)
	if n != len(source) || e != nil {
		t.Fatal("could not write all the bytes")
	}
	e = w.Close() // IMPORTANT !
	if e != nil {
		t.Fatal("could not close writer")
	}

	buf2 := make([]byte, len(source))
	n, e = r.Read(buf2)
	if n != len(buf2) || e != nil {
		t.Fatal("could not read all bytes")
	}

	if bytes.Compare(source, buf2) != 0 {
		t.Fatal("bytes do not match")
	}

}
