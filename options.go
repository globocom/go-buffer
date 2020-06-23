package buffer

import (
	"fmt"
	"time"
)

const (
	invalidInterval = "%s: interval must be greater than zero"
	invalidTimeout  = "%s: timeout cannot be negative"
)

type (
	// Configuration options.
	Options struct {
		FlushInterval time.Duration
		PushTimeout   time.Duration
		FlushTimeout  time.Duration
		CloseTimeout  time.Duration
	}

	// Option setter.
	Option func(*Options)
)

// WithFlushInterval sets the interval between automatic flushes.
func WithFlushInterval(interval time.Duration) Option {
	return func(options *Options) {
		options.FlushInterval = interval
	}
}

// WithPushTimeout sets how long a push should wait before giving up.
func WithPushTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.PushTimeout = timeout
	}
}

// WithFlushTimeout sets how long a manual flush should wait before giving up.
func WithFlushTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.FlushTimeout = timeout
	}
}

// WithCloseTimeout sets how long
func WithCloseTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.CloseTimeout = timeout
	}
}

func validateOptions(options *Options) error {
	if options.FlushInterval < 0 {
		return fmt.Errorf(invalidInterval, "FlushInterval")
	}

	if options.PushTimeout < 0 {
		return fmt.Errorf(invalidTimeout, "PushTimeout")
	}

	if options.FlushInterval < 0 {
		return fmt.Errorf(invalidTimeout, "FlushTimeout")
	}

	if options.CloseTimeout < 0 {
		return fmt.Errorf(invalidTimeout, "CloseTimeout")
	}

	return nil
}
