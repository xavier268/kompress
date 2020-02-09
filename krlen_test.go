package kompress

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestBasicKrlenCompressDecompress(t *testing.T) {

	source := "bbabb\x00\x00bsbbddvsb\x01bbdbdbbbbbbbbb"
	source = source + source + source
	source += string(bytes.Repeat([]byte{35}, 300))

	ls := len(source)

	res := bytes.NewBuffer(nil)
	k := NewKrlenWriter(res)
	n, err := k.Write([]byte(source))
	if n != ls || err != nil {
		t.Fatal(err, n, "<>", ls)
	}
	err = k.Close()
	if err != nil {
		t.Fatal("closing error")
	}

	kk := NewKrlenReader(res)
	buf2, err := ioutil.ReadAll(kk)

	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(buf2, []byte(source)) != 0 {
		fmt.Println([]byte(source))
		fmt.Println(buf2)
		t.Fatal("Compress  - decompress did not return the same bytes")
	}

}
