# errors [![GoDoc](https://godoc.org/github.com/jjeffery/errors?status.svg)](https://godoc.org/github.com/jjeffery/errors) [![License](http://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://raw.githubusercontent.com/jjeffery/errors/master/LICENSE.md) [![Build Status](https://travis-ci.org/jjeffery/errors.svg?branch=master)](https://travis-ci.org/jjeffery/errors) [![Coverage Status](https://coveralls.io/repos/github/jjeffery/errors/badge.svg?branch=master)](https://coveralls.io/github/jjeffery/errors?branch=master) [![GoReportCard](https://goreportcard.com/badge/github.com/jjeffery/errors)](https://goreportcard.com/report/github.com/jjeffery/errors)

Package `errors` provides simple error handling primitives that work well with structured logging.

- [Acknowledgement](#acknowledgement)
- [Background](#background)
- [Creating errors](#creating-errors)
- [Retrieving the cause of an error](#retrieving-the-cause-of-an-error)
- [Retrieving key value pairs for structured logging](#retrieving-key-value-pairs-for-structured-logging)

## Acknowledgement

This package is inspired by the excellent
[github.com/pkg/errors](https://github.com/pkg/errors) package. A significant
amount of code and documentation in this package has been adapted from that
source.

A key difference between this package and github.com/pkg/errors is that
this package has been designed to suit programs that make use of 
[structured logging](https://www.thoughtworks.com/radar/techniques/structured-logging). 
Some of the ideas in this package [were proposed](https://github.com/pkg/errors/issues/34) 
for package github.com/pkg/errors, but after a reasonable amount of consideration, were 
ultimately not included in that package.

> If you are not using structured logging in your application and have no intention
of doing so, you will probably be better off using the 
[github.com/pkg/errors](https://github.com/pkg/errors) package in preference to this one.

## Background

The traditional error handling idiom in Go is roughly akin to
```go
if err != nil {
        return err
}
```
which applied recursively up the call stack results in error reports without context or debugging information. The `errors` package allows programmers to add context to the failure path in their code in a way that does not destroy the original value of the error.

## Creating errors

The `errors` package provides three operations which combine to form a simple yet powerful system for enhancing the value of 
returned errors:

| Operation | Description                                      |
|-----------|--------------------------------------------------|
| New       | create a new error                               | 
| Wrap      | wrap an existing error with an optional message  |
| With      | attach key/value pairs to an error               |

### New &mdash; create a new error

The [`New`](https://godoc.org/github.com/jjeffery/errors#New) function
is used to create an error. This function is compatible with the Go standard 
library `errors` package:

```go
err := errors.New("emit macho dwarf: elf header corrupted")
```

### Wrap &mdash; add a message to an error

The [`Wrap`](https://godoc.org/github.com/jjeffery/errors#Wrap) function 
returns an error that adds a message to the original error. This additional
message can be useful for putting the original error in context. For example:
```go
err := errors.New("permission denied")
fmt.Println(err)

err = errors.Wrap(err, "cannot list directory contents")
fmt.Println(err)

// Output:
// permission denied
// cannot list directory contents: permission denied
```

### With &mdash; add key/value pairs to an error

The [`With`](https://godoc.org/github.com/jjeffery/errors#With) function 
accepts a variadic list of alternating key/value pairs, and returns an error
context that can be used to create a new error or wrap an existing error.

```go
// create new error
err = errors.With("file", "testrun", "line", 101).New("file locked")
fmt.Println(err)

// wrap existing error
err = errors.With("attempt", 3).Wrap(err, "retry failed")
fmt.Println(err)

// Output:
// file locked file=testrun line=101
// retry failed attempt=3: file locked file=testrun line=101
```
One useful pattern is to create an error context that is used for an entire
function scope:

```go
func doSomethingWith(file string, line int) error {
	// set error context
	errors := errors.With("file", file, "line", line)
	
	if number <= 0 {
		// file and line will be attached to the error
		return errors.New("invalid number")
	}
	
	// ... later ...
	
	if err := doOneThing(); err != nil {
		// file and line will be attached to the error
		return errors.Wrap(err, "cannot do one thing")
	}
	
	// ... and so on until ...
	
	return nil
}
```

The errors returned by `New` and `Wrap` provide a `With` method that enables
a fluent-style of error handling:

```go
// create new error
err = errors.New("file locked").With(
    "file", "testrun", 
	"line", 101,
)
fmt.Println(err)

// wrap existing error
err = errors.Wrap(err, "retry failed").With("attempt", 3)
fmt.Println(err)

// Output:
// file locked file=testrun line=101
// retry failed attempt=3: file locked file=testrun line=101
```

(Dave Cheney has written up some 
[good reasons](https://github.com/pkg/errors/issues/15#issuecomment-221194128) 
to avoid a fluent API. Experience will show if this presents a problem, but to 
date it has felt like it leads to simpler, more readable code).

## Retrieving the cause of an error

Using `errors.Wrap` constructs a stack of errors, adding context to the preceding error. Depending on the nature of the error it may be necessary to reverse the operation of `errors.Wrap` to retrieve the original error for inspection. Any error value which implements this interface can be inspected by [`errors.Cause`](https://godoc.org/github.com/jjeffery/errors#Cause).
```go
type causer interface {
        Cause() error
}
```
`errors.Cause` will recursively retrieve the topmost error which does not implement `causer`, which is assumed to be the original cause. For example:
```go
switch err := errors.Cause(err).(type) {
case *MyError:
    // handle specifically
default:
    // unknown error
}
```

## Retrieving key value pairs for structured logging

Errors created by `errors.Wrap` and `errors.New` implement the following interface.
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

[Read the package documentation for more information](https://godoc.org/github.com/jjeffery/errors).

## Licence

MIT

