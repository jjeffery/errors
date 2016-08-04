// Package errv provides a simple interface for
// error handling that works well with structured logging.
package errv

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

// KeyValue is a key-value pair containing additional information
// relevant to an error.
type keyValueT struct {
	key   string
	value interface{}
}

type context struct {
	pairs  []keyValueT
	caller string
}

// clone creates a deep copy of the context.
func (ctx context) clone() context {
	ctx2 := context{
		caller: ctx.caller,
	}
	if len(ctx.pairs) > 0 {
		ctx2.pairs = make([]keyValueT, len(ctx.pairs))
		copy(ctx2.pairs, ctx.pairs)
	}
	return ctx2
}

func (ctx *context) applyOptions(opts []Option) {
	for _, opt := range opts {
		opt(ctx)
	}
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
	return newError(ctx, nil, msg, opts)
}

// Wrap creates an error that wraps an existing error, optionally providing additional information.
func Wrap(err error, msg string, opts ...Option) error {
	var ctx context
	return newError(ctx, err, msg, opts)
}

func newError(ctx context, cause error, msg string, opts []Option) *_error {
	ctx = ctx.clone()
	ctx.applyOptions(opts)
	return &_error{
		context: ctx,
		msg:     msg,
		cause:   cause,
	}
}

// Option represents additional information that can be associated
// with an error.
type Option func(*context)

// KV associates a single key-value pair with an error.
func KV(key string, value interface{}) Option {
	return func(ctx *context) {
		kv := keyValueT{
			key:   key,
			value: value,
		}
		ctx.pairs = append(ctx.pairs, kv)
	}
}

func Keyvals(keyvals ...interface{}) Option {
	return func(ctx *context) {
		for i := 0; i < len(keyvals); i += 2 {
			if k, ok := keyvals[i].(string); ok {
				kv := keyValueT{
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
	return func(ctx *context) {
		if pc, file, line, ok := runtime.Caller(skip + 4); ok {
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

// A Context can be useful for specifying key-value pairs to be associated with
// any error messages that might be generated in a function.
type Context interface {
	New(msg string, opts ...Option) error
	Wrap(err error, msg string, opts ...Option) error
	NewContext(opts ...Option) Context
	Caller(skip int) Option
	KV(key string, value interface{}) Option
	Keyvals(keyvals ...interface{}) Option
}

func NewContext(opts ...Option) Context {
	var ctx context
	ctx.applyOptions(opts)
	return &ctx
}

func (ctx *context) New(msg string, opts ...Option) error {
	return newError(*ctx, nil, msg, opts)
}

func (ctx *context) Wrap(err error, msg string, opts ...Option) error {
	return newError(*ctx, err, msg, opts)
}

func (ctx *context) NewContext(opts ...Option) Context {
	ctx2 := ctx.clone()
	ctx2.applyOptions(opts)
	return &ctx2
}

func (ctx *context) Caller(skip int) Option {
	return Caller(skip + 1)
}

func (ctx *context) KV(key string, value interface{}) Option {
	return KV(key, value)
}

func (ctx *context) Keyvals(keyvals ...interface{}) Option {
	return Keyvals(keyvals)
}
