package utils

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

type limitedWriter struct {
	writer   io.Writer
	capacity uint64
}

var ErrWriteLimitExceeded = errors.New("write limit exceeded")

func (lw *limitedWriter) Write(p []byte) (n int, err error) {
	if lw.capacity <= 0 {
		return 0, ErrWriteLimitExceeded
	}

	if uint64(len(p)) > lw.capacity {
		p = p[:lw.capacity]
		err = ErrWriteLimitExceeded
	}

	n, writeErr := lw.writer.Write(p)
	lw.capacity -= uint64(n)
	if writeErr != nil {
		return n, writeErr
	}
	return n, err
}

func (lw *limitedWriter) String() string {
	return strings.TrimSpace(lw.writer.(*bytes.Buffer).String())
}

var _ io.Writer = (*limitedWriter)(nil)

func NewLimitedWriter(capacity int) *limitedWriter {
	var buf bytes.Buffer
	return &limitedWriter{
		writer:   &buf,
		capacity: uint64(capacity) * 1024,
	}
}
