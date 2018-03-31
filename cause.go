package errors

// Cause was copied from https://github.com/pkg/errors
// for compatibility. See CREDITS.md.

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
		if e := cause.Cause(); e != nil {
			err = e
		} else {
			break
		}
	}
	return err
}
