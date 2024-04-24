package main

import "fmt"

func main() {
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
	bridge := func(
		done <-chan any,
		chanChan <-chan <-chan any,
	) <-chan any {
		valChan := make(chan any)
		go func() {
			defer close(valChan)
			for {
				var c <-chan any
				select {
				case maybeChan, ok := <-chanChan:
					if !ok {
						return
					}
					c = maybeChan
				case <-done:
					return
				}
				for val := range orDone(done, c) {
					select {
					case valChan <- val:
					case <-done:
					}
				}
			}
		}()
		return valChan
	}

	genVals := func() <-chan <-chan any {
		chanChan := make(chan (<-chan any))
		go func() {
			defer close(chanChan)
			for i := 0; i < 10; i++ {
				c := make(chan any, 1)
				c <- i
				close(c)
				chanChan <- c
			}
		}()
		return chanChan
	}

	for v := range bridge(nil, genVals()) {
		fmt.Printf("%v ", v)
	}
}
