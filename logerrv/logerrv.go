// Package logerrv provides functions for structured logging.
//
// This package is intended as an interface with
// structured logging packages. It should not be
// used as a mechanism for extracting data from
// an error for further processing by an application.
package logerrv

// Private contains private implementation details shared
// with package errv.
var Private struct {
	Keyvals func(error) (keyvals []interface{}, ok bool)
}

// Keyvals extracts an array of alternating keys and values
// from the error value to be used as part of a structured
// logging message.
//
// The returned keyvals array will have at least one key-value
// pair whose key value is "msg". The following key-value pairs,
// if present, have special meaning.
//
//  msg    Error message
//  cause  Inner error message
//  caller The file and line number of the caller
//
// If the error has been created by the errv package then ok
// is set to true. Otherwise ok is false, and the returned
// keyvals array will only contain the "msg" key-value pair.
func Keyvals(err error) (keyvals []interface{}, ok bool) {
	if err == nil {
		return []interface{}{"msg", "(no error)"}, false
	}
	if Private.Keyvals == nil {
		return []interface{}{"msg", err.Error()}, false
	}
	return Private.Keyvals(err)
}
