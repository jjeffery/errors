package errv

import (
	"github.com/jjeffery/errv/logerrv"
)

func init() {
	logerrv.Private.Keyvals = logKeyvals
}

func logKeyvals(err error) ([]interface{}, bool) {
	if err == nil {
		return nil, false
	}

	e, ok := err.(*_error)
	if !ok {
		return []interface{}{
			"msg", err.Error(),
		}, false
	}

	var keyvals []interface{}

	if e.msg != "" {
		keyvals = append(keyvals, "msg", e.msg)
	} else if e.cause != nil {
		// This happens when an error is wrapped without a message,
		// eg errv.Wrap(err, "", ...). In this case use the cause message
		// as the error message.
		keyvals = append(keyvals, "msg", e.cause.Error())
	} else {
		// Unlikely to happen, but we want to guarantee the downstream logging library
		// that the first pair will have the key "msg".
		keyvals = append(keyvals, "msg", "(no message)")
	}
	for _, kv := range e.pairs {
		keyvals = append(keyvals, kv.key, kv.value)
	}
	if e.caller != "" {
		keyvals = append(keyvals, "caller", e.caller)
	}
	if e.cause != nil && e.msg != "" {
		keyvals = append(keyvals, "cause", e.cause.Error())
	}
	return keyvals, true
}
