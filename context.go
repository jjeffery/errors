package errors

import (
	"bytes"
	"strings"

	"github.com/jjeffery/kv"
)

// A context implements the public Context interface.
type context struct {
	keyvals []interface{}
}

// New creates a new context.
func (ctx context) New(msg string) Error {
	return ctx.newError(msg)
}

func (ctx context) Wrap(err error, msg ...string) Error {
	if err == nil {
		return nil
	}
	// strip out any empty strings in the msg slice
	{
		v := make([]string, 0, len(msg))
		for _, m := range msg {
			if m != "" {
				v = append(v, m)
			}
		}
		msg = v
	}
	if len(msg) == 0 {
		// A wrap without a message just attaches the options
		// to the error.
		return ctx.attachError(err)
	}
	return ctx.wrapError(err, strings.Join(msg, ": "))
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
	if len(ctx.keyvals) == 0 {
		return
	}
	// kv.List.MarshalText does not return a non-nil error.
	b, _ := kv.List(ctx.keyvals).MarshalText()
	if buf.Len() > 0 {
		buf.WriteRune(' ')
	}
	buf.Write(b)
}
