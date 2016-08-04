package errv

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
	return ctx.newError(nil, msg, opts)
}

func (ctx context) Wrap(err error, msg string, opts ...Option) error {
	return ctx.newError(err, msg, opts)
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

func (ctx context) newError(err error, msg string, opts []Option) error {
	ctx = ctx.clone()
	ctx.applyOptions(opts)
	return &errorT{
		context: ctx,
		msg:     msg,
		cause:   err,
	}
}
