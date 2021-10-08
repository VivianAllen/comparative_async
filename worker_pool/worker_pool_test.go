package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

type MockUrlSelector struct {
}

type MockVisitsSetter struct {
}

func (MockVisitsSetter) SetNumberOfVisits(scanner *bufio.Scanner) int {
	return rand.Intn(29) + 1
}

type MockContinueChecker struct {
}

func (MockContinueChecker) Cont(scanner *bufio.Scanner) bool {
	return false
}

func TestWorkerPool(t *testing.T) {
	//sets up mock functions
	mus := MockUrlSelector{}
	mvs := MockVisitsSetter{}
	mcc := MockContinueChecker{}
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
	fmt.Println("test")
	workerPool(urls, tasksCh, resultsCh, &wg, mus, mvs, mcc)
	fmt.Println("test2")
	go func() {
		// for i := 0; i < 4; i++ {
		// 	result := <-resultsCh
		// 	expectedResponseCode := expected[result.url]
		// 	if result.responseCode != expectedResponseCode {
		// 		//send errors to errs channel
		// 		errsCh <- fmt.Errorf("for url %v expected response code %v, got %v",
		// 			result.url, expectedResponseCode, result.responseCode)
		// 	}
		// }
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
	}()
	for err := range errsCh {
		if err != nil {
			t.Fatal(err)
		}
	}
}
