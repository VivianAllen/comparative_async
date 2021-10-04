package main

import (
	"fmt"
	"testing"
)

func TestWorkerPool(t *testing.T) {
	tasks := make(chan string, 10)
	results := make(chan Result)
	//because test errors cannot be raised inside a goroutine, we need to make an channel to send and receive errors
	errs := make(chan error, 10)
	var urls = []string{"https://www.google.com"}
	expected := make(map[string]int)
	expected["https://www.google.com"] = 404
	workerPool(urls, tasks, results)
	go func() {
		for result := range results {
			expectedResponseCode := expected[result.url]
			if result.responseCode != expectedResponseCode {
				//send errors to errs channel
				errs <- fmt.Errorf("for url %v expected response code %v, got %v", result.url, expectedResponseCode, result.responseCode)
			}

		}
		//close errs channel when results channel is closed, allowing for loop below to complete and test function to exit
		close(errs)
	}()
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
}
