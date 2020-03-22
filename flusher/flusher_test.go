package flusher_test

import (
	"errors"
	"go-buffer/flusher"
	"testing"
)

func flusherOptionsFixture() flusher.Options {
	return flusher.Options{
		OnStart: func(item []interface{}) error {
			return nil
		},
		OnEach: func(item interface{}) error {
			return nil
		},
		OnEnd: func(item []interface{}) error {
			return nil
		},
		OnStartError: func(item []interface{}, err error) {
		},
		OnEachError: func(item interface{}, err error) {

		},
		OnEndError: func(item []interface{}, err error) {

		},
	}
}

func TestFlusher(t *testing.T) {

	t.Run("OnStart function should be called just one time by Flush", func(t *testing.T) {
		// Arrange
		called := false
		flusherOptions := flusherOptionsFixture()
		flusherOptions.OnStart = func(item []interface{}) error {
			called = true
			return nil
		}
		flusherInstance, _ := flusher.NewFlusher(&flusherOptions)
		items := make([]interface{}, 3)
		items[0] = 1
		items[1] = 2
		items[2] = 3

		// Act
		flusherInstance.Flush(items)

		// Assert
		if !called {
			t.Errorf("OnStart should be called")
		}
	})

	t.Run("OnEnd function should be called just one time by Flush", func(t *testing.T) {
		// Arrange
		called := false
		flusherOptions := flusherOptionsFixture()
		flusherOptions.OnEnd = func(item []interface{}) error {
			called = true
			return nil
		}
		flusherInstance, _ := flusher.NewFlusher(&flusherOptions)
		items := make([]interface{}, 3)
		items[0] = 1
		items[1] = 2
		items[2] = 3

		// Act
		flusherInstance.Flush(items)

		// Assert
		if !called {
			t.Errorf("OnEnd should be called")
		}
	})

	t.Run("OnEach function should be called once per item on Flush", func(t *testing.T) {
		// Arrange
		calls := 0
		flusherOptions := flusherOptionsFixture()
		flusherOptions.OnEach = func(item interface{}) error {
			calls = calls + 1
			return nil
		}
		flusherInstance, _ := flusher.NewFlusher(&flusherOptions)
		items := make([]interface{}, 3)
		items[0] = 1
		items[1] = 2
		items[2] = 3

		// Act
		flusherInstance.Flush(items)

		// Assert
		if calls != 3 {
			t.Errorf("OnEach should be called 3 times")
		}
	})

	t.Run("When OnStart return an error, OnStartError should be called with that error", func(t *testing.T) {
		// Arrange
		var expectedErr = errors.New("Ops, try again")
		var returnedErr error
		flusherOptions := flusherOptionsFixture()
		flusherOptions.OnStart = func(item []interface{}) error {
			return expectedErr
		}
		flusherOptions.OnStartError = func(item []interface{}, err error) {
			returnedErr = err
		}

		flusherInstance, _ := flusher.NewFlusher(&flusherOptions)
		items := make([]interface{}, 3)
		items[0] = 1

		// Act
		flusherInstance.Flush(items)

		// Assert
		if returnedErr != expectedErr {
			t.Errorf("OnStart should return the expected error")
		}
	})

	t.Run("When OnEnd return an error, OnEndError should be called with that error", func(t *testing.T) {
		// Arrange
		var expectedErr = errors.New("Ops, try again")
		var returnedErr error
		flusherOptions := flusherOptionsFixture()
		flusherOptions.OnEnd = func(item []interface{}) error {
			return expectedErr
		}
		flusherOptions.OnEndError = func(item []interface{}, err error) {
			returnedErr = err
		}

		flusherInstance, _ := flusher.NewFlusher(&flusherOptions)
		items := make([]interface{}, 3)
		items[0] = 1

		// Act
		flusherInstance.Flush(items)

		// Assert
		if returnedErr != expectedErr {
			t.Errorf("OnEnd should return the expected error")
		}
	})

	t.Run("When OnEach return an error, OnEachError should be called with that error", func(t *testing.T) {
		// Arrange
		var expectedErr = errors.New("Ops, try again")
		var returnedErr error
		flusherOptions := flusherOptionsFixture()
		flusherOptions.OnEach = func(item interface{}) error {
			return expectedErr
		}
		flusherOptions.OnEachError = func(item interface{}, err error) {
			returnedErr = err
		}

		flusherInstance, _ := flusher.NewFlusher(&flusherOptions)
		items := make([]interface{}, 3)
		items[0] = 1

		// Act
		flusherInstance.Flush(items)

		// Assert
		if returnedErr != expectedErr {
			t.Errorf("OnEach should return the expected error")
		}
	})
}
