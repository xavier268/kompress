package kompress

import (
	"bufio"
	"fmt"
	"io"
)

// Kdyndelta compress and decompress suing a "dynamic delta" approach.
// For each byte, the next expected byte is computed, based on past occurences,
// and the actula value is replaced with the delta between actual and expected.
// The total size doex not change, but the frequency and variance of the bytes
// are expected to be easier to compress, escpecially because "long distance"
// correlation can be captured.
type Kdyndelta struct {
	// The buffer of past byte
	buf []byte
	// max capacity of buffer
	capa int
	// map from the computed hash to the expected byte value
	exp map[int]byte
}

// NewKdyndelta constructs the Compresser with the specified capacity.
func NewKdyndelta(capacity int) *Kdyndelta {
	k := new(Kdyndelta)
	k.capa = capacity
	if capacity <= 1 {
		k.capa = 1
		fmt.Println("Warning : Kdyndelta capacity set to 1")
	}
	k.Reset()
	return k
}

// Compress does the forward delta transformation.
func (k *Kdyndelta) Compress(in io.Reader, out io.Writer) error {
	bin := bufio.NewReader(in)
	bout := bufio.NewWriter(out)
	defer bout.Flush()

	for {
		b, err := bin.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		k.buf = append(k.buf, b)
		l := len(k.buf)
		if l > k.capa {
			k.buf = k.buf[l-k.capa:]
		}

		h := k.hash()
		bout.WriteByte(b ^ k.exp[h])
		k.exp[h] = b
	}
	return nil
}

// Decompress does the reverse transformation, from the delta encoding.
// Here, decompress is compress, the operation is its own inverse.
func (k *Kdyndelta) Decompress(in io.Reader, out io.Writer) error {
	return k.Compress(in, out)
}

// Reset returns to the initial state.
func (k *Kdyndelta) Reset() {
	k.buf = make([]byte, 0, k.capa)
	k.exp = make(map[int]byte, 0)
}

// hash will hash the buffer into an int to access the byte map.
func (k *Kdyndelta) hash() int {
	var r int
	for _, c := range k.buf {
		// create hash
		r = r*11 + int(c)*7 + 5
	}
	// limit potential return values to maintain memory size
	return r % (1024 * k.capa)
}
