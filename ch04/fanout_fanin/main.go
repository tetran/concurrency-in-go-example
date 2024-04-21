package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
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

	repeatFn = func(
		done <-chan any,
		fn func() any,
	) <-chan any {
		valueChan := make(chan any)
		go func() {
			defer close(valueChan)
			for {
				select {
				case <-done:
					return
				case valueChan <- fn():
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

	toInt = func(
		done <-chan any,
		valueChan <-chan any,
	) <-chan int {
		intChan := make(chan int)
		go func() {
			defer close(intChan)
			for v := range valueChan {
				select {
				case <-done:
					return
				case intChan <- v.(int):
				}
			}
		}()
		return intChan
	}

	primeFinder = func(
		done <-chan any,
		valueChan <-chan int,
	) <-chan any {
		primeChan := make(chan any)
		go func() {
			defer close(primeChan)
			for v := range valueChan {
				select {
				case <-done:
					return
				default:
					isPrime := true
					for i := 2; i < v; i++ {
						if v%i == 0 {
							isPrime = false
							break
						}
					}
					if isPrime {
						primeChan <- v
					}
				}
			}
		}()
		return primeChan
	}

	fanIn = func(
		done <-chan any,
		channels ...<-chan any,
	) <-chan any {
		var wg sync.WaitGroup
		mpChan := make(chan any)
		multiplex := func(c <-chan any) {
			defer wg.Done()
			for i := range c {
				select {
				case <-done:
					return
				case mpChan <- i:
				}
			}
		}

		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}

		go func() {
			wg.Wait()
			close(mpChan)
		}()

		return mpChan
	}
)

func main() {
	fmt.Println("--- Running no_fofi ---")
	no_fofi()

	fmt.Println("\n--- Running fofi ---")
	fofi()
}

func no_fofi() {
	rand := func() any { return rand.Intn(50000000) }

	done := make(chan any)
	defer close(done)

	start := time.Now()

	randIntChan := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, randIntChan), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v\n", time.Since(start))
}

func fofi() {
	rand := func() any { return rand.Intn(50000000) }

	done := make(chan any)
	defer close(done)

	start := time.Now()

	randIntChan := toInt(done, repeatFn(done, rand))

	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	finders := make([]<-chan any, numFinders)
	fmt.Println("Primes:")
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntChan)
	}

	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v\n", time.Since(start))
}
