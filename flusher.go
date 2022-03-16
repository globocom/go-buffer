package buffer

type (
	// Flusher represents a destination of buffered data.
	Flusher[T any] interface {
		Write(items []T)
	}

	// FlusherFunc represents a flush function.
	FlusherFunc[T any] func(items []T)
)

func (fn FlusherFunc[T]) Write(items []T) {
	fn(items)
}
