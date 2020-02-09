package kompress

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestKDeltaInternals(t *testing.T) {

	k := new(kdelta)
	k.reset(200)

	if k.expect(55) != 0 {
		t.Fatal("Invalid initial expect value")
	}

	k.learnAndUpdate(55, 12)
	k.learnAndUpdate(45, 22)
	if k.expect(55) != 12 || k.expect(45) != 22 {
		t.Fatal("values were note retrieved")
	}

	k.learnAndUpdate(55, 13)
	if k.expect(55) != 12 {
		t.Fatal("Unnecessary change")
	}
	k.learnAndUpdate(55, 13)
	if k.expect(55) != 13 {
		t.Fatal("Change did not happen")
	}

}

func TestKDeltaBasicWriter(t *testing.T) {

	source := "abcabcabcdabcdeaaacbaa"
	res := bytes.NewBuffer(nil)

	k := NewKdeltaWriter(res, 2)
	n, err := k.Write([]byte(source))
	k.Close()

	if n != len(source) || err != nil {
		t.Fatal(err)
	}
	if len(res.Bytes()) != len(source) {
		fmt.Println("\nSize went from ", len(source), " to ", len(res.Bytes()))
		t.Fatal("not the same byte length")
	}
}

func TestKDeltaCompressDecompress(t *testing.T) {
	source := "abcabcabcdabcdeaaacbaa"
	source = source + source + "lkjkj" + source
	source = source[3:]

	ls := len(source)

	res := bytes.NewBuffer(nil)

	k := NewKdeltaWriter(res, 4)
	n, err := k.Write([]byte(source))

	if err != nil || ls != n {
		t.Fatal(err, ls, n)
	}
	err = k.Close()
	if err != nil {
		t.Fatal(err)
	}

	kk := NewKdeltaReader(res, 4)

	buf2, err := ioutil.ReadAll(kk)
	if bytes.Compare(buf2, []byte(source)) != 0 {
		fmt.Println([]byte(source))
		fmt.Println(buf2)
		t.Fatal("Compress-decompress is not returning the same result")
	}

}
