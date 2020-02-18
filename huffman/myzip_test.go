package huffman

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestMyZipBasic(t *testing.T) {

	bb := bytes.NewBuffer(nil)
	w := NewMyZipWriter(bb)

	source := []byte("Hellooooooo \x00\xFF world !")

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

	fmt.Println("From : ", len(source), source)
	fmt.Println("To   : ", len(bb.Bytes()), bb.Bytes())

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
		fmt.Println("Seent : ", source)
		fmt.Println("Got   : ", res)
		t.Fatal("source and res are not the same ")
	}

	fmt.Println("Back : ", len(res), res)

	// Try one more read
	b := []byte{0}
	_, err = r.Read(b)
	if err != io.EOF {
		t.Fatal("Expected EOF, but got ", err)
	}

}
