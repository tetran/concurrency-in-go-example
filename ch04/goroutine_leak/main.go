package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("--- Running leak1 ---")
	leak1()

	fmt.Println("\n--- Running noleak1 ---")
	noleak1()

	fmt.Println("\n--- Running leak2 ---")
	leak2()

	fmt.Println("\n--- Running noleak2 ---")
	noleak2()
}

func leak1() {
	doWork := func(strings <-chan string) <-chan any {
		completed := make(chan any)
		go func() {
			// これは実行されない
			defer fmt.Println("doWork exited.")
			defer close(completed)
			// `strings`が`nil`の場合は何も取り出せず、延々と待つことになる。 => メモリ内に残り続ける。
			for s := range strings {
				// 何かやる
				fmt.Println(s)
			}
		}()

		return completed
	}

	doWork(nil)
	// 別の何かをやる
	fmt.Println("Done.")
}

func noleak1() {
	doWork := func(
		done <-chan any,
		strings <-chan string,
	) <-chan any {
		terminated := make(chan any)
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					// 何かかやる
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()

		return terminated
	}

	done := make(chan any)
	terminated := doWork(done, nil)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	<-terminated
	fmt.Println("Done.")
}

func leak2() {
	newRandChan := func() <-chan int {
		randChan := make(chan int)
		go func() {
			// これは実行されない
			defer fmt.Println("newRandChan closure exited.")
			defer close(randChan)
			for {
				randChan <- rand.Int()
			}
		}()

		return randChan
	}

	randChan := newRandChan()
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randChan)
	}
}

func noleak2() {
	newRandChan := func(done <-chan any) <-chan int {
		randChan := make(chan int)
		go func() {
			defer fmt.Println("newRandChan closure exited.")
			defer close(randChan)
			for {
				select {
				case randChan <- rand.Int():
				case <-done:
					return
				}
			}
		}()

		return randChan
	}

	done := make(chan any)
	randChan := newRandChan(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randChan)
	}
	close(done)

	// 処理が実行中であることをシミュレート
	time.Sleep(1 * time.Second)
}
