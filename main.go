package main

import (
	"fmt"
	"go-buffer/buffer"
	"go-buffer/flusher"
	"time"
)

func onStart(item []interface{}) error {
	fmt.Println("onStart")
	return nil
}
func OnEach(item interface{}) error {
	fmt.Println("onEach", item)
	time.Sleep(1 * time.Microsecond)
	return nil
}
func OnEnd() error {
	fmt.Print("onEnd")
	return nil
}

func onError(item []interface{}, err error) {
	fmt.Println(err)
}
func onErrorS(item interface{}, err error) {
	fmt.Println(err)
}

func main() {
	bufferOptions := buffer.Options{
		Size: 5,
		FlusherOptions: flusher.Options{
			OnStart:      onStart,
			OnEach:       OnEach,
			OnEnd:        OnEnd,
			OnStartError: onError,
			OnEndError:   onError,
			OnEachError:  onErrorS,
		},
	}
	buf, err := buffer.NewBuffer(bufferOptions)
	if err != nil {
		panic(err)
	}
	buf.Push(1)
	buf.Push(2)
	buf.Push(3)
	buf.Push(4)
	buf.Push(5)
	buf.Push(6)
	buf.Push(7)
	buf.Push(8)
	buf.Push(5)
	buf.Push(6)
	buf.Push(7)
	buf.Push(8)
	//buf.Push(4)
	//buf.Push(5)
	time.Sleep(5 * time.Second)

	select {}
}
