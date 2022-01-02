package buffer

import (
	"errors"
	"io"
	"time"
)

var (
	// ErrTimeout indicates an operation has timed out.
	ErrTimeout = errors.New("operation timed-out")
	// ErrClosed indicates the buffer is closed and can no longer be used.
	ErrClosed = errors.New("buffer is closed")
)

type (
	// Buffer represents a data buffer that is asynchronously flushed, either manually or automatically.
	Buffer[T any] struct {
		io.Closer
		flushFunc func([]T)
		dataCh    chan T
		flushCh   chan struct{}
		closeCh   chan struct{}
		doneCh    chan struct{}
		options   *Options
	}
)

// New creates a new buffer instance with the provided flush function and options.
// It panics if provided with a nil flush function.
func New[T any](flushFunc func([]T), opts ...Option) *Buffer[T] {
	if flushFunc == nil {
		panic("flush function cannot be nil")
	}

	buffer := &Buffer[T]{
		flushFunc: flushFunc,
		dataCh:    make(chan T),
		flushCh:   make(chan struct{}),
		closeCh:   make(chan struct{}),
		doneCh:    make(chan struct{}),
		options:   resolveOptions(opts...),
	}
	go buffer.consume()

	return buffer
}

// Push appends an item to the end of the buffer.
//
// It returns an ErrTimeout if it cannot be performed in a timely fashion, and
// an ErrClosed if the buffer has been closed.
func (buffer *Buffer[T]) Push(item T) error {
	if buffer.closed() {
		return ErrClosed
	}

	select {
	case buffer.dataCh <- item:
		return nil
	case <-time.After(buffer.options.PushTimeout):
		return ErrTimeout
	}
}

// Flush outputs the buffer to a permanent destination.
//
// It returns an ErrTimeout if if cannot be performed in a timely fashion, and
// an ErrClosed if the buffer has been closed.
func (buffer *Buffer[T]) Flush() error {
	if buffer.closed() {
		return ErrClosed
	}

	select {
	case buffer.flushCh <- struct{}{}:
		return nil
	case <-time.After(buffer.options.FlushTimeout):
		return ErrTimeout
	}
}

// Close flushes the buffer and prevents it from being further used.
//
// It returns an ErrTimeout if if cannot be performed in a timely fashion, and
// an ErrClosed if the buffer has already been closed.
//
// An ErrTimeout can either mean that a flush could not be triggered, or it can
// mean that a flush was triggered but it has not finished yet. In any case it is
// safe to call Close again.
func (buffer *Buffer[T]) Close() error {
	if buffer.closed() {
		return ErrClosed
	}

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

func (buffer *Buffer[T]) closed() bool {
	select {
	case <-buffer.doneCh:
		return true
	default:
		return false
	}
}

func (buffer *Buffer[T]) consume() {
	count := 0
	items := make([]T, buffer.options.Size)
	mustFlush := false
	ticker, stopTicker := newTicker(buffer.options.FlushInterval)

	isOpen := true
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
			buffer.flushFunc(items[:count])

			count = 0
			items = make([]T, buffer.options.Size)
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
