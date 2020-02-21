package kompress

import (
	"io"
	"testing"
)

func TestRepeatReadWrite(t *testing.T) {

	esc, data := repeatTestData1()
	testDataW(t, esc, data)
	testDataR(t, esc, data)

	esc, data = repeatTestData300()
	testDataW(t, esc, data)
	testDataR(t, esc, data)
}

// Test reading data
func testDataR(t *testing.T, esc Symbol, data [][]Symbol) {

	for i := 0; i < len(data); i += 2 {

		sb := NewSymbolBuffer()
		rr := NewRepeatReader(sb, int(esc))

		// Load buffer
		for _, s := range data[i+1] {
			sb.WriteSymbol(s)
		}

		// Read symbols
		for _, s := range data[i] {
			ss, err := rr.ReadSymbol()
			if err != nil {
				t.Log(err)
				panic("unexpected read error")
			}
			if ss != s {
				t.Log("Expected :", s)
				t.Log("Got      :", ss)
				t.Fatal("Wrong symbol returned")
			}
		}

		// Check no more symbol
		_, err := rr.ReadSymbol()
		if err != io.EOF {
			t.Fatal("Was expecting EOF, but got ", err)
		}
	}
}

// Test writing data
func testDataW(t *testing.T, esc Symbol, data [][]Symbol) {

	for i := 0; i < len(data); i += 2 {

		sb := NewSymbolBuffer()
		rw := NewRepeatWriter(sb, int(esc))

		// Write all symbols
		for _, s := range data[i] {
			rw.WriteSymbol(s)
		}

		// close to flush !
		rw.Close()

		// Check content written
		if len(sb.buf) != len(data[i+1]) {
			t.Log("Expected :", data[i+1])
			t.Log("Got      :", sb.buf)
			t.Fatal("Unexpected result length")
		}

		for j, s := range data[i+1] {
			if s != sb.buf[j] {
				t.Log("Expected :", data[i+1])
				t.Log("Got      :", sb.buf)
				t.Fatal("Unexpected result")
			}
		}

	}

}

func repeatTestData1() (Symbol, [][]Symbol) {
	esc := Symbol(10)
	return esc,
		[][]Symbol{
			[]Symbol{0}, []Symbol{0},
			[]Symbol{0, 0}, []Symbol{0, 0},
			[]Symbol{0, 0, 0}, []Symbol{esc, 0, 0},
			[]Symbol{0, 0, 0, 0}, []Symbol{esc, 1, 0},

			[]Symbol{1}, []Symbol{1},
			[]Symbol{1, 1}, []Symbol{1, 1},
			[]Symbol{1, 1, 1}, []Symbol{esc, 0, 1},
			[]Symbol{1, 1, 1, 1}, []Symbol{esc, 1, 1},

			[]Symbol{0, 1}, []Symbol{0, 1},
			[]Symbol{0, 1, 1}, []Symbol{0, 1, 1},
			[]Symbol{0, 1, 1, 1}, []Symbol{0, esc, 0, 1},
			[]Symbol{0, 1, 1, 1, 1}, []Symbol{0, esc, 1, 1},

			[]Symbol{0, 1, 2}, []Symbol{0, 1, 2},
			[]Symbol{0, 1, 1, 2}, []Symbol{0, 1, 1, 2},
			[]Symbol{0, 1, 1, 1, 2}, []Symbol{0, esc, 0, 1, 2},
			[]Symbol{0, 1, 1, 1, 1, 2}, []Symbol{0, esc, 1, 1, 2},

			[]Symbol{1, 1, 1, 1, 1, 1, 1, 1, 1}, []Symbol{esc, 6, 1},
			[]Symbol{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, []Symbol{esc, 7, 1},
			[]Symbol{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, []Symbol{esc, 8, 1},
			[]Symbol{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, []Symbol{esc, 9, 1},
			[]Symbol{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, []Symbol{esc, 10, 1},

			[]Symbol{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			[]Symbol{esc, 10, 1, 1},

			[]Symbol{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			[]Symbol{esc, 10, 1, 1, 1},

			[]Symbol{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			[]Symbol{esc, 10, 1, esc, 0, 1},

			[]Symbol{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			[]Symbol{esc, 10, 1, esc, 1, 1},
		}
}

func repeatTestData300() (Symbol, [][]Symbol) {
	esc := Symbol(300)

	var src []Symbol

	for i := 0; i < 350; i++ {
		src = append(src, Symbol(1))
	}
	res := []Symbol{esc, 255, 1, esc, (350 - 255 - 3 - 3), 1}
	return esc, [][]Symbol{
		src, res,
		append(src, 3), append(res, 3),
		append(src, 3, 3), append(res, 3, 3),
		append([]Symbol{3}, src...), append([]Symbol{3}, res...),
	}
}
