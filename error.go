package errors

import (
	"bytes"
)

// errorT represents an error with a message and context.
type errorT struct {
	ctx context
	msg string
}

// Error implements the error interface.
func (e *errorT) Error() string {
	var buf bytes.Buffer
	buf.WriteString(e.msg)
	e.ctx.writeToBuf(&buf)
	return buf.String()
}

// With returns an error with additional key/value pairs attached.
// It implements the Error interface.
func (e *errorT) With(keyvals ...interface{}) Error {
	return e.withKeyvals(keyvals)
}

// MarshalText implements the TextMarshaler interface.
func (e *errorT) MarshalText() ([]byte, error) {
	return []byte(e.Error()), nil
}

// Keyvals returns the contents of the error
// as an array of alternating keys and values.
func (e *errorT) Keyvals() []interface{} {
	var keyvals []interface{}
	keyvals = append(keyvals, "msg", e.msg)
	keyvals = e.ctx.appendKeyvals(keyvals)
	return keyvals
}

func (e *errorT) withKeyvals(keyvals []interface{}) *errorT {
	return &errorT{
		ctx: e.ctx.withKeyvals(keyvals),
		msg: e.msg,
	}
}

// causeT represents an error with a message, context, and an error which
// contains the original cause of the error condition.
type causeT struct {
	*errorT
	cause error
}

// Error implements the error interface.
func (c *causeT) Error() string {
	var buf bytes.Buffer
	buf.WriteString(c.msg)
	c.ctx.writeToBuf(&buf)
	buf.WriteString(": ")
	buf.WriteString(c.cause.Error())
	return buf.String()
}

// With returns an error with additional key/value pairs attached.
// It implements the Error interface.
func (c *causeT) With(keyvals ...interface{}) Error {
	return &causeT{
		errorT: c.errorT.withKeyvals(keyvals),
		cause:  c.cause,
	}
}

// MarshalText implements the TextMarshaler interface.
func (c *causeT) MarshalText() ([]byte, error) {
	return []byte(c.Error()), nil
}

// Cause implements the causer interface, and is compatible with
// the github.com/pkg/errors package.
func (c *causeT) Cause() error {
	return c.cause
}

// Keyvals returns the contents of the error
// as an array of alternating keys and values.
func (c *causeT) Keyvals() []interface{} {
	keyvals := c.errorT.Keyvals()

	// TODO(jpj): this might be improved by checking if cause
	// implements keyvalser, and appending keyvals.
	keyvals = append(keyvals, "cause", c.cause.Error())
	return keyvals
}

// attachT represents an error that has additional keyword/value pairs
// attached to it.
type attachT struct {
	ctx   context
	cause error
}

// Error implements the error interface.
func (a *attachT) Error() string {
	var buf bytes.Buffer
	buf.WriteString(a.cause.Error())
	a.ctx.writeToBuf(&buf)
	return buf.String()
}

// With returns an error with additional key/value pairs attached.
// It implements the Error interface.
func (a *attachT) With(keyvals ...interface{}) Error {
	return &attachT{
		ctx:   a.ctx.withKeyvals(keyvals),
		cause: a.cause,
	}
}

// MarshalText implements the TextMarshaler interface.
func (a *attachT) MarshalText() ([]byte, error) {
	return []byte(a.Error()), nil
}

// Cause implements the causer interface, and is compatible with
// the github.com/pkg/errors package.
func (a *attachT) Cause() error {
	return a.cause
}

// Keyvals returns the contents of the error
// as an array of alternating keys and values.
func (a *attachT) Keyvals() []interface{} {
	var keyvals []interface{}
	// TODO(jpj): this could be improved by checking if the
	// cause implements the keyvalser interface.
	keyvals = append(keyvals, "msg", a.cause.Error())
	keyvals = a.ctx.appendKeyvals(keyvals)
	return keyvals
}
