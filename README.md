# errors [![GoDoc](https://godoc.org/github.com/jjeffery/errors?status.svg)](https://godoc.org/github.com/jjeffery/errors) [![License](http://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://raw.githubusercontent.com/jjeffery/errors/master/LICENSE.md) [![Build Status](https://travis-ci.org/jjeffery/errors.svg?branch=master)](https://travis-ci.org/jjeffery/errors) [![Coverage Status](https://coveralls.io/repos/github/jjeffery/errors/badge.svg?branch=master)](https://coveralls.io/github/jjeffery/errors?branch=master) [![GoReportCard](https://goreportcard.com/badge/github.com/jjeffery/errors)](https://goreportcard.com/report/github.com/jjeffery/errors)

Package `errors` provides a simple error API that works well with structured logging.

- [Acknowledgement](#acknowledgement)
- [Background](#background)
- [Creating Errors](#creating_errors)
- [Retrieving the cause of an error](#retrieving-the-cause-of-an-error)
- [Retrieving key value pairs for structured logging](#retrieving-key-value-pairs-for-structured-logging)

## Acknowledgement

Many of the ideas, some of the code, and some of the text for the documentation
of this package is based on the excellent
[github.com/pkg/errors](https://github.com/pkg/errors) package. 

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

The `errors` package provides three operations which combine to for a simple yet powerful system for enhancing the value of 
returned errors:

| Operation | Description                                      |
|-----------|--------------------------------------------------|
| New       | create a new error                               | 
| Wrap      | wrap an existing error with an optional message  |
| With      | attach key/value pairs to an error               |

### New: create a new error

Create a new error with the The [`errors.New`](https://godoc.org/github.com/jjeffery/errors#New) 
function, which is source-compatible with the Go standard library `errors` package:

```go
err := errors.New("emit macho dwarf: elf header corrupted")
```

### Wrap: add a message to an error

The [`errors.Wrap`](https://godoc.org/github.com/jjeffery/errors#Wrap) function 
returns a new error that adds a message to the original error. For example:
```go
err := errors.New("original cause")
fmt.Println(err)

err = errors.Wrap(err, "cannot do something")
fmt.Println(err)

// Output:
// original cause
// cannot do something: original cause
```

### With: add key/value pairs to an error

The [`errors.With`](https://godoc.org/github.com/jjeffery/errors#With) function 
can be used to create or wrap an error that has key/value pairs attached:

```go
// create new file
err = errors.With("file", "testrun", "line", 101).New("syntax error")

// wrap existing file
err = errors.With("attempt", 2).Wrap(err, "cannot continue")

// Output:
// elf header corrupted file=testrun line=101
// emit macho dwarf attempt=2: elf header corrupted file=testrun line=101
```

The `With` function accepts a variadic list of alternating key/value pairs, and
returns a context that can be used to create a new error or wrap an existing error.
One useful pattern is to create an errors context that is used for an entire
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

This interface works well with the 
[github.com/jjeffery/kv](https://github.com/jjeffery/kv) package, which
provides improved type safety and clarity when working with key value pairs.

> **GOOD ADVICE:** Do not use the `Keyvals` method on an error to retrieve the
individual key/value pairs associated with an error for processing by the
calling program. 

[Read the package documentation for more information](https://godoc.org/github.com/jjeffery/errors).

## Licence

MIT

