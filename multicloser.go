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
func New(c ...io.Closer) MultiCloser {
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

// ReadMultiCloser implements a ReadCloser ontop of multiple closers.
type ReadMultiCloser struct {
	io.Reader
	MultiCloser
}

// NewReadMultiCloser returns a ReadCloser over multiple closers, following
// the same behavior as the MultiCloser.
func NewReadMultiCloser(r io.Reader, c ...io.Closer) ReadMultiCloser {
	return ReadMultiCloser{
		Reader:      r,
		MultiCloser: New(c...),
	}
}

// WriteMultiCloser implements a WriteCloser ontop of multiple closers.
type WriteMultiCloser struct {
	io.Writer
	MultiCloser
}

// NewWriteMultiCloser returns a WriteCloser over multiple closers, following
// the same behavior as the MultiCloser.
func NewWriteMultiCloser(r io.Writer, c ...io.Closer) WriteMultiCloser {
	return WriteMultiCloser{
		Writer:      r,
		MultiCloser: New(c...),
	}
}
