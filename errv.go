// Package errv provides a simple interface for
// error handling that works well with structured logging.
package errv

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

// keyValue is a key-value pair containing additional information
// relevant to an error.
type keyValue struct {
	key   string
	value interface{}
}

type _error struct {
	context
	msg   string
	cause error
}

// Error implements the error interface.
func (e *_error) Error() string {
	var buf bytes.Buffer
	buf.WriteString(e.msg)
	for _, kv := range e.pairs {
		if buf.Len() > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteString(kv.key)
		buf.WriteRune('=')
		buf.WriteString(fmt.Sprintf("%v", kv.value))
	}
	if e.caller != "" {
		if buf.Len() > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteString(e.caller)
	}
	if e.cause != nil {
		buf.WriteRune(':')
		buf.WriteRune(' ')
		buf.WriteString(e.cause.Error())
	}
	return buf.String()
}

// Cause implements the causer interface, for compatiblity with
// the github.com/pkg/errors package.
func (e *_error) Cause() error {
	return e.cause
}

// New creates a new error.
func New(msg string, opts ...Option) error {
	var ctx context
	return ctx.newError(nil, msg, opts)
}

// Wrap creates an error that wraps an existing error, optionally providing additional information.
func Wrap(err error, msg string, opts ...Option) error {
	var ctx context
	return ctx.newError(err, msg, opts)
}

// Option represents additional information that can be associated
// with an error.
type Option func(*context)

// KV associates a single key-value pair with an error.
func KV(key string, value interface{}) Option {
	return func(ctx *context) {
		kv := keyValue{
			key:   key,
			value: value,
		}
		ctx.pairs = append(ctx.pairs, kv)
	}
}

// Keyvals provides a way to specifiy multiple key-value pairs.
// The keyvals parameter is a variadic sequence of alternating keys and values.
// The keys must be of type string, otherwise they are ignored.
//
// Function KV provides a more typesafe alternative to Keyvals, although
// it is a little more verbose.
func Keyvals(keyvals ...interface{}) Option {
	return func(ctx *context) {
		for i := 0; i < len(keyvals); i += 2 {
			if k, ok := keyvals[i].(string); ok {
				kv := keyValue{
					key:   k,
					value: keyvals[i+1],
				}
				ctx.pairs = append(ctx.pairs, kv)
			}
		}
	}
}

// Caller is used to add a key-value pair to the error indicating
// the file and line number for the caller. The argument skip is
// the number of stack frames to ascend, with 0 identifying the
// caller of Caller.
func Caller(skip int) Option {
	// additionalSkip is the number of stack frames used by
	// this package in a call to the function returned by
	// this function. It needs to be added to the number
	// if skip frames requested by the calling program.
	const additionalSkip = 4

	return func(ctx *context) {
		if pc, file, line, ok := runtime.Caller(skip + additionalSkip); ok {
			fn := runtime.FuncForPC(pc)
			file = trimGOPATH(fn.Name(), file)
			ctx.caller = fmt.Sprintf("%s:%d", file, line)
		}
	}
}

// trimGOPATH was copied from https://github.com/pkg/errors (Author: Dave Cheney)
// which in turn was adapted from https://github.com/go-stack/stack (Author: Chris Hines).
func trimGOPATH(name, file string) string {
	// Here we want to get the source file path relative to the compile time
	// GOPATH. As of Go 1.6.x there is no direct way to know the compiled
	// GOPATH at runtime, but we can infer the number of path segments in the
	// GOPATH. We note that fn.Name() returns the function name qualified by
	// the import path, which does not include the GOPATH. Thus we can trim
	// segments from the beginning of the file path until the number of path
	// separators remaining is one more than the number of path separators in
	// the function name. For example, given:
	//
	//    GOPATH     /home/user
	//    file       /home/user/src/pkg/sub/file.go
	//    fn.Name()  pkg/sub.Type.Method
	//
	// We want to produce:
	//
	//    pkg/sub/file.go
	//
	// From this we can easily see that fn.Name() has one less path separator
	// than our desired output. We count separators from the end of the file
	// path until it finds two more than in the function name and then move
	// one character forward to preserve the initial path segment without a
	// leading separator.
	const sep = "/"
	goal := strings.Count(name, sep) + 2
	i := len(file)
	for n := 0; n < goal; n++ {
		i = strings.LastIndex(file[:i], sep)
		if i == -1 {
			// not enough separators found, set i so that the slice expression
			// below leaves file unmodified
			i = -len(sep)
			break
		}
	}
	// get back to 0 or trim the leading separator
	file = file[i+len(sep):]
	return file
}
