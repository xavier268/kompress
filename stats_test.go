package kompress

import (
	"bufio"
	"io"
	"math"
)

// Stats provides stats for a byte reader.
func Stats(r io.Reader) (int, float64, float64) {

	rr := bufio.NewReader(r)

	var sum, sum2 float64
	var n int
	var err error
	var b byte

	for {

		b, err = rr.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

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
