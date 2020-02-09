package kompress

import (
	"fmt"
	"io"
)

// NewKlogReader constructs a new Klog Reader
// nil reader will read from stdin.
func NewKlogReader(r io.Reader) io.Reader {
	rr := new(klogReader)
	rr.kreader.reset(r)
	return rr
}

// NewKlogWriter constructs a new Klog Writer
// nil sends to stdout.
func NewKlogWriter(w io.Writer) io.WriteCloser {
	rr := new(klogWriter)
	rr.kwriter.reset(w)
	return rr
}

// A klogReader implements the io.Reader interface.
// Reading means decompressing.
type klogReader struct {
	kreader
	count int
}

// A klogWriter that implements io.Writer.
// Writing means compressing.
type klogWriter struct {
	kwriter
	count int
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
