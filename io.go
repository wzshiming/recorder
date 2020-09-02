package recorder

import (
	"context"
	"io"
	"time"
)

type reader struct {
	ctx    context.Context
	buf    []byte
	reader io.Reader
	sum    time.Duration
	start  time.Time
}

// NewReader return io.Redaer that will be reproduce the stream with delay information.
func NewReader(ctx context.Context, r io.Reader) io.Reader {
	return &reader{
		ctx:    ctx,
		reader: r,
		start:  time.Now(),
	}
}

func (r *reader) Read(p []byte) (n int, err error) {
	err = r.ctx.Err()
	if err != nil {
		return 0, err
	}
	if len(r.buf) == 0 {
		d := record{}
		err := d.decode(r.reader)
		if err != nil {
			return 0, err
		}
		r.buf = d.message
		r.sum += d.duration
		duration := r.sum - time.Since(r.start)
		if duration > 0 {
			select {
			case <-r.ctx.Done():
				return 0, r.ctx.Err()
			case <-time.After(duration):
			}
		}
	}
	n = copy(p, r.buf)
	r.buf = r.buf[n:]
	return n, nil
}

type writer struct {
	ctx    context.Context
	last   time.Time
	writer io.Writer
}

// NewWriter return io.Writer that will be save the stream with delay information.
func NewWriter(ctx context.Context, w io.Writer) io.Writer {
	return &writer{
		ctx:    ctx,
		last:   time.Now(),
		writer: w,
	}
}

func (w *writer) Write(p []byte) (n int, err error) {
	err = w.ctx.Err()
	if err != nil {
		return 0, err
	}
	now := time.Now()
	ra := record{
		duration: now.Sub(w.last),
		message:  p,
	}
	err = ra.encode(w.writer)
	if err != nil {
		return 0, err
	}
	w.last = now
	return len(p), nil
}
