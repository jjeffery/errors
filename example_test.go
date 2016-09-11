package errorv_test

import (
	"fmt"

	"github.com/jjeffery/errorv"
)

func Example() {
	err := errorv.New("first error",
		"card", "ace",
		"suite", "spades")
	fmt.Println(err)

	err = errorv.Wrap(err, "second error",
		"piece", "rook",
		"color", "black",
	)
	fmt.Println(err)

	// Output:
	// first error card=ace suite=spades
	// second error piece=rook color=black: first error card=ace suite=spades
}

var userID = "u1"
var documentID = "d1"

func ExampleContext() {
	// ... if a function has been called with userID and DocumentID ...
	errorv := errorv.NewContext("userID", userID, "documentID", documentID)

	n, err := doOneThing()
	if err != nil {
		// will include key value pairs for userID and document ID
		fmt.Println(errorv.Wrap(err, "cannot do one thing"))
	}

	if err := doAnotherThing(n); err != nil {
		// will include key value pairs for userID, document ID and n
		fmt.Println(errorv.Wrap(err, "cannot do another thing", "n", n))
	}

	// Output:
	// cannot do one thing userID=u1 documentID=d1: doOneThing: unable to finish
	// cannot do another thing userID=u1 documentID=d1 n=0: doAnotherThing: not working properly
}

func doOneThing() (int, error) {
	return 0, fmt.Errorf("doOneThing: unable to finish")
}

func doAnotherThing(n int) error {
	return fmt.Errorf("doAnotherThing: not working properly")
}

var NotFound error

func getNameOfThing() string {
	return "!not-valid"
}

func thingExists(name string) bool {
	return false
}

func doSomethingWithThing(name string) error {
	return nil
}

func isValidName(name string) bool {
	return false
}

func ExampleNew() {
	name := getNameOfThing()

	if !isValidName(name) {
		fmt.Println(errorv.New("invalid name", "name", name))
	}
	// Output:
	// invalid name name=!not-valid
}

func ExampleCause() {
	// tests if an error is a not found error
	type notFounder interface {
		NotFound() bool
	}

	err := getError()

	if notFound, ok := errorv.Cause(err).(notFounder); ok {
		fmt.Printf("Not found: %v", notFound.NotFound())
	}
}

func getError() error {
	return fmt.Errorf("not a not found error")
}
