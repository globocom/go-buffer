package buffer

import (
	"errors"
	"io"
	"time"
)

const (
	invalidSize    = "go-buffer: size must be greater than zero"
	invalidFlusher = "go-buffer: flusher is required"
)

var (
	// ErrTimeout indicates an operation has timed out.
	ErrTimeout = errors.New("operation timed-out")
)

type (
	// Buffer represents a data buffer that is asynchronously flushed, either manually or automatically.
	Buffer struct {
		io.Closer
		dataCh  chan interface{}
		flushCh chan struct{}
		doneCh  chan struct{}
		size    uint
		flusher func([]interface{})
		options *Options
	}
)

// Push appends an item to the end of the buffer. It times out if it cannot be
// performed in a timely fashion.
func (buffer *Buffer) Push(item interface{}) error {
	select {
	case buffer.dataCh <- item:
		return nil
	case <-time.After(buffer.options.PushTimeout):
		return ErrTimeout
	}
}

// Flush outputs the buffer to a permanent destination. It times out if it cannot be
// performed in a timely fashion.
func (buffer *Buffer) Flush() error {
	select {
	case buffer.flushCh <- struct{}{}:
		return nil
	case <-time.After(buffer.options.FlushTimeout):
		return ErrTimeout
	}
}

// Close flushes the buffer and prevents it from being further used. It times
// out if it cannot be performed in a timely fashion.
// The buffer must not be used after it has been closed as all further operations will panic.
func (buffer *Buffer) Close() error {
	close(buffer.flushCh)

	select {
	case <-buffer.doneCh:
		return nil
	case <-time.After(buffer.options.CloseTimeout):
		return ErrTimeout
	}
}

func (buffer *Buffer) consume() {
	items := make([]interface{}, buffer.size)
	ticker := time.NewTicker(buffer.options.AutoFlushInterval)

	count := 0
	isOpen := true
	mustFlush := false

	for isOpen {
		select {
		case item := <-buffer.dataCh:
			items[count] = item
			count++
			mustFlush = count >= len(items)
		case <-ticker.C:
			mustFlush = count > 0
		case _, open := <-buffer.flushCh:
			isOpen = open
			mustFlush = count > 0
		}

		if mustFlush {
			ticker.Stop()
			buffer.flusher(items[:count])
			count = 0
			mustFlush = false
			ticker = time.NewTicker(buffer.options.AutoFlushInterval)
		}
	}

	ticker.Stop()
	buffer.doneCh <- struct{}{}
}

func New(size uint, flusher func([]interface{}), opts ...Option) (*Buffer, error) {
// New creates a new buffer instance with the provided options.
	if size == 0 {
		panic(invalidSize)
	}
	if flusher == nil {
		panic(invalidFlusher)
	}

	options := &Options{
		AutoFlushInterval: time.Hour,
		PushTimeout:       time.Second,
		FlushTimeout:      time.Second,
		CloseTimeout:      time.Second,
	}

	for _, opt := range opts {
		opt(options)
	}

	buffer := &Buffer{
		dataCh:  make(chan interface{}),
		flushCh: make(chan struct{}),
		doneCh:  make(chan struct{}),
		size:    size,
		flusher: flusher,
		options: options,
	}
	go buffer.consume()

	return buffer, nil
}
