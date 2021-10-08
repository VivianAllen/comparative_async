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

func main() {
	var urls = []string{
		"https://www.google.com",
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
	//because the urls will be sent to the tasks channel synchronously, the channel must be sufficiently buffered to
	//have space to receive them (I think)
	tasksCh := make(chan string, 8)
	//goroutines allow the workerPool function to return the results channel before workerTasks have completed, and
	//before the results channel has been closed.  This means the results channel does not need to be buffered, it can
	//print the results one by one as they are sent
	resultsCh := make(chan Result)
	//wait group set up to keep track of tasks completed so that results channel can be closed and main function can
	//return
	var wg sync.WaitGroup
	//workerPool is passed pointer to waitGroup to ensure that workerPool execution is scheduled by the same
	//waitGroup, rather than a copy
	for result := range workerPool(urls, tasksCh, resultsCh, &wg) {
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

func workerPool(urls []string, tasksCh chan string, resultsCh chan Result, wg *sync.WaitGroup) chan Result {
	//sets up worker pool
	for worker := 1; worker <= 3; worker++ {
		go workerTask(worker, tasksCh, resultsCh, wg)
	}

	for _, url := range urls {
		tasksCh <- url
		//waitGroup is incremented by 1 for each task
		wg.Add(1)
	}
	go func() {
		wg.Wait()
		userInput(tasksCh, resultsCh, urls, wg)
	}()

	return resultsCh
}

//function set up to receive from the urls channel and send to the results channel
func workerTask(id int, tasksCh <-chan string, resultsCh chan<- Result, wg *sync.WaitGroup) {
	//each time a url from the tasks channel is used inside the workerTask function, it is removed from the tasks.
	//because the tasks channel is never closed, workerTask goroutines will wait at start of for loop, and will
	//only exit when the workerPool function is returned
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
	SelectUrl(scanner *bufio.Scanner, tasksCh chan<- string, resultsCh chan Result, urls []string,
		wg *sync.WaitGroup) int
}

type StdInUrlSelector struct {
}

func (StdInUrlSelector) SelectUrl(scanner *bufio.Scanner, tasksCh chan<- string, resultsCh chan Result, urls []string,
	wg *sync.WaitGroup) int {
	fmt.Println("Enter a number between 1 and 11 and press return to run an extended speed test on one of these urls.")
	scanner.Scan()
	index, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Input error")
		userInput(tasksCh, resultsCh, urls, wg)
	} else if index < 1 || index > 11 {
		fmt.Println("Error: please enter a number between 1 and 11.")
		userInput(tasksCh, resultsCh, urls, wg)
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

type ContinueSetter interface {
	Cont(scanner *bufio.Scanner) bool
}

type StdInContinueSetter struct {
}

func (cs StdInContinueSetter) Cont(scanner *bufio.Scanner) bool {
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

func userInput(tasksCh chan<- string, resultsCh chan Result, urls []string, wg *sync.WaitGroup) {
	fmt.Println("URLS:")
	for i := 0; i < len(urls); i++ {
		fmt.Println(i+1, urls[i])
	}
	scanner := bufio.NewScanner(os.Stdin)
	urlSelector := StdInUrlSelector{}
	index := urlSelector.SelectUrl(scanner, tasksCh, resultsCh, urls, wg)
	visitsSetter := StdinVisitsSetter{}
	visits := visitsSetter.SetNumberOfVisits(scanner)
	for j := 0; j < visits; j++ {
		tasksCh <- urls[index]
		wg.Add(1)
	}
	go func() {
		wg.Wait()
		cs := StdInContinueSetter{}
		c := cs.Cont(scanner)
		if !c {
			close(resultsCh)
			return
		}
		userInput(tasksCh, resultsCh, urls, wg)
	}()
}
