package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("--- Running err_noreturn ---")
	err_noreturn()

	fmt.Println("--- Running err_return ---")
	err_return()

	fmt.Println("--- Running err_return2 ---")
	err_return2()
}

func err_noreturn() {
	checkStatus := func(
		done <-chan any,
		urls ...string,
	) <-chan *http.Response {
		responses := make(chan *http.Response)
		go func() {
			defer close(responses)
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					// エラーを返すことができない。
					fmt.Println(err)
					continue
				}
				select {
				case <-done:
					return
				case responses <- resp:
				}
			}
		}()

		return responses
	}

	done := make(chan any)
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}
	for responses := range checkStatus(done, urls...) {
		fmt.Printf("Response: %v\n", responses.Status)
	}
}

type ErrorHandlingResult struct {
	Error    error
	Response *http.Response
}

func err_return() {
	checkStatus := func(done <-chan any, urls ...string) <-chan ErrorHandlingResult {
		results := make(chan ErrorHandlingResult)
		go func() {
			defer close(results)

			for _, url := range urls {
				var result ErrorHandlingResult
				resp, err := http.Get(url)
				result = ErrorHandlingResult{Error: err, Response: resp}
				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()
		return results
	}

	done := make(chan any)
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}

func err_return2() {
	checkStatus := func(done <-chan any, urls ...string) <-chan ErrorHandlingResult {
		results := make(chan ErrorHandlingResult)
		go func() {
			defer close(results)

			for _, url := range urls {
				var result ErrorHandlingResult
				resp, err := http.Get(url)
				result = ErrorHandlingResult{Error: err, Response: resp}
				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()
		return results
	}

	done := make(chan any)
	defer close(done)

	errCount := 0
	urls := []string{"a", "https://www.google.com", "b", "c", "d"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}
