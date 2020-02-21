package kompress

import (
	"testing"
)

// Defer a call to this function to recover on panic,
// and panic if no panic while testing ...
// NB : FailNow/Fatal do not call defer, and will not be recovered ...
func expectPanic(t *testing.T) {
	if r := recover(); r != nil {
		t.Log("panic was expected, and DID happen. All OK")
	} else {
		t.Fatal("panic was expected and DID NOT happen, FAIL !")
	}
}

func TestExpectPanic(t *testing.T) {

	defer expectPanic(t)
	panic("do you see me ?")

}
