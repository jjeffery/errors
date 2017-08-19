package errors_test

import (
	"fmt"

	"github.com/jjeffery/errors"
)

func Example() {
	err := errors.New("first error").With(
		"card", "ace",
		"suite", "spades",
	)
	fmt.Println(err)

	err = errors.Wrap(err, "second error").With(
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

func ExampleWith() {
	// ... if a function has been called with userID and DocumentID ...
	errors := errors.With("userID", userID, "documentID", documentID)

	n, err := doOneThing()
	if err != nil {
		// will include key value pairs for userID and document ID
		fmt.Println(errors.Wrap(err, "cannot do one thing"))
	}

	if err := doAnotherThing(n); err != nil {
		// will include key value pairs for userID, document ID and n
		fmt.Println(errors.Wrap(err, "cannot do another thing").With("n", n))
	}

	if !isValid(userID) {
		// will include key value pairs for userID and document ID
		fmt.Println(errors.New("invalid user"))
	}

	// Output:
	// cannot do one thing userID=u1 documentID=d1: doOneThing: unable to finish
	// cannot do another thing userID=u1 documentID=d1 n=0: doAnotherThing: not working properly
	// invalid user userID=u1 documentID=d1
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

func isValid(s string) bool {
	return false
}

func ExampleNew() {
	name := getNameOfThing()

	if !isValidName(name) {
		fmt.Println(errors.New("invalid name").With("name", name))
	}
	// Output:
	// invalid name name="!not-valid"
}

func doSomething() error {
	return fmt.Errorf("not implemented")
}

func doSomethingWith(name string) error {
	return fmt.Errorf("permission denied")
}

func ExampleWrap() {
	if err := doSomething(); err != nil {
		fmt.Println(errors.Wrap(err, "cannot do something"))
	}

	name := "otherthings.dat"
	if err := doSomethingWith(name); err != nil {
		fmt.Println(errors.Wrap(err, "cannot do something with").With("name", name))
	}

	// Output:
	// cannot do something: not implemented
	// cannot do something with name="otherthings.dat": permission denied
}

func ExampleCause() {
	// tests if an error is a not found error
	type notFounder interface {
		NotFound() bool
	}

	err := getError()

	if notFound, ok := errors.Cause(err).(notFounder); ok {
		fmt.Printf("Not found: %v", notFound.NotFound())
	}
}

func getError() error {
	return fmt.Errorf("not a not found error")
}
