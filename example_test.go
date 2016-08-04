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

var userID, documentID string

func ExampleContext() error {
	// ... if a function has been called with userID and DocumentID ...
	errv := errv.NewContext(errv.KV("userID", userID), errv.KV("documentID", documentID))

	n, err := doOneThing()
	if err != nil {
		// will include key value pairs for userID and document ID
		return errv.Wrap(err, "cannot do one thing")
	}

	if err := doAnotherThing(n); err != nil {
		// will include key value pairs for userID, document ID and n
		return errv.Wrap(err, "cannot do another thing", errv.KV("n", n))
	}

	return nil
}

func doOneThing() (int, error) {
	return 0, nil
}

func doAnotherThing(n int) error {
	return nil
}
