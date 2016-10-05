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
func Wrap(err error, message string) Error {
	var ctx context
	return ctx.Wrap(err, message)
}

// With creates a context with the key/value pairs.
func With(keyvals ...interface{}) Context {
	var ctx context
	return ctx.With(keyvals...)
}

// A Context contains key/value pairs that will be attached to any
// error message created from that context.
type Context interface {
	With(keyvals ...interface{}) Context
	New(message string) Error
	Wrap(err error, message string) Error
}
