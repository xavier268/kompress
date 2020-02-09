package e2e

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/xavier268/kompress"
)

func TestBasicKrlenWriter(t *testing.T) {

	source := "abbbbcd"
	res := bytes.NewBuffer(nil)

	k := kompress.NewKlogWriter(kompress.NewKrlenWriter(kompress.NewKlogWriter(res)))
	n, e := k.Write([]byte(source))
	if e != nil || n != len(source) {
		fmt.Println(e, n, "<>", len(source))
		t.Fail()
	}
}

func TestDataKrlenCompress(t *testing.T) {
	data := getTestData()
	for i := 0; i < len(data); i += 2 {
		res := bytes.NewBuffer(nil)
		k := kompress.NewKrlenWriter(res)
		n, e := k.Write(data[i])
		k.Close() // IMPORTANT !!!
		if n != len(data[i]) || e != nil {
			fmt.Println(e, n, "<>", len(data[i]))
			t.Fail()
		}
		if bytes.Compare(res.Bytes(), data[i+1]) != 0 {
			fmt.Printf("%d\tFrom     : %v\n", i, data[i])
			fmt.Printf("%d\tExpected : %v\n", i, data[i+1])
			fmt.Printf("%d\tGot      : %v\n", i, res.Bytes())
			t.Fail()
		}
	}
}
func TestDataKrlenDecompress(t *testing.T) {
	data := getTestData()
	for i := 0; i < len(data); i += 2 {
		k := kompress.NewKrlenReader(bytes.NewReader(data[i+1]))
		res, e := ioutil.ReadAll(k)

		if e != nil {
			t.Fatal(e)
		}
		if bytes.Compare(res, data[i]) != 0 {
			fmt.Printf("%d\tExpected : %v\n", i, data[i])
			fmt.Printf("%d\tGot      : %v\n", i, res)
			t.Fail()
		}
	}
}

func getTestData() [][]byte {
	data := [][]byte{
		{1}, {1},
		{1, 2}, {1, 2},
		{1, 2, 3}, {1, 2, 3},
		{1, 2, 3, 4}, {1, 2, 3, 4},

		{}, {},

		{0}, {0, 0},
		{0, 0}, {0, 1, 0},
		{0, 0, 0}, {0, 2, 0},
		{0, 0, 0, 0}, {0, 3, 0},

		{1, 3, 3, 3}, {1, 0, 2, 3},
		{1, 3, 3, 3, 4}, {1, 0, 2, 3, 4},
		{1, 3, 3, 3, 0}, {1, 0, 2, 3, 0, 0},
		{0, 1, 3, 3, 3, 0}, {0, 0, 1, 0, 2, 2, 3, 0},
		{0, 3, 3, 3, 0}, {0, 0, 1, 2, 3, 0},

		{1, 8, 8, 5}, {1, 8, 8, 5},
		{8, 8, 5}, {8, 8, 5},
		{8, 8, 0}, {8, 8, 0, 0},
		{5, 8, 8}, {5, 8, 8},
		{0, 8, 8}, {0, 0, 8, 8},

		{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}, {1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3},
	}

	// Add some long sequences
	seq := []byte{22, 22}
	for i := 2; i <= 255; i++ {
		seq = append(seq, 22)
		data = append(data, seq)
		data = append(data, []byte{0, byte((i % 256)), 22})
	}
	// 256 and beyond ...
	seq = append(seq, 22)
	data = append(data, seq)
	data = append(data, []byte{0, 255, 22, 22})

	seq = append(seq, 22)
	data = append(data, seq)
	data = append(data, []byte{0, 255, 22, 22, 22})

	seq = append(seq, 22)
	data = append(data, seq)
	data = append(data, []byte{0, 255, 22, 0, 2, 22})

	seq = append(seq, 22)
	data = append(data, seq)
	data = append(data, []byte{0, 255, 22, 0, 3, 22})

	return data
}
