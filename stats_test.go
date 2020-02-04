package kompress

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

// Stats provides stats for a byte reader.
func Stats(r io.Reader) (int, float64, float64) {

	rr := bufio.NewReader(r)

	var sum, sum2 float64
	var n int
	var b byte
	var err error

	for b, err = rr.ReadByte(); err != nil; {

		n++
		f := float64(b)
		sum += f
		sum2 += f * f
	}

	mean := sum / float64(n)
	vr := sum2/float64(n) - mean*mean
	sigma := math.Sqrt(vr)

	return n, mean, sigma

}

func TestStats(t *testing.T) {
	in, err := os.Open("LICENSE")
	if err != nil {
		panic(err)
	}
	defer in.Close()
	n, m, s := Stats(in)
	fmt.Println("Stats tests : ", n, m, s)

}
