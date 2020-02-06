package klog

import (
	"fmt"
	"io"
	"os"
)

// NewReader constructs a new Klog Reader
// nil reader will read from stdin.
func NewReader(r io.Reader) *Reader {
	rr := new(Reader)
	rr.reader = r
	if r == nil {
		rr.reader = os.Stdin
	}
	return rr
}

// A Reader implements the io.Reader interface.
type Reader struct {
	count  int
	reader io.Reader
	err    error
}

// Read will fetch up to len(b) bytes,
// and logs the contents as it flows.
func (r *Reader) Read(buf []byte) (n int, err error) {

	if r.err != nil {
		return 0, r.err
	}

	n, err = r.reader.Read(buf)
	r.err = err

	for i := 0; i < n; i++ {
		if r.count%16 == 0 {
			fmt.Println()
		}
		fmt.Printf("%2X ", buf[i])
		r.count++
	}

	return n, err
}
