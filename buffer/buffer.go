package buffer

import (
	"errors"
	"go-buffer/flusher"
	"time"
)

var (
	ErrFullBuffer   = errors.New("buffer is full")
	ErrCloseTimeout = errors.New("close timed-out")
	defaultOptions  = Options{
		Size:          100,
		PushTimeout:   time.Second,
		CloseTimeout:  time.Second,
		FlushInterval: time.Second,
	}
)

type Options struct {
	Size           int64
	PushTimeout    time.Duration
	CloseTimeout   time.Duration
	FlushInterval  time.Duration
	FlusherOptions flusher.Options
}

type Buffer struct {
	bufferChannel  chan interface{}
	flusherChannel chan struct{}
	doneChannel    chan struct{}
	options        Options
	flusher        *flusher.Flusher
}

func (buffer *Buffer) Push(item interface{}) error {
	select {
	case buffer.bufferChannel <- item:
		return nil
	case <-time.After(buffer.options.PushTimeout):
		return ErrFullBuffer
	}
}

func (buffer *Buffer) ForceFlush() error {
	buffer.flusherChannel <- struct{}{}
	return nil
}

func (buffer *Buffer) Close() error {
	close(buffer.flusherChannel)

	select {
	case <-buffer.doneChannel:
		return nil
	case <-time.After(buffer.options.CloseTimeout):
		return ErrCloseTimeout
	}
}

func (buffer *Buffer) consume() {
	bufferArray := make([]interface{}, buffer.options.Size)
	ticker := time.NewTicker(buffer.options.FlushInterval)
	defer ticker.Stop()

	count := 0
	closed := false
	flush := false

	for !closed {
		select {
		case item := <-buffer.bufferChannel:
			bufferArray[count] = item
			count++
			flush = count >= len(bufferArray)
		case <-ticker.C:
			flush = count > 0
		case _, open := <-buffer.flusherChannel:
			closed = !open
			flush = count > 0
		}

		if flush {
			ticker.Stop()
			buffer.flusher.Flush(bufferArray[:count])
			count = 0
			flush = false
			ticker = time.NewTicker(buffer.options.FlushInterval)
		}
	}

	buffer.doneChannel <- struct{}{}
}

func loadDefaultOptions(options Options) Options {
	if options.Size == 0 {
		options.Size = defaultOptions.Size
	}
	if options.CloseTimeout == 0 {
		options.CloseTimeout = defaultOptions.CloseTimeout
	}
	if options.FlushInterval == 0 {
		options.FlushInterval = defaultOptions.FlushInterval
	}
	if options.PushTimeout == 0 {
		options.PushTimeout = defaultOptions.PushTimeout
	}
	return options

}

func NewBuffer(options Options) (*Buffer, error) {
	mergedOptions := loadDefaultOptions(options)
	flusherInstance, err := flusher.NewFlusher(&mergedOptions.FlusherOptions)

	if err != nil {
		return nil, err
	}

	buff := &Buffer{
		bufferChannel:  make(chan interface{}),
		flusherChannel: make(chan struct{}),
		doneChannel:    make(chan struct{}),
		options:        mergedOptions,
		flusher:        flusherInstance,
	}
	go buff.consume()
	return buff, nil
}
