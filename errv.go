// Package errv provides a simple interface for
// error handling that works well with structured logging.
package errv

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// New creates a new error with fields.
func New(msg string, fields ...Field) error {
	// TODO(jpj): implement
	return errors.New(flatten(msg, fields))
}

// Wrap creates an error that wraps an existing error, providing additional fields.
func Wrap(err error, msg string, fields ...Field) error {
	// TODO(jpj): implement this
	return errors.Wrap(err, flatten(msg, fields))
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

// V returns a field.
func V(key string, value interface{}) Field {
	return Field{
		key:   key,
		value: value,
	}
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

func flatten(msg string, fields []Field) string {
	var vals []string
	if msg != "" {
		vals = append(vals, msg)
	}
	for _, field := range fields {
		vals = append(vals, fmt.Sprintf("%s=%v", field.Key(), field.Value()))
	}
	return strings.Join(vals, " ")
}
