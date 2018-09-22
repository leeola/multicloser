package multicloser_test

import (
	"errors"
	"testing"

	"github.com/leeola/multicloser"
)

type mockCloser struct {
	log  *string
	name string
	err  error
}

func (c mockCloser) Close() error {
	if c.log != nil {
		*c.log += c.name
	}
	return c.err
}

func TestCloseOrder(t *testing.T) {
	var log string
	mc := multicloser.New(
		mockCloser{log: &log, name: "a"},
		mockCloser{log: &log, name: "b"},
		mockCloser{log: &log, name: "c"},
	)

	err := mc.Close()
	if err != nil {
		t.Errorf("multicloser returned an error when it shouldn't: %v", err)
	}

	if want := "cba"; log != want {
		t.Errorf("call order is incorrect. got: %q, want: %q", log, want)
	}
}

func TestErrors(t *testing.T) {
	mc := multicloser.New(
		mockCloser{},
		mockCloser{err: errors.New("foo")},
		mockCloser{},
	)

	err := mc.Close()
	if err == nil {
		t.Errorf("multicloser should return an error from closer list, but did not")
	}

	// checking the error string binds this test to the multierror library, but
	// that seems acceptable rather than using the multi-closer lib to match
	// the individual error values. Either we i think we bind to the multierr
	// lib.

	gotErrMsg := err.Error()
	if want := "1 error occurred:\n\t* foo\n\n"; gotErrMsg != want {
		t.Errorf("multicloser error did not match expected error msg. got: %q, want: %q",
			gotErrMsg, want)
	}
}
