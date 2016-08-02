// Package errv provides a simple interface for
// error handling that works well with structured logging.
package errv

import (
	"errors"
)

// New creates a new error with fields.
func New(msg string, fields ...Field) error {
	// TODO(jpj): implement this
	return errors.New(msg)
}

// Wrap creates an error that wraps an existing error, providing additional fields.
func Wrap(err error, msg string, fields ...Field) error {
	// TODO(jpj): implement this
	return errors.New(msg)
}

// Add fields to an existing error, leaving the message unchanged.
//
// A field is not added to the error if an identical field already exists.
func Add(err error, fields ...Field) error {
	// TODO(jpj): implement this
	return err
}

// Field contains additional information about an error.
type Field struct {
	key   string
	value interface{}
}

// Key returns the key for the field.
func (f Field) Key() string {
	return f.key
}

// Value returns the field value.
func (f Field) Value() interface{} {
	return f.value
}

// String returns a field with a string value.
func String(key string, value string) Field {
	return Field{
		key:   key,
		value: value,
	}
}

// Int returns a field with an integer value.
func Int(key string, value int) Field {
	return Field{
		key:   key,
		value: value,
	}
}
