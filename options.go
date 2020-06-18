package buffer

import "time"

type (
	Options struct {
		AutoFlushInterval time.Duration
		PushTimeout       time.Duration
		FlushTimeout      time.Duration
		CloseTimeout      time.Duration
	}

	Option func(*Options)
)

func WithAutoFlush(interval time.Duration) Option {
	return func(options *Options) {
		options.AutoFlushInterval = interval
	}
}

func WithPushTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.PushTimeout = timeout
	}
}

func WithFlushTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.FlushTimeout = timeout
	}
}

func WithCloseTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.CloseTimeout = timeout
	}
}
