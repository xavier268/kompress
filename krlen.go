package kompress

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// KrlenWriter is a compressor for encoding repeated bytes.
// It encodes using an escape byte, <esc>
// x or xx  where x != <esc>, is unchanged
// <esc> (single esc) is encoded as <esc><0>
// <esc><esc> is encoded as <esc><1><esc>
// xxx... is encoded as <esc><l><x>, where 2 <= l <= 255
// <esc><esc><esc> ... is encoded as <esc><l><esc>, where 2 <= l <= 255
// where l is the length of the sequence minus 1 ( l==2, means a sequence of 3, as in xxx)
// so, the max sequence that can be encoded is 255 + 1 = 256
type krlenWriter struct {
	kwriter
	kfreq
	last byte // Last byte seen
	len  int  // Sequence length so far ...
}

type krlenReader struct {
	kreader
	kfreq
	// buffer available for the next read
	// since we may have too many bytes to send at once.
	buf []byte
}

// NewKrlenWriter will create a writer that compress run-length sequences.
// Ensure you close it at the end, to flush pending bytes.
func NewKrlenWriter(w io.Writer) io.WriteCloser {
	k := new(krlenWriter)
	if w == nil {
		k.writer = os.Stdout
	} else {
		k.writer = w
	}
	return k
}

// NewKrlenReader will create a reader that decompress run-length sequences.
func NewKrlenReader(r io.Reader) io.Reader {
	k := new(krlenReader)
	if r == nil {
		k.reader = os.Stdin
	} else {
		k.reader = r
	}
	return k
}

// Read (ie decompress) from the provided reader
func (k *krlenReader) Read(res []byte) (int, error) {

	if k.err != nil {
		return 0, k.err
	}
	// This is how many bytes we want ...
	nres := len(res)

	var e error
	var b byte
	// Loop until we have enough in buffer
	for len(k.buf) <= nres {
		b, e = k.readByte()
		if e != nil {
			break
		}
		switch b {
		case k.esc:
			bb, e := k.readByte()
			switch {
			case e != nil:
				k.err = errors.New(e.Error() + " : invalid compression format")
			case bb == 0:
				k.buf = append(k.buf, k.esc) // single escape
				k.update(k.esc)
			default:
				x, e := k.readByte()
				if e != nil {
					k.err = errors.New(e.Error() + " : invalid compression format")
				} else {
					for i := 0; i < int(bb)+1; i++ {
						k.buf = append(k.buf, x)
						k.update(x)
					}

				}
			}

		default:
			k.buf = append(k.buf, b)
			k.update(b)
		}

	}

	if e == nil {
		// we should have enough ...
		if len(k.buf) < nres {
			panic("internal error - we should have a full buffer")
		}
		copy(res, k.buf[0:nres])
		k.buf = k.buf[nres:]
		return nres, nil
	}

	// e != nil, but DO NOT ASSUME buffer is smaller than res !!
	n := len(k.buf)
	if n < nres {
		copy(res, k.buf)
		k.err = e
		return n, e
	}
	// Now, buffer is larger than res - do not send the error (EOF) yet !
	copy(res, k.buf[:nres])
	k.buf = k.buf[:nres]
	// Don't send error yet, until buffer is empty !
	return nres, nil

}

// Write( ie Compress )  into to the provided writer.
func (k *krlenWriter) Write(buf []byte) (int, error) {

	if k.err != nil {
		return 0, k.err
	}

	// read byte after byte ...
	for _, b := range buf {

		switch {
		case k.len == 0:
			// just start a new sequence ...
			k.last = b
			k.len = 1
		case k.len == 1:
			switch {
			case k.last != k.esc && b != k.last: // single non-escape
				k.emit(k.last)
				k.update(k.last)
				k.last = b
				k.len = 1
			case k.last == k.esc && b != k.last: // single esc
				k.emit(k.esc, 0)
				k.update(k.last)
				k.last = b
				k.len = 1
			case b == k.last:
				k.len++
				// dont update freq as we follow the sequence ...
			default:
				panic("internal error")
			}

		case k.len == 2:
			switch {
			case k.last == k.esc && b != k.last: // double esc confirmed
				k.emit(k.esc, 1, k.esc)
				k.update(k.last)
				k.update(k.last)
				k.last = b
				k.len = 1
			case k.last != k.esc && b != k.last: // double non esc
				k.emit(k.last, k.last)
				k.update(k.last)
				k.update(k.last)
				k.last = b
				k.len = 1
			case b == k.last:
				k.len++
			default:
				panic("internal error")
			}
		case k.len > 2 && k.len < 256:
			switch {
			case b != k.last:
				k.emit(k.esc, byte(k.len-1), k.last)
				for i := 0; i < k.len; i++ {
					k.update(k.last)
				}
				k.last = b
				k.len = 1
			case b == k.last:
				k.len++
			}
		case k.len == 256: // emit now !
			k.emit(k.esc, 255, k.last)
			for i := 0; i < 256; i++ {
				k.update(k.last)
			}
			k.last = b
			k.len = 1
		default:
			panic("internal error")
		}

	}
	return len(buf), nil
}

// Close the writer, flushing remaining data if any, then closing.
func (k *krlenWriter) Close() error {

	switch {
	case k.len == 0: // do nothing
	case k.last == k.esc && k.len == 1:
		k.emit(k.esc, 0)
	case k.last == k.esc && k.len == 2:
		k.emit(k.esc, 1, k.esc)
	case k.last != k.esc && k.len == 1:
		k.emit(k.last)
	case k.last != k.esc && k.len == 2:
		k.emit(k.last, k.last)
	case k.len <= 256:
		k.emit(k.esc, byte(k.len-1), k.last)
	default:
		fmt.Printf("Invalid state : len = %d, last = %02X\n", k.len, k.last)
		panic("invalid state")
	}
	return k.kwriter.Close()
}
