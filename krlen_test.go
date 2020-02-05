package kompress

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// Compiler checks
var _ Compresser = new(Krlen)

func TestKrlenBasic(t *testing.T) {
	k := NewKrlen()
	s := "ababbaaa\x00bbbc"
	in := strings.NewReader(s)
	fmt.Printf("From\t%v \nTo  \t", []byte(s))

	err := k.Compress(in, new(Dumper))
	if err != nil && err != io.EOF {
		panic(err)
	}
	fmt.Println()
}

func TestKrlenCompress(t *testing.T) {
	k := NewKrlen()
	data := getTestData()

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
		//fmt.Println("Test sample : ", i/2)
		if bytes.Compare(expect, out.Bytes()) != 0 {
			fmt.Printf("From    \t%v\n", source)
			fmt.Printf("Got     \t%v\n", out.Bytes())
			fmt.Printf("Expected\t%v\n\n", expect)
			panic("Unexpected test result")
		}

	}

}

func TestKrlenDecompress(t *testing.T) {
	k := NewKrlen()
	data := getTestData()

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
		//fmt.Println("Test sample : ", i/2)
		if bytes.Compare(expect, out.Bytes()) != 0 {
			fmt.Printf("From    \t%v\n", source)
			fmt.Printf("Got     \t%v\n", out.Bytes())
			fmt.Printf("Expected\t%v\n\n", expect)
			panic("Unexpected test result")
		}

	}
}

func TestKrlenCompressDecompress(t *testing.T) {

	k := NewKrlen()
	kk := NewKrlen()
	data := getTestData()

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

		//fmt.Println("Test sample : ", i)
		if bytes.Compare(source, out2.Bytes()) != 0 {
			fmt.Printf("From    \t%v\n", source)
			fmt.Printf("Got     \t%v\n", out1.Bytes())
			fmt.Printf("Then    \t%v\n\n", out2.Bytes())
			panic("Unexpected non idempotent test")
		}

	}

}

func TestKrlenDecompressCompress(t *testing.T) {
	// This test would fail, and this is NORMAL.
	// Some formats are not fully anticipated as invalid compressed format,
	// such as : 0000 uncompr>>> 00 compress>>> 010
	t.Skip()

	k := NewKrlen()
	kk := NewKrlen()
	data := getTestData()

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

			//fmt.Println("Test sample : ", i)
			if bytes.Compare(source, out2.Bytes()) != 0 {
				fmt.Printf("From    \t%v\n", source)
				fmt.Printf("Got     \t%v\n", out1.Bytes())
				fmt.Printf("Then    \t%v\n\n", out2.Bytes())
				panic("Unexpected non idempotent test")
			}
		}
	}

}

func TestKrlenStats(t *testing.T) {
	in, err := os.Open("LICENSE")
	if err != nil {
		panic(err)
	}

	n, m, s := Stats(in)
	fmt.Println("Stats before compress : ", n, m, s)
	in.Close()

	in, err = os.Open("LICENSE")
	if err != nil {
		panic(err)
	}
	defer in.Close()

	out := bytes.NewBuffer(nil)
	k := NewKrlen()
	k.Compress(in, out)

	n, m, s = Stats(bytes.NewReader(out.Bytes()))
	fmt.Println("Stats after compress : ", n, m, s)
}
