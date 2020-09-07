package buffer_test

import (
	"testing"

	"github.com/globocom/go-buffer"
)

func BenchmarkBuffer(b *testing.B) {
	noop := buffer.FlusherFunc(func([]interface{}) {})

	b.Run("push only", func(b *testing.B) {
		sut := buffer.New(
			buffer.WithSize(uint(b.N)+1),
			buffer.WithFlusher(noop),
		)
		defer sut.Close()

		for i := 0; i < b.N; i++ {
			err := sut.Push(i)
			if err != nil {
				b.Fail()
			}
		}
	})

	b.Run("push and flush", func(b *testing.B) {
		sut := buffer.New(
			buffer.WithSize(1),
			buffer.WithFlusher(noop),
		)
		defer sut.Close()

		for i := 0; i < b.N; i++ {
			err := sut.Push(i)
			if err != nil {
				b.Fail()
			}
		}
	})
}
