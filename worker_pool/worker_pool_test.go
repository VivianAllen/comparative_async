package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestWorkerPool(t *testing.T) {
	//sets up mock functions
	vs := VisitsSetter{}
	vs

	tasksCh := make(chan string, 10)
	resultsCh := make(chan Result)
	var wg sync.WaitGroup
	//because test errors cannot be raised inside a goroutine, we need to make an channel to send and receive errors
	errsCh := make(chan error, 10)
	var urls = []string{
		"https://www.google.com",
		"https://www.methods.co.uk",
		"https://www.github.com",
		"https://www.stackoverflow.com"}
	var expected = map[string]int{
		"https://www.google.com":        200,
		"https://www.methods.co.uk":     200,
		"https://www.github.com":        200,
		"https://www.stackoverflow.com": 200,
	}
	workerPool(urls, tasksCh, resultsCh, &wg)
	go func() {
		for result := range resultsCh {
			expectedResponseCode := expected[result.url]
			if result.responseCode != expectedResponseCode {
				//send errors to errs channel
				errsCh <- fmt.Errorf("for url %v expected response code %v, got %v",
					result.url, expectedResponseCode, result.responseCode)
			}

		}
		//close errs channel when results channel is closed, allowing for loop below to complete and test function to
		//exit
		close(errsCh)
		fmt.Println("test4")
	}()
	for err := range errsCh {
		if err != nil {
			t.Fatal(err)
		}
	}
}
