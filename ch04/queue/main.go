package main

import (
	"fmt"
	"time"
)

var (
	repeat = func(
		done <-chan any,
		values ...any,
	) <-chan any {
		valueChan := make(chan any)
		go func() {
			defer close(valueChan)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueChan <- v:
					}
				}
			}
		}()
		return valueChan
	}

	take = func(
		done <-chan any,
		valueChan <-chan any,
		num int,
	) <-chan any {
		takeChan := make(chan any)
		go func() {
			defer close(takeChan)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeChan <- <-valueChan:
				}
			}
		}()
		return takeChan
	}

	sleep = func(
		done <-chan any,
		wait time.Duration,
		c <-chan any,
	) <-chan any {
		sleepChan := make(chan any)
		go func() {
			defer close(sleepChan)
			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if !ok {
						fmt.Printf("sleep(%v) closed.\n", wait)
						return
					}
					time.Sleep(wait)
					sleepChan <- v
				}
			}
		}()
		return sleepChan
	}

	buffer = func(
		done <-chan any,
		num int,
		c <-chan any,
	) <-chan any {
		buf := make(chan any, num)
		go func() {
			defer close(buf)
			for v := range c {
				select {
				case <-done:
					return
				case buf <- v:
				}
			}
		}()
		return buf
	}
)

func main() {
	fmt.Println("--- Running noBuffer ---")
	noBuffer()

	fmt.Println("\n--- Running withBuffer ---")
	withBuffer()
}

func noBuffer() {
	done := make(chan any)
	defer close(done)

	zeros := take(done, repeat(done, 0), 3)
	short := sleep(done, 500*time.Millisecond, zeros)
	long := sleep(done, 2*time.Second, short)
	start := time.Now()
	for v := range long {
		fmt.Printf("Read %v after %v\n", v, time.Since(start))
	}
	fmt.Println("Total time: ", time.Since(start))
}

func withBuffer() {
	done := make(chan any)
	defer close(done)

	zeros := take(done, repeat(done, 0), 3)
	short := sleep(done, 500*time.Millisecond, zeros)
	buffered := buffer(done, 2, short)
	long := sleep(done, 2*time.Second, buffered)
	start := time.Now()
	for v := range long {
		fmt.Printf("Read %v after %v\n", v, time.Since(start))
	}
	fmt.Println("Total time: ", time.Since(start))
}
