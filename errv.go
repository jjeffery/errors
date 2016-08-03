// Package errv provides a simple interface for
// error handling that works well with structured logging.
package errv

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

type _error struct {
	msg    string
	pairs  []keyValueT
	cause  error
	caller string
}

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

// New creates a new error.
func New(msg string, opts ...Option) error {
	return newError(msg, opts...)
}

// Wrap creates an error that wraps an existing error, providing additional fields.
func Wrap(err error, msg string, opts ...Option) error {
	e := newError(msg, opts...)
	e.cause = err
	return e
}

func newError(msg string, opts ...Option) *_error {
	e := &_error{
		msg: msg,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Option represents additional information that can be associated
// with an error.
type Option func(*_error)

// KV associates a key-value pair with an error.
func KV(key string, value interface{}) Option {
	return func(e *_error) {
		kv := keyValueT{
			key:   key,
			value: value,
		}
		e.pairs = append(e.pairs, kv)
	}
}

// KeyValue is a key-value pair containing additional information
// relevant to an error.
type keyValueT struct {
	key   string
	value interface{}
}

// Caller is used to add a key-value pair to the error indicating
// the file and line number for the caller. The argument skip is
// the number of stack frames to ascend, with 0 identifying the
// caller of Caller.
func Caller(skip int) Option {
	return func(e *_error) {
		if pc, file, line, ok := runtime.Caller(skip + 3); ok {
			fn := runtime.FuncForPC(pc)
			file = trimGOPATH(fn.Name(), file)
			e.caller = fmt.Sprintf("%s:%d", file, line)
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
