package errorv

import (
	"io"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		msg    string
		opts   []Option
		expect string
	}{
		{
			msg:    "",
			opts:   nil,
			expect: "",
		},
		{
			msg:    "xx",
			opts:   nil,
			expect: "xx",
		},
		{
			msg: "xx",
			opts: []Option{
				KV("key1", "val1"),
				KV("key2", 2),
			},
			expect: "xx key1=val1 key2=2",
		},
		{
			msg: "msg",
			opts: []Option{
				KV("key", time.Time{}),
			},
			expect: "msg key=0001-01-01 00:00:00 +0000 UTC",
		},
		{
			msg: "msg",
			opts: []Option{
				Caller(0),
			},
			// WARNING: this test is pretty brittle: if you move
			// any lines in this file you will have to change expect.
			expect: "msg github.com/jjeffery/errorv/errorv_test.go:53",
		},
	}

	for _, tt := range tests {
		got := New(tt.msg, tt.opts...)
		if got.Error() != tt.expect {
			t.Errorf("New.Error(): got: %q, want %q", got, tt.expect)
		}
	}
}

func TestWrapNil(t *testing.T) {
	got := Wrap(nil, "no error")
	if got != nil {
		t.Errorf("Wrap(nil, \"no error\"): got %#v, expected nil", got)
	}
}

type nilError struct{}

func (nilError) Error() string { return "nil error" }

func TestCause(t *testing.T) {
	x := New("error")
	tests := []struct {
		err  error
		want error
	}{{
		// nil error is nil
		err:  nil,
		want: nil,
	}, {
		// explicit nil error is nil
		err:  (error)(nil),
		want: nil,
	}, {
		// typed nil is nil
		err:  (*nilError)(nil),
		want: (*nilError)(nil),
	}, {
		// uncaused error is unaffected
		err:  io.EOF,
		want: io.EOF,
	}, {
		// caused error returns cause
		err:  Wrap(io.EOF, "ignored"),
		want: io.EOF,
	}, {
		err:  x, // return from errors.New
		want: x,
	}}

	for i, tt := range tests {
		got := Cause(tt.err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("test %d: got %#v, want %#v", i+1, got, tt.want)
		}
	}
}

func TestAttachNil(t *testing.T) {
	got := Attach(nil, KV("no error", "no error"))
	if got != nil {
		t.Errorf("Attach(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestAttach(t *testing.T) {
	tests := []struct {
		cause         error
		opts          []Option
		expectedMsg   string
		expectedCause error
		expectedErr   error
	}{
		{
			// this test case tests that when no options are passed, the
			// original error is returned
			cause:         io.EOF,
			opts:          nil,
			expectedMsg:   "EOF",
			expectedCause: io.EOF,
			expectedErr:   io.EOF,
		},
		{
			cause:         io.EOF,
			opts:          []Option{KV("k1", "v1"), KV("k2", "v2")},
			expectedMsg:   "EOF k1=v1 k2=v2",
			expectedCause: io.EOF,
		},
		{
			cause:         Wrap(io.EOF, "something failed", KV("k3", "v3")),
			opts:          []Option{KV("k1", "v1"), KV("k2", "v2")},
			expectedMsg:   "something failed k3=v3 k1=v1 k2=v2: EOF",
			expectedCause: io.EOF,
		},
	}

	for i, tt := range tests {
		err := Attach(tt.cause, tt.opts...)
		actualMsg := err.Error()
		if actualMsg != tt.expectedMsg {
			t.Errorf("%d: expected=%q, actual=%q", i, tt.expectedMsg, actualMsg)
		}
		actualCause := Cause(err)
		if actualCause != tt.expectedCause {
			t.Errorf("%d: cause: expected=%v, actual=%v", i, tt.expectedCause, actualCause)
		}

		// only test if non-nil in the test case
		if tt.expectedErr != nil {
			if tt.expectedErr != err {
				t.Errorf("%d: error: expected=%v, actual=%v", i, tt.expectedErr, err)
			}
		}
	}
}

func TestKeyvals(t *testing.T) {
	tests := []struct {
		err     error
		keyvals []interface{}
	}{
		{
			err:     New("message"),
			keyvals: []interface{}{"msg", "message"},
		},
		{
			err:     New("message", KV("k1", "v1"), KV("k2", 2)),
			keyvals: []interface{}{"msg", "message", "k1", "v1", "k2", 2},
		},
		{
			err:     Wrap(io.EOF, "message", KV("k1", "v1"), KV("k2", 2)),
			keyvals: []interface{}{"msg", "message", "k1", "v1", "k2", 2, "cause", "EOF"},
		},
		{
			err:     Attach(io.EOF, KV("k1", "v1"), KV("k2", 2)),
			keyvals: []interface{}{"msg", "EOF", "k1", "v1", "k2", 2},
		},
	}

	type keyvalser interface {
		Keyvals() []interface{}
	}

	for i, tt := range tests {
		keyvals, ok := tt.err.(keyvalser)
		if !ok {
			t.Errorf("%d: expected Keyvals(), none available", i)
			continue
		}
		kvs := keyvals.Keyvals()
		if !reflect.DeepEqual(tt.keyvals, kvs) {
			t.Errorf("%d: expected %v, actual %v", i, tt.keyvals, kvs)
		}
	}
}
