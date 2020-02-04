package kompress

import (
	"bufio"
	"fmt"
	"io"
)

// Krlen is a compressor for encoding repeated bytes.
// It encodes the following way :
// x or xx  where x != 0, is unchanged
// 0 (single zero) is encoded as <0><0>
// 00 is encoded as <0><1><0>
// xxx... is encoded as <0><l><x>, where 2 <= l <= 255
// 000... is encoded as <0><l><0>, where 2 <= l <= 255
// where l is the length of the sequence minus 1 ( l==2, means a sequence of 3, as in xxx)
// so, the max sequence that can be encoded is 255 + 1 = 256
type Krlen struct {
	//  last byte seen
	last byte
	// length of similar bytes already seen.
	len int
}

// NewKrlen constructs a new Kompressor/Decompressor.
func NewKrlen() *Krlen {
	k := new(Krlen)
	k.Reset()
	return k
}

// Reset compressor/decompressor status.
func (k *Krlen) Reset() {
	k.last = 0
	k.len = 0
}

// Compress from into to.
// The error returned is the io.Reader error, including closing the file.
func (k *Krlen) Compress(from io.Reader, to io.Writer) (err error) {

	// Buffer input/output for performance
	in := bufio.NewReader(from)
	out := bufio.NewWriter(to)
	defer out.Flush()

	for {

		// Read byte after byte
		b, err := in.ReadByte()

		if err != nil {
			// purge last byte(s) read
			switch k.len {
			case 0:
			case 1:
				if k.last != 0 {
					out.WriteByte(k.last)
				} else {
					out.Write([]byte{0, 0})
				}
			case 2:
				if k.last != 0 {
					out.WriteByte(k.last)
					out.WriteByte(k.last)
				} else {
					out.Write([]byte{0, 1, 0})
				}
			default:
				if k.len >= 256 {
					panic("logic error - too long sequence at the end")
				}
				out.Write([]byte{0, byte(k.len - 1), k.last})
			}
			return err
		}

		// Starting reading ...
		if k.len == 0 {
			k.last = b
			k.len = 1
			continue
		}

		// One byte already read
		if k.len == 1 {
			if k.last == 0 && k.last != b {
				// emit 0 0
				_, err = out.Write([]byte{0, 0})
				if err != nil {
					return err
				}
				k.last = b
				k.len = 1
				continue
			}
			if k.last != 0 && k.last != b {
				// emit x
				err = out.WriteByte(k.last)
				if err != nil {
					return err
				}
				k.last = b
				k.len = 1
				continue
			}
			if k.last == b {
				// continue the sequence ...
				k.len = k.len + 1
				continue
			}
		}
		if k.len == 2 && k.last != 0 && k.last != b {
			// special case "xx" should not be encoded
			_, err = out.Write([]byte{k.last, k.last})
			if err != nil {
				return err
			}
			k.len = 1
			k.last = b
			continue
		}
		if k.len >= 2 && k.len <= 254 {
			if k.last == b {
				k.len = k.len + 1
				continue
			} else { // k.last != b
				_, err = out.Write([]byte{0, byte(k.len - 1), k.last})
				if err != nil {
					return err
				}
				k.last = b
				k.len = 1
				continue
			}

		}
		if k.len == 255 {
			// emit in any case
			if k.last != b {
				_, err = out.Write([]byte{0, byte(k.len - 1), k.last})
				if err != nil {
					return err
				}
				k.len = 1
				k.last = b
				continue
			} else { // k.last == b
				_, err = out.Write([]byte{0, byte(k.len), k.last})
				if err != nil {
					return err
				}
				k.len = 0
				continue
			}
		}
		panic("State error")

	}
}

// Decompress from into to.
func (k *Krlen) Decompress(from io.Reader, to io.Writer) (err error) {

	// Buffer input/output for performance
	in := bufio.NewReader(from)
	out := bufio.NewWriter(to)
	defer out.Flush()

	for {

		// Read byte after byte
		b, err := in.ReadByte()

		if err == io.EOF && k.len == 0 {
			return err
		}

		if err != nil {
			fmt.Println(err)
			return ErrorInvalidCompressionFormat
		}

		// emit double zero
		if k.len == 1 && k.last == 0 && b == 0 {
			err = out.WriteByte(0)
			if err != nil {
				return err
			}
			k.len = 0
			continue
		}

		// emit compressed data - first step
		if k.len == 1 && k.last == 0 && b != 0 {
			// read more !
			k.len = 2
			k.last = b
			continue
		}

		// emit compressed data - second step
		if k.len == 2 {
			for i := 0; i <= int(k.last); i++ {
				out.WriteByte(b)
			}
			k.len = 0
			k.last = 0
			continue
		}

		// Handle start 0
		if k.len == 0 && k.last == 0 && b == 0 {
			k.len = 1
			k.last = 0
			continue
		}

		// emit immediately the rest
		out.WriteByte(b)
		k.len = 0
		k.last = 0

	}
}
