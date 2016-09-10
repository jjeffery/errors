/*
Package errorv provides a simple error API that works well with structured logging.

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
context or debugging information. The errorv package allows programmers to
add context to the failure path in their code in a way that does not destroy
the original value of the error.

Adding context to an error

The errorv.Wrap function returns a new error that adds context to the
original error. For example
 name := "some-file"
 number := 53
 err := doSomethingWith(name, number)
 if err != nil {
     return errorv.Wrap(err, "cannot do something",
         "name", name,
         "number", number,
     )
 }

Retrieving the cause of an error

Using errorv.Wrap constructs a stack of errors, adding context to the
preceding error. Depending on the nature of the error it may be necessary
to reverse the operation of errorv.Wrap to retrieve the original error for
inspection. Any error value which implements this interface can be inspected
by errorv.Cause.

 type causer interface {
     Cause() error
 }
errorv.Cause will recursively retrieve the topmost error which does not
implement causer, which is assumed to be the original cause. For example:

 switch err := errorv.Cause(err).(type) {
 case *MyError:
     // handle specifically
 default:
     // unknown error
 }

Retrieving key value pairs for structured logging

Errors created by `errorv.Wrap` and `errorv.New` implement the following
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
package errorv
