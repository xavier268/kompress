package kompress

import "fmt"

type Writer struct{}

func (w *Writer) Write(b []byte) (n int, err error) {

	fmt.Printf("%v\n", b)
	return len(b), nil

}
