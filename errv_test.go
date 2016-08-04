package errv

import (
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
			expect: "msg github.com/jjeffery/errv/errv_test.go:51",
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
