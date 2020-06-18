package buffer

import (
	"errors"
	"io"
	"time"
)

var (
	ErrFull             = errors.New("buffer is full")
	ErrOperationTimeout = errors.New("operation timed-out")
)

type (
	// Buffer represents a data buffer that is asynchronously flushed, either manually or
	// automatically.
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

// Push appends an item to the end of the Buffer. It might return an error if the buffer is full.
func (buffer *Buffer) Push(item interface{}) error {
	select {
	case buffer.dataCh <- item:
		return nil
	case <-time.After(buffer.options.PushTimeout):
		return ErrFull
	}
}

// Flush outputs the buffer to a permanent destination. It might return an error if the buffer is already being flushed.
func (buffer *Buffer) Flush() error {
	select {
	case buffer.flushCh <- struct{}{}:
		return nil
	case <-time.After(buffer.options.FlushTimeout):
		return ErrOperationTimeout
	}
}

// Close flushes the buffer and prevents it from being further used.
func (buffer *Buffer) Close() error {
	close(buffer.flushCh)

	select {
	case <-buffer.doneCh:
		return nil
	case <-time.After(buffer.options.CloseTimeout):
		return ErrOperationTimeout
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

// New creates a new buffer instance
func New(size uint, flusher func([]interface{}), opts ...Option) (*Buffer, error) {
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
