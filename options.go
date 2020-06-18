package buffer

import "time"

const (
	invalidInterval = "invalid interval"
	invalidTimeout  = "invalid timeout"
)

type (
	Options struct {
		AutoFlush         bool
		AutoFlushInterval time.Duration
		PushTimeout       time.Duration
		FlushTimeout      time.Duration
		CloseTimeout      time.Duration
	}

	Option func(*Options)
)

func WithAutoFlush(interval time.Duration) Option {
	if interval <= 0 {
		panic(invalidInterval)
	}

	return func(options *Options) {
		options.AutoFlush = true
		options.AutoFlushInterval = interval
	}
}

func WithPushTimeout(timeout time.Duration) Option {
	if timeout <= 0 {
		panic(invalidTimeout)
	}

	return func(options *Options) {
		options.PushTimeout = timeout
	}
}

func WithFlushTimeout(timeout time.Duration) Option {
	if timeout <= 0 {
		panic(invalidTimeout)
	}

	return func(options *Options) {
		options.FlushTimeout = timeout
	}
}

func WithCloseTimeout(timeout time.Duration) Option {
	if timeout <= 0 {
		panic(invalidTimeout)
	}

	return func(options *Options) {
		options.CloseTimeout = timeout
	}
}
