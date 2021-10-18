package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

//sets up mock structs to implement user input interfaces (allowing tests to execute without user input)

type MockUrlSelector struct {
}

func (MockUrlSelector) SelectUrl(scanner *bufio.Scanner, tasksCh chan<- string, resultsCh chan Result,
	wg *sync.WaitGroup, us UrlSelector, vs VisitsSetter, ucc ContinueChecker) int {
	return rand.Intn(10) + 1
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
	//initialises mock structs
	mus := MockUrlSelector{}
	mvs := MockVisitsSetter{}
	mcc := MockContinueChecker{}
	tasksCh := make(chan string, 10)
	resultsCh := make(chan Result)
	var wg sync.WaitGroup
	//test errors cannot be raised inside a goroutine, so instead we need to make a channel to send and receive errors
	errsCh := make(chan error, 10)
	var expected = map[string]int{
		"https://www.google.com":                         200,
		"https://www.methods.co.uk":                      200,
		"https://www.github.com":                         200,
		"https://www.stackoverflow.com":                  200,
		"https://go.dev":                                 200,
		"https://www.youtube.com":                        200,
		"https://www.ons.gov.uk":                         200,
		"https://coronavirusresources.phe.gov.uk/":       200,
		"https://campaignresources.phe.gov.uk/resources": 200,
		"https://www.twitter.com":                        200,
		"https://www.facebook.com":                       200,
	}
	workerPool(tasksCh, resultsCh, &wg, mus, mvs, mcc)
	go func() {
		//for loop ranges over results channel which will eventually be closed by userInput function
		for result := range resultsCh {
			expectedResponseCode := expected[result.url]
			if result.responseCode != expectedResponseCode {
				//send errors to errs channel
				errsCh <- fmt.Errorf("for url %v expected response code %v, got %v",
					result.url, expectedResponseCode, result.responseCode)
			}
			wg.Done()
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
