package main

import (
	"fmt"
	"testing"
)

func TestWorkerPool(t *testing.T) {
	tasks := make(chan string, 10)
	results := make(chan Result)
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
	workerPool(urls, tasks, results)
	go func() {
		for result := range results {
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
	}()
	for err := range errsCh {
		if err != nil {
			t.Fatal(err)
		}
	}
}
