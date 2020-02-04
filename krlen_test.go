package kompress

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

// Compiler checks
var _ Compresser = new(Krlen)

func TestKBasic(t *testing.T) {
	k := NewKrlen()
	s := "ababbaaa\x00bbbc"
	in := strings.NewReader(s)
	fmt.Printf("From\t%v \nTo  \t", []byte(s))

	err := k.Compress(in, new(Writer))
	if err != nil && err != io.EOF {
		panic(err)
	}
	fmt.Println()
}

func TestKCompress(t *testing.T) {
	k := NewKrlen()
	data := getKrlenTestData()

	for i := 0; i < len(data); i += 2 {
		k.Reset()
		source := data[i]
		expect := data[i+1]
		in := bytes.NewReader(source)
		out := bytes.NewBuffer(nil)
		err := k.Compress(in, out)
		if err != io.EOF && err != nil {
			panic(err)
		}
		fmt.Println("Test sample : ", i/2)
		if bytes.Compare(expect, out.Bytes()) != 0 {
			fmt.Printf("From    \t%v\n", source)
			fmt.Printf("Got     \t%v\n", out.Bytes())
			fmt.Printf("Expected\t%v\n\n", expect)
			panic("Unexpected test result")
		}

	}

}

func TestKDecompress(t *testing.T) {
	k := NewKrlen()
	data := getKrlenTestData()

	for i := 1; i < len(data); i += 2 {
		k.Reset()
		source := data[i]
		expect := data[i-1]
		in := bytes.NewReader(source)
		out := bytes.NewBuffer(nil)
		err := k.Decompress(in, out)
		if err != io.EOF && err != nil {
			panic(err)
		}
		fmt.Println("Test sample : ", i/2)
		if bytes.Compare(expect, out.Bytes()) != 0 {
			fmt.Printf("From    \t%v\n", source)
			fmt.Printf("Got     \t%v\n", out.Bytes())
			fmt.Printf("Expected\t%v\n\n", expect)
			panic("Unexpected test result")
		}

	}
}

func TestCompressDecompress(t *testing.T) {

	k := NewKrlen()
	kk := NewKrlen()
	data := getKrlenTestData()

	for i := 0; i < len(data); i++ {
		k.Reset()
		kk.Reset()

		source := data[i]
		in1 := bytes.NewReader(source)
		out1 := bytes.NewBuffer(nil)

		err := k.Compress(in1, out1)
		if err != io.EOF && err != nil {
			panic(err)
		}

		in2 := bytes.NewReader(out1.Bytes())
		out2 := bytes.NewBuffer(nil)

		err = kk.Decompress(in2, out2)
		if err != io.EOF && err != nil {
			panic(err)
		}

		fmt.Println("Test sample : ", i)
		if bytes.Compare(source, out2.Bytes()) != 0 {
			fmt.Printf("From    \t%v\n", source)
			fmt.Printf("Got     \t%v\n", out1.Bytes())
			fmt.Printf("Then    \t%v\n\n", out2.Bytes())
			panic("Unexpected non idempotent test")
		}

	}

}

func TestDecompressCompress(t *testing.T) {
	// This test would fail, and this is NORMAL.
	// Some formats are not fully anticipated as invalid compressed format,
	// such as : 0000 uncompr>>> 00 compress>>> 010
	t.Skip()

	k := NewKrlen()
	kk := NewKrlen()
	data := getKrlenTestData()

	for i := 0; i < len(data); i++ {
		k.Reset()
		kk.Reset()

		source := data[i]
		in1 := bytes.NewReader(source)
		out1 := bytes.NewBuffer(nil)

		err := k.Decompress(in1, out1)
		if err != io.EOF && err != nil && err != ErrorInvalidCompressionFormat {
			panic(err)
		}
		if err == ErrorInvalidCompressionFormat {
			fmt.Println(err)
			fmt.Printf("From    \t%v\n", source)
			fmt.Printf("Got     \t%v\n\n", out1.Bytes())
		} else {
			// Only continue testing if uncompression was possible ...
			in2 := bytes.NewReader(out1.Bytes())
			out2 := bytes.NewBuffer(nil)

			err = kk.Compress(in2, out2)
			if err != io.EOF && err != nil {
				panic(err)
			}

			fmt.Println("Test sample : ", i)
			if bytes.Compare(source, out2.Bytes()) != 0 {
				fmt.Printf("From    \t%v\n", source)
				fmt.Printf("Got     \t%v\n", out1.Bytes())
				fmt.Printf("Then    \t%v\n\n", out2.Bytes())
				panic("Unexpected non idempotent test")
			}
		}
	}

}

func getKrlenTestData() [][]byte {
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
		{0, 1, 3, 3, 3, 0}, {0, 0, 1, 0, 2, 3, 0, 0},
		{0, 3, 3, 3, 0}, {0, 0, 0, 2, 3, 0, 0},

		{1, 8, 8, 5}, {1, 8, 8, 5},
		{8, 8, 5}, {8, 8, 5},
		{8, 8, 0}, {8, 8, 0, 0},
		{5, 8, 8}, {5, 8, 8},
		{0, 8, 8}, {0, 0, 8, 8},
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
