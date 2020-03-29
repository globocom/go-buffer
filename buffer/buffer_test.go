package buffer_test

import (
	"go-buffer/buffer"
	"testing"
	"time"
)

func bufferOptionsFixture() buffer.Options {
	return buffer.Options{
		Size:          5,
		PushTimeout:   time.Second,
		CloseTimeout:  time.Second,
		FlushInterval: time.Second,
	}
}

func TestBuffer(t *testing.T) {

	t.Run("When buffer is full push should return an error", func(t *testing.T) {
		// Arrange
		bufferOptions := bufferOptionsFixture()
		bufferOptions.Flusher = func(items []interface{}) {
			time.Sleep(time.Second)
		}
		bufferOptions.Size = 2

		bufferInstance, _ := buffer.NewBuffer(bufferOptions)

		// Act
		push1 := bufferInstance.Push(1)
		push2 := bufferInstance.Push(2)
		push3 := bufferInstance.Push(3)

		// Assert
		if push1 != nil {
			t.Errorf("Buffer should not is full")
		}
		if push2 != nil {
			t.Errorf("Buffer should not is full")
		}
		if push3 != buffer.ErrFullBuffer {
			t.Errorf("Buffer should is full")
		}
	})

	t.Run("Should flush buffer when wait time is greater than FlushInterval", func(t *testing.T) {
		// Arrange
		flushed := false
		bufferOptions := bufferOptionsFixture()
		bufferOptions.FlushInterval = time.Millisecond * 500
		bufferOptions.Flusher = func(items []interface{}) {
			flushed = true
		}

		bufferInstance, _ := buffer.NewBuffer(bufferOptions)

		// Act
		bufferInstance.Push(1)

		time.Sleep(time.Second)

		// Assert
		if !flushed {
			t.Errorf("Buffer should have been flushed")
		}
	})

	t.Run("Should not flush buffer when wait time is less than FlushInterval", func(t *testing.T) {
		// Arrange
		flushed := false
		bufferOptions := bufferOptionsFixture()
		bufferOptions.Flusher = func(items []interface{}) {
			flushed = true
		}

		bufferInstance, _ := buffer.NewBuffer(bufferOptions)

		// Act
		bufferInstance.Push(1)

		// Assert
		if flushed {
			t.Errorf("Buffer should have been flushed")
		}
	})

	t.Run("Should flush buffer when size is full", func(t *testing.T) {
		// Arrange
		flushed := false
		bufferOptions := bufferOptionsFixture()
		bufferOptions.Flusher = func(items []interface{}) {
			flushed = true
		}
		bufferOptions.Size = 1

		bufferInstance, _ := buffer.NewBuffer(bufferOptions)

		// Act
		bufferInstance.Push(1)
		bufferInstance.Push(2)

		time.Sleep(time.Millisecond)

		// Assert
		if !flushed {
			t.Errorf("Buffer should have been flushed")
		}
	})

}
