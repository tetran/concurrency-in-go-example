package main

import "fmt"

func main() {
	repeat := func(
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
	take := func(
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
	tee := func(
		done <-chan any,
		in <-chan any,
	) (_, _ <-chan any) {
		out1 := make(chan any)
		out2 := make(chan any)
		go func() {
			defer close(out1)
			defer close(out2)
			for val := range orDone(done, in) {
				// nil代入が外に影響しないようにローカルコピー
				var out1, out2 = out1, out2
				// out1, out2の両方に確実に書き込まれるように2回繰り返す
				for i := 0; i < 2; i++ {
					select {
					case out1 <- val:
						// ただ一度のみ書き込まれるようにnilを代入。以降はもう一方に書き込まれるようにする。
						out1 = nil
					case out2 <- val:
						out2 = nil
					}
				}
			}
		}()
		return out1, out2
	}

	done := make(chan any)
	defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 1, 2), 4))

	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
}
