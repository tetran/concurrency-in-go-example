package main

import (
	"fmt"
	"time"
)

func main() {
	repeat := func(
		values ...any,
	) <-chan any {
		valueChan := make(chan any)
		go func() {
			defer close(valueChan)
			for {
				for _, v := range values {
					valueChan <- v
				}
			}
		}()
		return valueChan
	}
	orDone := func(done, c <-chan any) <-chan any {
		valChan := make(chan any)
		go func() {
			defer close(valChan)
			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case valChan <- v:
					case <-done:
					}
				}
			}
		}()
		return valChan
	}

	done := make(chan any)
	go func() {
		time.Sleep(3 * time.Millisecond)
		close(done)
	}()

	for val := range orDone(done, repeat(1)) {
		fmt.Print(val)
	}
}
