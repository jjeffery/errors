package errorv

import (
	"bytes"
	"fmt"

	"github.com/jjeffery/kv"
)

// A Context can be useful for specifying common key/value pairs that will
// be attached to all error messages created from that context.
type Context interface {
	New(msg string, keyvals ...interface{}) error
	Wrap(err error, msg string, keyvals ...interface{}) error
	NewContext(keyvals ...interface{}) Context
}

// NewContext creates a new error context with information that will be
// associated with any errors created from that context.
func NewContext(keyvals ...interface{}) Context {
	var ctx context
	ctx.keyvals = append(ctx.keyvals, keyvals...)
	return ctx
}

// A context implements the public Context interface.
type context struct {
	keyvals []interface{}
}

func (ctx context) New(msg string, keyvals ...interface{}) error {
	return ctx.newError(msg, keyvals)
}

func (ctx context) Wrap(err error, msg string, keyvals ...interface{}) error {
	if msg == "" {
		// A wrap without a message just attaches the options
		// to the error.
		return ctx.attachError(err, keyvals)
	}
	return ctx.wrapError(err, msg, keyvals)
}

func (ctx context) NewContext(keyvals ...interface{}) Context {
	ctx2 := ctx.clone()
	ctx.keyvals = append(ctx.keyvals, keyvals...)
	return ctx2
}

// clone creates a deep copy of the context.
func (ctx context) clone() context {
	ctx2 := context{}
	if len(ctx.keyvals) > 0 {
		// set the capacity of the slice to ensure that any
		// append will allocate a new underlying array
		ctx2.keyvals = ctx.keyvals[0:len(ctx.keyvals):len(ctx.keyvals)]
	}
	return ctx2
}

func (ctx *context) mergeFrom(other context) {
	ctx.keyvals = append(ctx.keyvals, other.keyvals...)
}

func (ctx context) newError(msg string, keyvals []interface{}) error {
	ctx = ctx.clone()
	ctx.keyvals = append(ctx.keyvals, keyvals...)
	return &errorT{
		context: ctx,
		msg:     msg,
	}
}

func (ctx context) wrapError(cause error, msg string, keyvals []interface{}) error {
	if cause == nil {
		return nil
	}
	if msg == "" {
		return ctx.attachError(cause, keyvals)
	}
	ctx = ctx.clone()
	ctx.keyvals = append(ctx.keyvals, keyvals...)
	return &causeT{
		errorT: errorT{
			msg:     msg,
			context: ctx,
		},
		cause: cause,
	}
}

func (ctx context) attachError(cause error, keyvals []interface{}) error {
	if cause == nil {
		return nil
	}
	if len(keyvals) == 0 {
		return cause
	}
	ctx = ctx.clone()
	ctx.keyvals = append(ctx.keyvals, keyvals...)

	type contextGetter interface {
		getContext() *context
	}

	if getContext, ok := cause.(contextGetter); ok {
		// the cause already has a context, so we can
		// append to it
		otherCtx := getContext.getContext()
		otherCtx.mergeFrom(ctx)
		return cause
	}

	// the cause does not have a context, so create an attach
	// wrapper error
	return &attachT{
		context: ctx,
		cause:   cause,
	}
}

func (ctx context) appendKeyvals(keyvals []interface{}) []interface{} {
	return append(keyvals, ctx.keyvals...)
}

func (ctx context) errorBuf(buf *bytes.Buffer) {
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
