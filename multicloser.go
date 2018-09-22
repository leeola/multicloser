package multicloser

import (
	"io"

	multierror "github.com/hashicorp/go-multierror"
)

// MultiCloser implements closing a slice of io.Closer in reverse order,
// combining encountered errors via a multierror implementation.
type MultiCloser struct {
	Closers []io.Closer
}

// New creates a MultiCloser which closes given io.Closers in reverse
// order.
func New(c ...io.Closer) io.Closer {
	return MultiCloser{
		Closers: c,
	}
}

// Close all closers in reverse order.
func (mc MultiCloser) Close() error {
	var multierr error

	for i := len(mc.Closers) - 1; i >= 0; i-- {
		if err := mc.Closers[i].Close(); err != nil {
			multierr = multierror.Append(err)
		}
	}

	return multierr
}
