/*
Package errors provides a simple error API that works well with structured logging.

Many of the ideas, much of the code, and even the text for the documentation
of this package is based on the excellent
github.com/pkg/errors package (https://github.com/pkg/errors).

A key difference between this package and github.com/pkg/errors is that
this package has been designed to suit programs that make use of
structured logging.
Some of the ideas in this package were proposed for package
github.com/pkg/errors, but after a reasonable amount of consideration, were
ultimately not included in that package.
(See https://github.com/pkg/errors/issues/34 for details).

If you are not using structured logging in your application and have no intention
of doing so, use the github.com/pkg/errors package in preference to this one.

Background

The traditional error handling idiom in Go is roughly akin to
 if err != nil {
     return err
 }
which applied recursively up the call stack results in error reports without
context or debugging information. The errors package allows programmers to
add context to the failure path in their code in a way that does not destroy
the original value of the error.

Adding context to an error

The errors.Wrap function returns a new error that adds context to the
original error. For example
 name := "some-file"
 number := 53
 err := doSomethingWith(name, number)
 if err != nil {
     return errors.Wrap(err, "cannot do something").With(
         "name", name,
         "number", number,
     )
 }

Retrieving the cause of an error

Using errors.Wrap constructs a stack of errors, adding context to the
preceding error. Depending on the nature of the error it may be necessary
to reverse the operation of errors.Wrap to retrieve the original error for
inspection. Any error value which implements this interface can be inspected
by errors.Cause.

 type causer interface {
     Cause() error
 }
errors.Cause will recursively retrieve the topmost error which does not
implement causer, which is assumed to be the original cause. For example:

 switch err := errors.Cause(err).(type) {
 case *MyError:
     // handle specifically
 default:
     // unknown error
 }

Retrieving key value pairs for structured logging

Errors created by `errors.Wrap` and `errors.New` implement the following
interface:

 type keyvalser interface {
     Keyvals() []interface{}
 }

The Keyvals method returns an array of alternating keys and values. The
first key will always be "msg" and its value will be a string containing
the message associated with the wrapped error.

Example using go-kit logging (https://github.com/go-kit/kit/tree/master/log):

 // logError logs details of an error to a structured error log.
 func logError(logger log.Logger, err error) {
     // start with timestamp and error level
     keyvals := []interface{}{
         "ts",    time.Now().Format(time.RFC3339Nano),
         "level", "error",
     }

     type keyvalser interface {
         Keyvals() []interface{}
     }
     if kv, ok := err.(keyvalser); ok {
         // error contains structured information, first key/value
         // pair will be "msg".
         keyvals = append(keyvals, kv.Keyvals()...)
     } else {
         // error does not contain structured information, use the
         // Error() string as the message.
         keyvals = append(keyvals, "msg", err.Error())
     }
     logger.Log(keyvals...)
 }

This interface works well with the github.com/jjeffery/kv package, which
provides improved type safety and clarity when working with key value pairs.

GOOD ADVICE: Do not use the `Keyvals` method on an error to retrieve the
individual key/value pairs associated with an error for processing by the
calling program.
*/
package errors

// BUG(jpj): This package makes use of a fluent API for attaching key/value
// pairs to an error. Dave Cheney has written up some good reasons to avoid
// this approach: see https://github.com/pkg/errors/issues/15#issuecomment-221194128.
// Experience will show if this presents a problem.

// BUG(jpj): Attaching key/value pairs to an error was considered for package
// github.com/pkg/errors, but in the end it was not implemented because
// of the potential for abusing the information in the error. See Dave Cheney's
// comment at https://github.com/pkg/errors/issues/34#issuecomment-228231192.
// This package has used the `keyvalser` interface as a mechanism for extracting
// key/value pairs from an error. In practice this seems to work quite well, but
// it would be possible to write code that abuses this interface by extractng 
// information from the error for use by the program.
