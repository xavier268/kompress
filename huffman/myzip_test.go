package huffman

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"testing"
)

func TestMyZipBasic(t *testing.T) {

	source := []byte("Hellooooooo \x00\xFF world !")
	source = append(source, source...)
	source = append(source, source...)

	testReadWriteMyZip(t, source)

	source = []byte("This is standard english text, \x00 et \xff et français, avec des caractères accentués.")
	testReadWriteMyZip(t, source)

	source = []byte("This is standard english text, \x00 et \xff et français, avec des caractères accentués.")
	source = append(source, source...)
	source = append(source, source...)
	testReadWriteMyZip(t, source)

	source = []byte(text2)

	testReadWriteMyZip(t, source)

	source = append(source, source...)
	testReadWriteMyZip(t, source)

}

func BenchmarkMyZip(b *testing.B) {
	source := []byte("Hello world !")
	for n := 0; n < b.N; n++ {
		testReadWriteMyZip(nil, source)
	}
}

func testReadWriteMyZip(t *testing.T, source []byte) {

	bb := bytes.NewBuffer(nil)
	w := NewMyZipWriter(bb)

	n, err := w.Write(source)
	if n != len(source) {
		t.Fatal("Could not write till the end !")
	}
	if err != nil {
		t.Fatal(err)
	}
	if err = w.Close(); err != nil {
		t.Fatal(err)
	}

	fmt.Print("MyZip : ", len(source), "\t==> ", len(bb.Bytes()))

	// read back

	r := NewMyZipReader(bb)
	res := make([]byte, len(source))
	n, err = r.Read(res)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(source) {
		t.Fatal("Unexpected length, initial : ", len(source), ", got :", len(res))
	}

	if bytes.Compare(res, source) != 0 {
		fmt.Println("Sent : ", source)
		fmt.Println("Got  : ", res)
		t.Fatal("source and res are not the same ")
	}

	fmt.Println("  \t==> ", len(res))

	// Try one more read
	b := []byte{0}
	_, err = r.Read(b)
	if err != io.EOF {
		t.Fatal("Expected EOF, but got ", err)
	}

	// print gzip for refernec
	testReadWriteGZIP(t, source)

}

// compress with gzip for reference
func testReadWriteGZIP(t *testing.T, source []byte) {

	bb := bytes.NewBuffer(nil)
	g := gzip.NewWriter(bb)

	g.Write(source)
	g.Close()

	fmt.Print("GZip  : ", len(source), "\t==> ", len(bb.Bytes()))

	u, _ := gzip.NewReader(bb)
	res := make([]byte, len(source))

	u.Read(res)
	fmt.Println("  \t==> ", len(res), " \n ")

}
