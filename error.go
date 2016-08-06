package errorv

import (
	"bytes"
	"fmt"
)

// New creates a new error.
func New(msg string, opts ...Option) error {
	var ctx context
	return ctx.newError(msg, opts)
}

// Wrap creates an error that wraps an existing error, optionally providing additional information.
// If err is nil, Wrap returns nil.
func Wrap(err error, msg string, opts ...Option) error {
	var ctx context
	return ctx.wrapError(err, msg, opts)
}

// Cause was copied from https://github.com/pkg/errors
// for compatibility.

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

type errorT struct {
	context
	msg string
}

// Error implements the error interface.
func (e *errorT) Error() string {
	var buf bytes.Buffer
	e.errorBuf(&buf)
	return buf.String()
}

// errorBuf fills a buffer with text for an error message.
// Shared with causeT Error implementation.
func (e *errorT) errorBuf(buf *bytes.Buffer) {
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
}

// Keyvals returns the contents of the error
// as an array of alternating keys and values.
func (e *errorT) Keyvals() []interface{} {
	var keyvals []interface{}

	if e.msg != "" {
		keyvals = append(keyvals, "msg", e.msg)
	} else {
		keyvals = append(keyvals, "msg", "(no message)")
	}
	for _, kv := range e.pairs {
		keyvals = append(keyvals, kv.key, kv.value)
	}
	if e.caller != "" {
		keyvals = append(keyvals, "caller", e.caller)
	}
	return keyvals
}

type causeT struct {
	errorT
	cause error
}

// Error implements the error interface.
func (c *causeT) Error() string {
	var buf bytes.Buffer
	c.errorBuf(&buf)
	buf.WriteRune(':')
	buf.WriteRune(' ')
	buf.WriteString(c.cause.Error())
	return buf.String()
}

// Cause implements the causer interface, for compatiblity with
// the github.com/pkg/errors package.
func (c *causeT) Cause() error {
	return c.cause
}

// Keyvals returns the contents of the error
// as an array of alternating keys and values.
func (c *causeT) Keyvals() []interface{} {
	keyvals := c.errorT.Keyvals()
	keyvals = append(keyvals, "cause", c.cause.Error())
	return keyvals
}
