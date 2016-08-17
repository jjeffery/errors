package errorv

import (
	"bytes"
	"fmt"
)

// A Context can be useful for specifying key-value pairs to be associated with
// any error messages that might be generated in a function.
type Context interface {
	New(msg string, opts ...Option) error
	Wrap(err error, msg string, opts ...Option) error
	Attach(err error, opts ...Option) error
	NewContext(opts ...Option) Context
	Caller(skip int) Option
	KV(key string, value interface{}) Option
	Keyvals(keyvals ...interface{}) Option
}

// NewContext creates a new error context with information that will be
// associated with any errors created from that context.
func NewContext(opts ...Option) Context {
	var ctx context
	ctx.applyOptions(opts)
	return ctx
}

// keyValue is a key-value pair containing additional information
// relevant to an error.
type keyValue struct {
	key   string
	value interface{}
}

// A context implements the public Context interface.
type context struct {
	pairs  []keyValue
	caller string
}

func (ctx context) New(msg string, opts ...Option) error {
	return ctx.newError(msg, opts)
}

func (ctx context) Wrap(err error, msg string, opts ...Option) error {
	if msg == "" {
		// A wrap without a message just attaches the options
		// to the error.
		return ctx.attachError(err, opts)
	}
	return ctx.wrapError(err, msg, opts)
}

func (ctx context) Attach(err error, opts ...Option) error {
	return ctx.attachError(err, opts)
}

func (ctx context) NewContext(opts ...Option) Context {
	ctx2 := ctx.clone()
	ctx2.applyOptions(opts)
	return ctx2
}

func (ctx context) Caller(skip int) Option {
	return Caller(skip + 1)
}

func (ctx context) KV(key string, value interface{}) Option {
	return KV(key, value)
}

func (ctx context) Keyvals(keyvals ...interface{}) Option {
	return Keyvals(keyvals)
}

// clone creates a deep copy of the context.
func (ctx context) clone() context {
	ctx2 := context{
		caller: ctx.caller,
	}
	if len(ctx.pairs) > 0 {
		ctx2.pairs = make([]keyValue, len(ctx.pairs))
		copy(ctx2.pairs, ctx.pairs)
	}
	return ctx2
}

func (ctx *context) applyOptions(opts []Option) {
	for _, opt := range opts {
		opt(ctx)
	}
}

func (ctx *context) mergeFrom(other context) {
	for _, kv := range other.pairs {
		// TODO: check for duplicates
		ctx.pairs = append(ctx.pairs, kv)
	}

	// TODO: fix
	if ctx.caller == "" {
		ctx.caller = other.caller
	}
}

func (ctx context) newError(msg string, opts []Option) error {
	ctx = ctx.clone()
	ctx.applyOptions(opts)
	return &errorT{
		context: ctx,
		msg:     msg,
	}
}

func (ctx context) wrapError(cause error, msg string, opts []Option) error {
	if cause == nil {
		return nil
	}
	ctx = ctx.clone()
	ctx.applyOptions(opts)
	return &causeT{
		errorT: errorT{
			msg:     msg,
			context: ctx,
		},
		cause: cause,
	}
}

func (ctx context) attachError(cause error, opts []Option) error {
	if cause == nil {
		return nil
	}
	if len(opts) == 0 {
		return cause
	}
	ctx = ctx.clone()
	ctx.applyOptions(opts)

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
	for _, kv := range ctx.pairs {
		keyvals = append(keyvals, kv.key, kv.value)
	}
	if ctx.caller != "" {
		keyvals = append(keyvals, "caller", ctx.caller)
	}
	return keyvals
}

func (ctx context) errorBuf(buf *bytes.Buffer) {
	for _, kv := range ctx.pairs {
		if buf.Len() > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteString(kv.key)
		buf.WriteRune('=')
		buf.WriteString(fmt.Sprintf("%v", kv.value))
	}
	if ctx.caller != "" {
		if buf.Len() > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteString(ctx.caller)
	}
}
