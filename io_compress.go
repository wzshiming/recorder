package recorder

import (
	"compress/zlib"
	"context"
	"fmt"
	"io"
)

var (
	errFormat = fmt.Errorf("invalid foramt")
)

// NewReaderWithCompress return io.ReadCloser that will be save the stream with delay information.
// will compress and need to be closed manually.
func NewReaderWithCompress(ctx context.Context, r io.Reader) (io.ReadCloser, error) {
	zrc, err := zlib.NewReader(r)
	if err != nil {
		if err == zlib.ErrChecksum || err == zlib.ErrDictionary || err == zlib.ErrHeader {
			err = errFormat
		}
		return nil, err
	}
	ctx, cancel := context.WithCancel(ctx)
	r = NewReader(ctx, zrc)
	rc := struct {
		io.Reader
		io.Closer
	}{
		Reader: r,
		Closer: withCloser(func() error {
			err := zrc.Close()
			cancel()
			return err
		}),
	}
	return rc, nil
}

// NewWriterWithCompress return io.WriteCloser that will be reproduce the stream with delay information.
// will compress and need to be closed manually.
func NewWriterWithCompress(ctx context.Context, w io.Writer) io.WriteCloser {
	zwc := zlib.NewWriter(w)
	ctx, cancel := context.WithCancel(ctx)
	w = NewWriter(ctx, zwc)
	wc := struct {
		io.Writer
		io.Closer
	}{
		Writer: w,
		Closer: withCloser(func() error {
			err := zwc.Close()
			cancel()
			return err
		}),
	}
	return wc
}

type withCloser func() error

func (w withCloser) Close() error {
	return w()
}
