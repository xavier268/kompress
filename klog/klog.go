package klog

import (
	"fmt"
	"io"
	"os"
)

// NewKlogReader constructs a new Klog Reader
// nil reader will read from stdin.
func NewKlogReader(r io.Reader) io.Reader {
	rr := new(klogReader)
	rr.reader = r
	if r == nil {
		rr.reader = os.Stdin
	}
	return rr
}

// NewKlogWriter constructs a new Klog Writer
// nil sends to stdout.
func NewKlogWriter(w io.Writer) io.Writer {
	rr := new(klogWriter)
	rr.writer = w
	if w == nil {
		rr.writer = os.Stdout
	}
	return rr
}

// A klogReader implements the io.Reader interface.
// Redaing means decompressing.
type klogReader struct {
	count  int
	reader io.Reader
	err    error
}

// A klogWriter that implements io.Writer.
// Writing means compressing.
type klogWriter struct {
	count  int
	writer io.Writer
	err    error
}

// Read will fetch up to len(b) bytes,
// and logs the contents as it flows.
func (r *klogReader) Read(buf []byte) (n int, err error) {

	if r.err != nil {
		return 0, r.err
	}

	n, err = r.reader.Read(buf)
	r.err = err

	for i := 0; i < n; i++ {
		if r.count%16 == 0 {
			fmt.Println()
			fmt.Printf("%6X:\t", r.count)
		}
		fmt.Printf("%2X ", buf[i])
		r.count++
	}

	if err == io.EOF {
		fmt.Println("<EOF>")
	}

	return n, err
}

// Write will write up to len(b) bytes,
// and logs the contents as it flows.
func (r *klogWriter) Write(buf []byte) (n int, err error) {

	if r.err != nil {
		return 0, r.err
	}

	n, err = r.writer.Write(buf)
	r.err = err

	for i := 0; i < n; i++ {
		if r.count%16 == 0 {
			fmt.Println()
			fmt.Printf("%6X:\t", r.count)
		}
		fmt.Printf("%2X ", buf[i])
		r.count++
	}
	if err == io.EOF {
		fmt.Println("<EOF>")
	}
	return n, err
}
