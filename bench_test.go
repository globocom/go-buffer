package buffer_test

import (
	"testing"

	"github.com/globocom/go-buffer/v3"
)

func BenchmarkBuffer(b *testing.B) {
	noop := func([]int) {}

	b.Run("push only", func(b *testing.B) {
		sut := buffer.New(
			noop,
			buffer.WithSize(uint(b.N)+1),
		)
		defer sut.Close()

		for b.Loop() {
			err := sut.Push(1)
			if err != nil {
				b.Fail()
			}
		}
	})

	b.Run("push and flush", func(b *testing.B) {
		sut := buffer.New(
			noop,
			buffer.WithSize(1),
		)
		defer sut.Close()

		for b.Loop() {
			err := sut.Push(1)
			if err != nil {
				b.Fail()
			}
		}
	})
}
