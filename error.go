package errorv

import (
	"bytes"
	"fmt"
)

type errorT struct {
	context
	msg   string
	cause error
}

// Error implements the error interface.
func (e *errorT) Error() string {
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
	var ctx context
	return ctx.newError(nil, msg, opts)
}

// Wrap creates an error that wraps an existing error, optionally providing additional information.
func Wrap(err error, msg string, opts ...Option) error {
	if err == nil {
		return nil
	}
	var ctx context
	return ctx.newError(err, msg, opts)
}

// Cause implements the causer interface, for compatiblity with
// the github.com/pkg/errors package.
func (e *errorT) Cause() error {
	return e.cause
}

// Keyvals returns the contents of the error
// as an array of alternating keys and values.
func (e *errorT) Keyvals() []interface{} {
	var keyvals []interface{}

	if e.msg != "" {
		keyvals = append(keyvals, "msg", e.msg)
	} else if e.cause != nil {
		// This happens when an error is wrapped without a message,
		// eg errorv.Wrap(err, "", ...). In this case use the cause message
		// as the error message.
		keyvals = append(keyvals, "msg", e.cause.Error())
	} else {
		// Unlikely to happen, but we want to guarantee the downstream logging library
		// that the first pair will have the key "msg".
		keyvals = append(keyvals, "msg", "(no message)")
	}
	for _, kv := range e.pairs {
		keyvals = append(keyvals, kv.key, kv.value)
	}
	if e.caller != "" {
		keyvals = append(keyvals, "caller", e.caller)
	}
	if e.cause != nil && e.msg != "" {
		keyvals = append(keyvals, "cause", e.cause.Error())
	}
	return keyvals
}
