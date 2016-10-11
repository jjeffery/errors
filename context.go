package errors

import (
	"bytes"
	"fmt"

	"github.com/jjeffery/kv"
)

// A context implements the public Context interface.
type context struct {
	keyvals []interface{}
}

func (ctx context) New(msg string) Error {
	return ctx.newError(msg)
}

func (ctx context) Wrap(err error, msg string) Error {
	if err == nil {
		return nil
	}
	if msg == "" {
		// A wrap without a message just attaches the options
		// to the error.
		return ctx.attachError(err)
	}
	return ctx.wrapError(err, msg)
}

// Keyvals implements the keyvalser interface.
func (ctx context) Keyvals() []interface{} {
	return ctx.keyvals
}

func (ctx context) With(keyvals ...interface{}) Context {
	return ctx.withKeyvals(keyvals)
}

// safeSlice returns a slice whose capacity is the same as its length.
// This slice is safe for concurrent operations because any attempt to
// append to the slice will result in a new underlying array being allocated.
func safeSlice(keyvals []interface{}) []interface{} {
	if len(keyvals) == 0 {
		return nil
	}
	return keyvals[0:len(keyvals):len(keyvals)]
}

// clone creates a deep copy of the context.
func (ctx context) clone() context {
	return context{
		keyvals: safeSlice(ctx.keyvals),
	}
}

func (ctx context) withKeyvals(keyvals []interface{}) context {
	ctx = ctx.clone()
	ctx.keyvals = append(ctx.keyvals, keyvals...)
	return ctx
}

func (ctx context) newError(msg string) *errorT {
	return &errorT{
		ctx: ctx.clone(),
		msg: msg,
	}
}

func (ctx context) wrapError(cause error, msg string) *causeT {
	return &causeT{
		errorT: &errorT{
			msg: msg,
			ctx: ctx.clone(),
		},
		cause: cause,
	}
}

func (ctx context) attachError(cause error) Error {
	if causeError, ok := cause.(Error); ok {
		// the cause already implements the Error interface,
		// so attach any context key/value pairs directly.
		return causeError.With(ctx.keyvals...)
	}

	// the cause does not implement the Error interface,
	// so create a wrapper and attach the context key/value pairs.
	return &attachT{
		ctx:   ctx.clone(),
		cause: cause,
	}
}

func (ctx context) appendKeyvals(keyvals []interface{}) []interface{} {
	return append(keyvals, ctx.keyvals...)
}

// writeToBuf writes the context's key/value pairs to a buffer.
func (ctx context) writeToBuf(buf *bytes.Buffer) {
	keyvals := kv.Flatten(ctx.keyvals)
	for i := 0; i < len(keyvals); i += 2 {
		// kv.Flatten guarantees that every even-numbered index
		// will contain a string, and that it will be followed by
		// an odd-numbered index
		key := keyvals[i].(string)
		value := keyvals[i+1]

		if buf.Len() > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteString(key)
		buf.WriteRune('=')
		buf.WriteString(fmt.Sprintf("%v", value))
	}
}
