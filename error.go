package errorv

import (
	"bytes"
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

// Attach attaches additional context to an existing error.
func Attach(err error, opts ...Option) error {
	var ctx context
	return ctx.attachError(err, opts)
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
	e.context.errorBuf(buf)
}

// Keyvals returns the contents of the error
// as an array of alternating keys and values.
func (e *errorT) Keyvals() []interface{} {
	var keyvals []interface{}
	keyvals = append(keyvals, "msg", e.msg)
	keyvals = e.appendKeyvals(keyvals)
	return keyvals
}

// getContext implements contextGetter interface
func (e *errorT) getContext() *context {
	return &e.context
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

// getContext implements contextGetter interface
func (c *causeT) getContext() *context {
	return &c.context
}

type attachT struct {
	context
	cause error
}

// Error implements the error interface.
func (a *attachT) Error() string {
	var buf bytes.Buffer
	a.errorBuf(&buf)
	return buf.String()
}

// Cause implements the causer interface, for compatiblity with
// the github.com/pkg/errors package.
func (a *attachT) Cause() error {
	return a.cause
}

// errorBuf fills a buffer with text for an error message.
// Shared with causeT Error implementation.
func (a *attachT) errorBuf(buf *bytes.Buffer) {
	buf.WriteString(a.cause.Error())
	a.context.errorBuf(buf)
}

// Keyvals returns the contents of the error
// as an array of alternating keys and values.
func (a *attachT) Keyvals() []interface{} {
	var keyvals []interface{}
	keyvals = append(keyvals, "msg", a.cause.Error())
	keyvals = a.appendKeyvals(keyvals)
	return keyvals
}

// getContext implements contextGetter interface
func (a *attachT) getContext() *context {
	return &a.context
}
