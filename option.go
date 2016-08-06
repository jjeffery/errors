package errorv

import (
	"fmt"
	"runtime"
	"strings"
)

// Option represents additional information that can be associated
// with an error.
type Option func(*context)

// KV associates a single key-value pair with an error.
func KV(key string, value interface{}) Option {
	return func(ctx *context) {
		kv := keyValue{
			key:   key,
			value: value,
		}
		ctx.pairs = append(ctx.pairs, kv)
	}
}

// Keyvals provides a way to associate multiple key-value pairs with the error.
// The keyvals parameter is a variadic sequence of alternating keys and values.
//
// Function KV provides a more typesafe alternative to Keyvals, although
// it is a little more verbose.
func Keyvals(keyvals ...interface{}) Option {
	return func(ctx *context) {
		for i := 0; i < len(keyvals); i += 2 {
			if k, ok := keyvals[i].(string); ok {
				kv := keyValue{
					key:   k,
					value: keyvals[i+1],
				}
				ctx.pairs = append(ctx.pairs, kv)
			}
		}
	}
}

// Caller associates the file and line number of the caller with
// the error. The argument skip is the number of stack frames to
// ascend, with 0 identifying the caller of Caller.
func Caller(skip int) Option {
	// additionalSkip is the number of stack frames used by
	// this package in a call to the function returned by
	// this function. It needs to be added to the number
	// if skip frames requested by the calling program.
	const additionalSkip = 4

	return func(ctx *context) {
		if pc, file, line, ok := runtime.Caller(skip + additionalSkip); ok {
			fn := runtime.FuncForPC(pc)
			file = trimGOPATH(fn.Name(), file)
			ctx.caller = fmt.Sprintf("%s:%d", file, line)
		}
	}
}

// trimGOPATH was copied from https://github.com/pkg/errors (Author: Dave Cheney)
// which in turn was adapted from https://github.com/go-stack/stack (Author: Chris Hines).
func trimGOPATH(name, file string) string {
	// Here we want to get the source file path relative to the compile time
	// GOPATH. As of Go 1.6.x there is no direct way to know the compiled
	// GOPATH at runtime, but we can infer the number of path segments in the
	// GOPATH. We note that fn.Name() returns the function name qualified by
	// the import path, which does not include the GOPATH. Thus we can trim
	// segments from the beginning of the file path until the number of path
	// separators remaining is one more than the number of path separators in
	// the function name. For example, given:
	//
	//    GOPATH     /home/user
	//    file       /home/user/src/pkg/sub/file.go
	//    fn.Name()  pkg/sub.Type.Method
	//
	// We want to produce:
	//
	//    pkg/sub/file.go
	//
	// From this we can easily see that fn.Name() has one less path separator
	// than our desired output. We count separators from the end of the file
	// path until it finds two more than in the function name and then move
	// one character forward to preserve the initial path segment without a
	// leading separator.
	const sep = "/"
	goal := strings.Count(name, sep) + 2
	i := len(file)
	for n := 0; n < goal; n++ {
		i = strings.LastIndex(file[:i], sep)
		if i == -1 {
			// not enough separators found, set i so that the slice expression
			// below leaves file unmodified
			i = -len(sep)
			break
		}
	}
	// get back to 0 or trim the leading separator
	file = file[i+len(sep):]
	return file
}
