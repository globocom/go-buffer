package buffer

import "time"

type (
	Options struct {
		AutoFlushInterval time.Duration
		PushTimeout       time.Duration
		FlushTimeout      time.Duration
		CloseTimeout      time.Duration
	}
)
