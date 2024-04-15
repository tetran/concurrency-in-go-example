package main

import (
	"bytes"
	"fmt"
	"sync"
)

func main() {
	fmt.Println("--- Running adhoc ---")
	adhoc()

	fmt.Println("--- Running lexical ---")
	lexical()
}

func adhoc() {
	// 規約により、dataへはloopData関数内のみからアクセスするようにしている。
	data := make([]int, 4)
	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}

func lexical() {
	// `printData`は`data`変数宣言の前にあるので、`data`変数に直接アクセスできない。
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("golang")
	// 各goroutineが`data`の一部にしかアクセスできないように「拘束」している。
	// => レキシカルスコープによって間違ったアクセスを不可能にしている。
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])

	wg.Wait()
}
