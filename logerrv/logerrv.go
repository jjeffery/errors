// Package logerrv provides functions for structured logging.
//
// This package is intended as an interface with
// structured logging packages. It should not be
// used as a mechanism for extracting data from
// an error for further processing.
package logerrv

// Private contains private implementation details shared
// with package errv.
var Private struct {
	Keyvals func(error) (keyvals []interface{}, ok bool)
	Map     func(error) (msg string, fields map[string]interface{}, ok bool)
}

// Keyvals returns an array of alternating keys and values.
// Useful for interfacing with gokit log package.
func Keyvals(err error) []interface{} {
	if Private.Keyvals != nil {
		if v, ok := Private.Keyvals(err); ok {
			return v
		}
	}
	return []interface{}{
		"msg",
		err.Error(),
	}
}

// Map returns a message and a map of key-values.
// Useful for interfacing with logrus.
func Map(err error) (msg string, fields map[string]interface{}) {
	if Private.Map != nil {
		if msg, fields, ok := Private.Map(err); ok {
			return msg, fields
		}
	}
	return err.Error(), nil
}
