package errors

import (
	"bytes"
	"encoding"
	"fmt"
	"reflect"
	"strings"

	"github.com/jjeffery/kv"
)

// A context implements the public Context interface.
type context struct {
	keyvals []interface{}
}

// New creates a new context.
func (ctx context) New(msg string) Error {
	return ctx.newError(msg)
}

func (ctx context) Wrap(err error, msg ...string) Error {
	if err == nil {
		return nil
	}
	// strip out any empty strings in the msg slice
	{
		v := make([]string, 0, len(msg))
		for _, m := range msg {
			if m != "" {
				v = append(v, m)
			}
		}
		msg = v
	}
	if len(msg) == 0 {
		// A wrap without a message just attaches the options
		// to the error.
		return ctx.attachError(err)
	}
	return ctx.wrapError(err, strings.Join(msg, ": "))
}

// Keyvals implements the keyvalser interface.
func (ctx context) Keyvals() []interface{} {
	return ctx.keyvals
}

func (ctx context) With(keyvals ...interface{}) Context {
	return ctx.withKeyvals(keyvals)
}

// safeSlice returns a slice whose capacity is the same as its length.
// This slice is safe for concurrent operations because any attempt to
// append to the slice will result in a new underlying array being allocated.
func safeSlice(keyvals []interface{}) []interface{} {
	if len(keyvals) == 0 {
		return nil
	}
	return keyvals[0:len(keyvals):len(keyvals)]
}

// clone creates a deep copy of the context.
func (ctx context) clone() context {
	return context{
		keyvals: safeSlice(ctx.keyvals),
	}
}

func (ctx context) withKeyvals(keyvals []interface{}) context {
	ctx = ctx.clone()
	ctx.keyvals = append(ctx.keyvals, keyvals...)
	return ctx
}

func (ctx context) newError(msg string) *errorT {
	return &errorT{
		ctx: ctx.clone(),
		msg: msg,
	}
}

func (ctx context) wrapError(cause error, msg string) *causeT {
	return &causeT{
		errorT: &errorT{
			msg: msg,
			ctx: ctx.clone(),
		},
		cause: cause,
	}
}

func (ctx context) attachError(cause error) Error {
	return &attachT{
		ctx:   ctx.clone(),
		cause: cause,
	}
}

func (ctx context) appendKeyvals(keyvals []interface{}) []interface{} {
	return append(keyvals, ctx.keyvals...)
}

// writeToBuf writes the context's key/value pairs to a buffer.
func (ctx context) writeToBuf(buf *bytes.Buffer) {
	keyvals := kv.Flatten(ctx.keyvals)
	for i := 0; i < len(keyvals); i += 2 {
		// kv.Flatten guarantees that every even-numbered index
		// will contain a string, and that it will be followed by
		// an odd-numbered index
		key := keyvals[i].(string)
		value := keyvals[i+1]

		if buf.Len() > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteString(key)
		buf.WriteRune('=')
		writeValue(buf, value)
	}
}

// constant byte values
var (
	bytesNull  = []byte("null")
	bytesPanic = []byte(`"<PANIC>"`)
	bytesError = []byte(`"<ERROR>"`)
)

func writeValue(buf *bytes.Buffer, value interface{}) {
	defer func() {
		if r := recover(); r != nil {
			if buf != nil {
				buf.Write(bytesPanic)
			}
		}
	}()
	switch v := value.(type) {
	case nil:
		writeBytesValue(buf, bytesNull)
		return
	case []byte:
		writeBytesValue(buf, v)
		return
	case string:
		writeStringValue(buf, v)
		return
	case bool, byte, int8, int16, uint16, int32, uint32, int64, uint64, int, uint, uintptr, float32, float64, complex64, complex128:
		fmt.Fprint(buf, v)
		return
	case encoding.TextMarshaler:
		writeTextMarshalerValue(buf, v)
		return
	case error:
		writeStringValue(buf, v.Error())
		return
	case fmt.Stringer:
		writeStringValue(buf, v.String())
		return
	default:
		// handle pointer to any of the above
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				buf.Write(bytesNull)
				return
			}
			writeValue(buf, rv.Elem().Interface())
			return
		}
		writeStringValue(buf, fmt.Sprint(value))
	}
}

func writeBytesValue(buf *bytes.Buffer, b []byte) {
	if b == nil {
		buf.Write(bytesNull)
		return
	}
	index := bytes.IndexFunc(b, needsQuote)
	if index < 0 {
		buf.Write(b)
		return
	}
	buf.WriteRune('"')
	if index > 0 {
		buf.Write(b[0:index])
		b = b[index:]
	}
	for {
		index = bytes.IndexFunc(b, needsBackslash)
		if index < 0 {
			break
		}
		if index > 0 {
			buf.Write(b[:index])
			b = b[index:]
		}
		buf.WriteRune('\\')
		// we know that the rune will be a single byte
		buf.WriteByte(b[0])
		b = b[1:]
	}
	buf.Write(b)
	buf.WriteRune('"')
}

func writeStringValue(buf *bytes.Buffer, s string) {
	index := strings.IndexFunc(s, needsQuote)
	if index < 0 {
		buf.WriteString(s)
		return
	}
	buf.WriteRune('"')
	if index > 0 {
		buf.WriteString(s[0:index])
		s = s[index:]
	}
	for {
		index = strings.IndexFunc(s, needsBackslash)
		if index < 0 {
			break
		}
		if index > 0 {
			buf.WriteString(s[0:index])
			s = s[index:]
		}
		buf.WriteRune('\\')
		// we know that the rune will be a single byte
		buf.WriteByte(s[0])
		s = s[1:]
	}
	buf.WriteString(s)
	buf.WriteRune('"')
}

func writeTextMarshalerValue(buf *bytes.Buffer, t encoding.TextMarshaler) {
	b, err := t.MarshalText()
	if err != nil {
		buf.Write(bytesError)
		return
	}
	writeBytesValue(buf, b)
}

func needsQuote(c rune) bool {
	// the single quote '\'' is not strictly necessary, but
	// is more human readable if quoted
	return c <= ' ' || c == '"' || c == '\\' || c == '\''
}

func needsBackslash(c rune) bool {
	return c == '\\' || c == '"'
}
