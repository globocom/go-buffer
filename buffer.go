package buffer

import (
	"errors"
	"io"
	"time"
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
		closeCh chan struct{}
		doneCh  chan struct{}
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

// Close flushes the buffer and prevents it from being further used. The buffer
// cannot be used after it has been closed as all further operations will panic.
func (buffer *Buffer) Close() error {
	close(buffer.closeCh)

	var err error
	select {
	case <-buffer.doneCh:
		err = nil
	case <-time.After(buffer.options.CloseTimeout):
		err = ErrTimeout
	}

	close(buffer.dataCh)
	close(buffer.flushCh)
	return err
}

func (buffer *Buffer) consume() {
	items := make([]interface{}, buffer.options.Size)
	ticker, stopTicker := newTicker(buffer.options.FlushInterval)

	count := 0
	isOpen := true
	mustFlush := false

	for isOpen {
		select {
		case item := <-buffer.dataCh:
			items[count] = item
			count++
			mustFlush = count >= len(items)
		case <-ticker:
			mustFlush = count > 0
		case <-buffer.flushCh:
			mustFlush = count > 0
		case <-buffer.closeCh:
			isOpen = false
			mustFlush = count > 0
		}

		if mustFlush {
			stopTicker()
			buffer.options.Flusher.Write(items[:count])
			count = 0
			mustFlush = false
			ticker, stopTicker = newTicker(buffer.options.FlushInterval)
		}
	}

	stopTicker()
	close(buffer.doneCh)
}

func newTicker(interval time.Duration) (<-chan time.Time, func()) {
	if interval == 0 {
		return nil, func() {}
	}

	ticker := time.NewTicker(interval)
	return ticker.C, ticker.Stop
}

// New creates a new buffer instance with the provided options.
func New(opts ...Option) *Buffer {
	buffer := &Buffer{
		dataCh:  make(chan interface{}),
		flushCh: make(chan struct{}),
		closeCh: make(chan struct{}),
		doneCh:  make(chan struct{}),
		options: resolveOptions(opts...),
	}
	go buffer.consume()

	return buffer
}
