package buffer_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/globocom/go-buffer"
)

var _ = Describe("Options", func() {
	It("sets up auto flush", func() {
		// arrange
		opts := &buffer.Options{}

		// act
		buffer.WithAutoFlush(5 * time.Second)(opts)

		// assert
		Expect(opts.AutoFlush).To(BeTrue())
		Expect(opts.AutoFlushInterval).To(Equal(5 * time.Second))
	})

	It("sets up push timeout", func() {
		// arrange
		opts := &buffer.Options{}

		// act
		buffer.WithPushTimeout(10 * time.Second)(opts)

		// assert
		Expect(opts.PushTimeout).To(Equal(10 * time.Second))
	})

	It("sets up flush timeout", func() {
		// arrange
		opts := &buffer.Options{}

		// act
		buffer.WithFlushTimeout(15 * time.Second)(opts)

		// assert
		Expect(opts.FlushTimeout).To(Equal(15 * time.Second))
	})

	It("sets up close timeout", func() {
		// arrange
		opts := &buffer.Options{}

		// act
		buffer.WithCloseTimeout(3 * time.Second)(opts)

		// assert
		Expect(opts.CloseTimeout).To(Equal(3 * time.Second))
	})
})
