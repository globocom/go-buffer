package buffer_test

import (
	"testing"

	"github.com/globocom/go-buffer/v3"
)

func BenchmarkBuffer(b *testing.B) {
	noop := buffer.FlusherFunc[int](func([]int) {})

	b.Run("push only", func(b *testing.B) {
		sut := buffer.New[int](
			noop,
			buffer.WithSize(uint(b.N)+1),
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
		sut := buffer.New[int](
			noop,
			buffer.WithSize(1),
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
