package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var URLS = [11]string{"https://www.google.com",
	"https://www.methods.co.uk",
	"https://www.github.com",
	"https://www.stackoverflow.com",
	"https://go.dev",
	"https://www.youtube.com",
	"https://www.ons.gov.uk",
	"https://coronavirusresources.phe.gov.uk/",
	"https://campaignresources.phe.gov.uk/resources",
	"https://www.twitter.com",
	"https://www.facebook.com"}

func main() {
	//because the urls will be sent to the tasks channel synchronously, the channel must be sufficiently buffered to
	//have space to receive them (I think).  As there are three workers three of them will be sent immediately, so the
	//channel needs space for the remaining 8
	tasksCh := make(chan string, 8)
	//goroutines allow the workerPool function to return the results channel before workerTasks have completed, and
	//before the results channel has been closed.  This means the results channel does not need to be buffered, it can
	//print the results one by one as they are sent
	resultsCh := make(chan Result)
	//wait group is a counter that we can use to keep track of asynchronous tasks as they are completed so that when
	//finished, results channel can be closed and main function can return
	var wg sync.WaitGroup
	userUrlSelector := StdInUrlSelector{}
	userVisitsSetter := StdinVisitsSetter{}
	userContinueChecker := StdInContinueChecker{}
	//workerPool is passed pointer to waitGroup to ensure that workerPool execution is scheduled by the same
	//waitGroup, rather than a copy
	for result := range workerPool(tasksCh, resultsCh, &wg, userUrlSelector, userVisitsSetter, userContinueChecker) {
		fmt.Println(result)
		//waitGroup decremented by 1 as each task is completed
		wg.Done()
	}
}

type Result struct {
	workerId     int
	url          string
	responseCode int
	speed        float64
}

//function takes in interfaces for UrlSelector, VisitsSetter and ContinueChecker to allow for mocking
func workerPool(tasksCh chan string, resultsCh chan Result, wg *sync.WaitGroup, us UrlSelector, vs VisitsSetter, ucc ContinueChecker) chan Result {
	//sets up worker pool
	for worker := 1; worker <= 3; worker++ {
		go workerTask(worker, tasksCh, resultsCh, wg)
	}

	//sends urls to tasks channel
	for _, url := range URLS {
		tasksCh <- url
		//waitGroup is incremented by 1 for each task
		wg.Add(1)
	}

	//waits for wg.Done() to execute before calling userInput
	go func() {
		wg.Wait()
		userInput(tasksCh, resultsCh, wg, us, vs, ucc)
	}()

	return resultsCh
}

//function set up to receive from the urls channel and send to the results channel
func workerTask(id int, tasksCh <-chan string, resultsCh chan<- Result, wg *sync.WaitGroup) {
	//each time a url from the tasks channel is used inside the workerTask function, it is removed from the tasks.
	//because the tasks channel is never closed, workerTask goroutines will wait at start of for loop, and will only
	//exit when the results channel is closed in the userInput function, closing the for loop in the main routine
	for url := range tasksCh {
		start := time.Now()
		resp, err := http.Get(url)
		timeElapsed := time.Since(start).Seconds()
		if err != nil {
			fmt.Println(err)
		}
		//Result is sent to results channel
		resultsCh <- Result{workerId: id, url: url, responseCode: resp.StatusCode, speed: timeElapsed}
	}
}

type UrlSelector interface {
	SelectUrl(scanner *bufio.Scanner, tasksCh chan<- string, resultsCh chan Result, wg *sync.WaitGroup,
		us UrlSelector, vs VisitsSetter, ucc ContinueChecker) int
}

type StdInUrlSelector struct {
}

func (StdInUrlSelector) SelectUrl(scanner *bufio.Scanner, tasksCh chan<- string, resultsCh chan Result,
	wg *sync.WaitGroup, us UrlSelector, vs VisitsSetter, ucc ContinueChecker) int {
	fmt.Println("Enter a number between 1 and 11 and press return to run an extended speed test on one of these urls.")
	scanner.Scan()
	index, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Input error")
		userInput(tasksCh, resultsCh, wg, us, vs, ucc)
	} else if index < 1 || index > 11 {
		fmt.Println("Error: please enter a number between 1 and 11.")
		userInput(tasksCh, resultsCh, wg, us, vs, ucc)
	}
	return index - 1
}

type VisitsSetter interface {
	SetNumberOfVisits(scanner *bufio.Scanner) int
}

type StdinVisitsSetter struct {
}

func (StdinVisitsSetter) SetNumberOfVisits(scanner *bufio.Scanner) int {
	fmt.Println("How many times would you like to check this url?  Enter a number between 1 and 30.")
	scanner.Scan()
	visits, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Input error: programme will default to 10 checks.")
		visits = 10
	} else if visits < 1 {
		visits = 1
	} else if visits > 30 {
		visits = 30
	}
	return visits
}

type ContinueChecker interface {
	Cont(scanner *bufio.Scanner) bool
}

type StdInContinueChecker struct {
}

func (cs StdInContinueChecker) Cont(scanner *bufio.Scanner) bool {
	fmt.Println("Continue? y/n")
	scanner.Scan()
	input := scanner.Text()
	switch input {
	case "y", "Y":
		return true
	case "n", "N":
		return false
	default:
		return cs.Cont(scanner)
	}
}

func userInput(tasksCh chan<- string, resultsCh chan Result, wg *sync.WaitGroup, us UrlSelector, vs VisitsSetter, ucc ContinueChecker) {
	fmt.Println("URLS:")
	for i := 0; i < len(URLS); i++ {
		fmt.Println(i+1, URLS[i])
	}
	scanner := bufio.NewScanner(os.Stdin)
	index := us.SelectUrl(scanner, tasksCh, resultsCh, wg, us, vs, ucc)
	visits := vs.SetNumberOfVisits(scanner)
	for j := 0; j < visits; j++ {
		tasksCh <- URLS[index]
		wg.Add(1)
	}
	go func() {
		wg.Wait()
		c := ucc.Cont(scanner)
		if !c {
			close(resultsCh)
			return
		}
		userInput(tasksCh, resultsCh, wg, us, vs, ucc)
	}()
}
