package errors

// The Error interface implements the builtin error interface, and
// implements an additional method that attaches key value pairs to
// the error.
type Error interface {
	Error() string
	With(keyvals ...interface{}) Error
}

// New returns a new error with a given message.
func New(message string) Error {
	var ctx context
	return ctx.New(message)
}

// Wrap creates an error that wraps an existing error.
// If err is nil, Wrap returns nil.
func Wrap(err error, message ...string) Error {
	var ctx context
	return ctx.Wrap(err, message...)
}

// With creates a context with the key/value pairs.
func With(keyvals ...interface{}) Context {
	var ctx context
	return ctx.With(keyvals...)
}

// A Context contains key/value pairs that will be attached to any
// error created or wrapped from that context.
//
// One useful pattern applies to functions that can return errors
// from many places. Define an `errors` variable early in the function:
//  func doSomethingWith(id string, n int) error {
//      // defines a new context with common key/value pairs
//      errors := errors.With("id", id, "n", n)
//
//      // ... later on ...
//
//      if err := doSomething(); err != nil {
//          return errors.Wrap(err, "cannot do something")
//      }
//
//      // ... and later still ...
//
//      if somethingBadHasHappened() {
//          return errors.New("something bad has happened")
//      }
//
//      // ... and so on ...
//
// This pattern ensures that all errors created or wrapped in a function
// have the same key/value pairs attached.
type Context interface {
	With(keyvals ...interface{}) Context
	New(message string) Error
	Wrap(err error, message ...string) Error
}
