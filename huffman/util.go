package huffman

import (
	"fmt"
	"io"
)

// LogWriter is a Writer that can write bits or Symbols
// to stdout for debugging.
type LogWriter struct{}

// compiler checks
var _ BitWriteCloser = new(LogWriter)
var _ SymbolWriteCloser = new(LogWriter)

// Close file
func (l *LogWriter) Close() error {
	fmt.Println("\n--- LogWriter was closed ---")
	return nil
}

// WriteBit to stdout
func (l *LogWriter) WriteBit(b Bit) error {
	if b == 0 {
		fmt.Print("0")
	} else {
		fmt.Print("1")
	}
	return nil
}

// WriteBits to stdout
func (l *LogWriter) WriteBits(bb ...Bit) error {
	for _, b := range bb {
		l.WriteBit(b)
	}
	return nil
}

// WriteSymbol to stdout
func (l *LogWriter) WriteSymbol(s Symbol) error {
	fmt.Printf("symbol(%d)", s)
	return nil
}

// ========================  BitBuffer ============================

// A BitBuffer can be read or written bits, or even bytes.
// It size adapts dynamically.
type BitBuffer struct {
	bits []Bit
}

// NewBitBuffer with optionnal preset bits.
func NewBitBuffer(bits ...Bit) *BitBuffer {
	bb := new(BitBuffer)
	bb.bits = make([]Bit, len(bits))
	for i, b := range bits {
		bb.bits[i] = b
	}
	return bb
}

// Compiler checks
var _ BitReader = new(BitBuffer)
var _ BitWriteCloser = new(BitBuffer)

// ReadBit from buffer
func (bb *BitBuffer) ReadBit() (Bit, error) {
	if len(bb.bits) == 0 {
		return 0, io.EOF
	}
	b := bb.bits[0]
	bb.bits = bb.bits[1:]
	if b == 0 {
		return 0, nil
	}
	return 1, nil
}

// Close does nothing
func (bb *BitBuffer) Close() error {
	return nil
}

// WriteBit to buffer
func (bb *BitBuffer) WriteBit(b Bit) error {
	bb.bits = append(bb.bits, b)
	return nil
}

// WriteBits to buffer
func (bb *BitBuffer) WriteBits(b ...Bit) error {
	bb.bits = append(bb.bits, b...)
	return nil
}

// Size of the buffer.
func (bb *BitBuffer) Size() int {

	return len(bb.bits)
}

// ReadByte reads a byte,
// byte is 0 padded on the right in needed
func (bb *BitBuffer) ReadByte() (byte, error) {

	var by byte
	var b Bit
	var err error

	for i := 0; i < 8; i++ {
		by = by << 1
		b, err = bb.ReadBit()
		if err == nil {
			by = by | byte(b)
		}
	}
	return by, err
}

// WriteByte writes a byte, bit by bit.
func (bb *BitBuffer) WriteByte(b byte) error {

	var err error
	var mask byte = byte(128)

	for i := 0; i < 8; i++ {
		err = bb.WriteBit(Bit(mask & b))
		mask = mask / 2
	}
	return err
}

// =============== BitFromByteReader ======================

// A BitFromByteReader object encapsulates an io.Reader and provides
// bit read . Internal buffer is kept to a minimum
// (less than 8 bits).
type BitFromByteReader struct {
	reader io.Reader
	// mask pointing to next available bit
	mask byte
	// buffer
	buf [1]byte
}

// compiler cheks
var _ BitReader = new(BitFromByteReader)

// NewBitFromByteReader constrstructed from io.Reader.
// Performance buffering is up to the io.Reader.
func NewBitFromByteReader(ior io.Reader) *BitFromByteReader {
	r := new(BitFromByteReader)
	r.reader = ior
	return r
}

// ReadBit from the underlying io.Reader.
func (r *BitFromByteReader) ReadBit() (Bit, error) {

	if r.mask <= 0 {
		r.mask = 0b_1000_0000
		_, err := r.reader.Read(r.buf[:])
		if err != nil {
			return 0, err
		}
	}

	bit := r.buf[0] & r.mask
	r.mask = r.mask >> 1
	if bit == 0 {
		return 0, nil
	}
	return 1, nil

}

// =============== BitToByteWriter ======================

// A BitToByteWriter object encapsulates an io.Writer  and provides
// bit write capabilities. Internal buffer is kept to a minimum
// (less than 8 bits). Closing writes is essential to trigger padding  and flush.
type BitToByteWriter struct {
	writer io.Writer
	// mask pointing to next available bit
	mask byte
	// buffer
	buf [1]byte
}

// compiler cheks
var _ BitWriteCloser = new(BitToByteWriter)

// NewBitToByteWriter from an io.Writer.
func NewBitToByteWriter(ior io.Writer) *BitToByteWriter {
	w := new(BitToByteWriter)
	w.mask = 0b_1000_0000
	w.writer = ior
	return w
}

// WriteBit a single bit, and flush if full byte is available.
func (w *BitToByteWriter) WriteBit(b Bit) (err error) {

	if b != 0 {
		w.buf[0] |= w.mask
	}

	w.mask = w.mask >> 1
	if w.mask == 0 {
		w.mask = 128
		_, err = w.writer.Write(w.buf[:])
		w.buf[0] = 0
	}

	return err

}

// WriteBits zero or more bits,  flushing every full byte available.
func (w *BitToByteWriter) WriteBits(b ...Bit) (err error) {
	for i := range b {
		err = w.WriteBit(b[i])
	}
	return err
}

// Close and flush, padding with 0.
// It tries to close the underlying io.Writer if it is also a Closer.
func (w *BitToByteWriter) Close() (err error) {

	// Write 0 until pending byte is flushed
	for w.mask != 128 {
		err = w.WriteBit(0)
	}

	if c, ok := w.writer.(io.Closer); ok {
		return c.Close()
	}
	return err

}
