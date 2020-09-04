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

// Close flushes the buffer and prevents it from being further used. If it succeeds,
// the buffer cannot be used after it has been closed as all further operations will panic.
func (buffer *Buffer) Close() error {
	select {
	case buffer.closeCh <- struct{}{}:
		// noop
	case <-time.After(buffer.options.CloseTimeout):
		return ErrTimeout
	}

	select {
	case <-buffer.doneCh:
		close(buffer.dataCh)
		close(buffer.flushCh)
		close(buffer.closeCh)
		return nil
	case <-time.After(buffer.options.CloseTimeout):
		return ErrTimeout
	}
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
	options := &Options{
		Size:          0,
		Flusher:       nil,
		FlushInterval: 0,
		PushTimeout:   time.Second,
		FlushTimeout:  time.Second,
		CloseTimeout:  time.Second,
	}

	for _, opt := range opts {
		opt(options)
	}

	if err := validateOptions(options); err != nil {
		panic(err.Error())
	}

	buffer := &Buffer{
		dataCh:  make(chan interface{}),
		flushCh: make(chan struct{}),
		closeCh: make(chan struct{}),
		doneCh:  make(chan struct{}),
		options: options,
	}
	go buffer.consume()

	return buffer
}
