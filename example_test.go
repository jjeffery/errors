package errv_test

import (
	"fmt"

	"github.com/jjeffery/errv"
)

func Example() {
	err := errv.New("first error",
		errv.KV("card", "ace"),
		errv.KV("suite", "spades"))
	fmt.Println(err)

	err = errv.Wrap(err, "second error",
		errv.KV("piece", "rook"),
		errv.KV("color", "black"),
		errv.Caller(0))
	fmt.Println(err)

	// Output:
	// first error card=ace suite=spades
	// second error piece=rook color=black github.com/jjeffery/errv/example_test.go:18: first error card=ace suite=spades
}
