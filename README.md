# errorv [![GoDoc](https://godoc.org/github.com/jjeffery/errorv?status.svg)](https://godoc.org/github.com/jjeffery/errorv)

Package `errorv` provides a simple error API that works well with structured logging.

Many of the ideas, much of the code, and even the text for the documentation
of this package is based on the excellent
[github.com/pkg/errors](https://github.com/pkg/errors) package. 
Please check out this package as it is more mature and has undergone more 
testing and scrutiny. It may suit your purpose better.

A key difference between this package and `github.com/pkg/errors` is that
this package has been designed to suit programs that make use of 
[structured logging](http://dev.splunk.com/view/logging-best-practices/SP-CAAADP6). 
Some of the ideas in this package [were proposed](https://github.com/pkg/errors/issues/34) 
for package `github.com/pkg/errors`, but after a reasonable amount of consideration, were 
ultimately not included in that package.

> If you are not using structured logging in your application and have no intention
of doing so, use the [github.com/pkg/errors](https://github.com/pkg/errors) package
in preference to this one.

## Background

The traditional error handling idiom in Go is roughly akin to
```go
if err != nil {
        return err
}
```
which applied recursively up the call stack results in error reports without context or debugging information. The `errorv` package allows programmers to add context to the failure path in their code in a way that does not destroy the original value of the error.

## Adding context to an error

The [`errorv.Wrap`](https://godoc.org/github.com/jjeffery/errorv#Wrap) function 
returns a new error that adds context to the original error. For example
```go
name := "some-file"
err := doSomethingWith(name)
if err != nil {
        return errorv.Wrap(err, "cannot do something",
		        errorv.KV("name", name))
}
```

Adding context to an error with the `errorv.Wrap` function involves attaching
a message and any number of key value pairs. There are alternatives for attaching
key value pairs to a message, depending on the convenience and preferences of
the caller:

```go
return errorv.Wrap(err, "message for wrapped error",
        errorv.KV("key1", value1),
		errorv.KV("key2", value2),
		errorv.KV("key3", value3))
```
can also be expressed as:
```go
return errorv.Wrap(err, "message for wrapped error", errorv.Keyvals{
	    "key1", value1,
		"key2", value2,
		"key3", value3,
})
```
The first way is a bit more verbose, but has stronger type checks.

## Retrieving the cause of an error

Using `errorv.Wrap` constructs a stack of errors, adding context to the preceding error. Depending on the nature of the error it may be necessary to reverse the operation of `errorv.Wrap` to retrieve the original error for inspection. Any error value which implements this interface can be inspected by [`errorv.Cause`](https://godoc.org/github.com/jjeffery/errorv#Cause).
```go
type causer interface {
        Cause() error
}
```
`errorv.Cause` will recursively retrieve the topmost error which does not implement `causer`, which is assumed to be the original cause. For example:
```go
switch err := errorv.Cause(err).(type) {
case *MyError:
        // handle specifically
default:
        // unknown error
}
```

## Retrieving the error for structured logging

Errors created by `errorv.Wrap` implement the following interface.
```go
type keyvalser interface {
	Keyvals() []interface{}
}
```
The `Keyvals` method returns an array of alternating keys and values. The
first key will always be "msg" and its value will be a string containing
the message associated with the wrapped error.

Example using [go-kit logging](https://github.com/go-kit/kit/tree/master/log):

```go
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
```

> **GOOD ADVICE:** Do not use the `Keyvals` method on an error to retrieve the
individual key/value pairs associated with an error for processing by the
calling program.

[Read the package documentation for more information](https://godoc.org/github.com/jjeffery/errorv).

## Licence

MIT

