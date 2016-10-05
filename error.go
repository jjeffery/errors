package errors

import (
	"bytes"
	"fmt"
)

var _ = fmt.Printf

type errorT struct {
	ctx context
	msg string
}

// Error implements the error interface.
func (e *errorT) Error() string {
	var buf bytes.Buffer
	buf.WriteString(e.msg)
	e.ctx.errorBuf(&buf)
	return buf.String()
}

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

type causeT struct {
	*errorT
	cause error
}

// Error implements the error interface.
func (c *causeT) Error() string {
	var buf bytes.Buffer
	buf.WriteString(c.msg)
	c.ctx.errorBuf(&buf)
	buf.WriteString(": ")
	buf.WriteString(c.cause.Error())
	return buf.String()
}

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

// Cause implements the causer interface, for compatibility with
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

type attachT struct {
	ctx   context
	cause error
}

// Error implements the error interface.
func (a *attachT) Error() string {
	var buf bytes.Buffer
	buf.WriteString(a.cause.Error())
	a.ctx.errorBuf(&buf)
	return buf.String()
}

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

// Cause implements the causer interface, for compatibility with
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
