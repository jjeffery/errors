package errorv

import (
	"bytes"
	"fmt"
)

var _ = fmt.Printf

// New creates a new error.
func New(msg string, keyvals ...interface{}) error {
	var ctx context
	return ctx.newError(msg, keyvals)
}

// Wrap creates an error that wraps an existing error, optionally providing
// additional information. If err is nil, Wrap returns nil.
func Wrap(err error, msg string, keyvals ...interface{}) error {
	var ctx context
	return ctx.wrapError(err, msg, keyvals)
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
//
// Cause is compatible with the Cause function in package "github.com/pkg/errors".
// The implementation and documentation of Cause has been copied from that package.
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
	buf.WriteString(e.msg)
	e.context.errorBuf(&buf)
	return buf.String()
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
	buf.WriteString(c.msg)
	c.context.errorBuf(&buf)
	buf.WriteString(": ")
	buf.WriteString(c.cause.Error())
	return buf.String()
}

// Cause implements the causer interface, for compatibility with
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
	buf.WriteString(a.cause.Error())
	a.context.errorBuf(&buf)
	return buf.String()
}

// Cause implements the causer interface, for compatibility with
// the github.com/pkg/errors package.
func (a *attachT) Cause() error {
	return a.cause
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
