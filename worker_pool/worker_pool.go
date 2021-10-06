package main

import (
	"fmt"
	"net/http"
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
		"https://www.twitter.com", "https://www.facebook.com"}
	//because the urls will be sent to the tasks channel synchronously, the channel must be sufficiently buffered to
	//have space to receive them (I think)
	tasksCh := make(chan string, 8)
	//goroutines allow the workerPool function to return the results channel before workerTasks have completed, and
	//before the results channel has been closed.  This means the results channel does not need to be buffered, it can
	//print the results one by one as they are sent
	resultsCh := make(chan Result)
	for result := range workerPool(urls, tasksCh, resultsCh) {
		fmt.Println(result)
	}
}

type Result struct {
	workerId     int
	url          string
	responseCode int
	speed        float64
}

func workerPool(urls []string, tasksCh chan string, resultsCh chan Result) chan Result {
	//wait group set up to keep track of tasks completed so that results channel can be closed and main function can
	//return
	var wg sync.WaitGroup

	//sets up worker pool
	for worker := 1; worker <= 3; worker++ {
		//goroutine is passed pointer to waitGroup to ensure that workerPool execution is scheduled by the same
		//waitGroup, rather than a copy
		go workerTask(worker, tasksCh, resultsCh, &wg)
	}

	for _, url := range urls {
		tasksCh <- url
	}

	//closing the results channel allows main function to return once all results have been printed (is this bad
	//practice - not closing channel from sender function?)
	go func() {
		wg.Wait()
		close(resultsCh)
	}()
	return resultsCh
}

//function set up to receive from the urls channel and send to the results channel
func workerTask(id int, tasksCh <-chan string, resultsCh chan<- Result, wg *sync.WaitGroup) {
	//each time a url from the tasks channel is used inside the workerTask function, it is removed from the tasks
	//because the tasks channel is never closed, workerTask goroutines will wait at start of for loop, and will only
	//exit when the workerPool function is returned
	for url := range tasksCh {
		//waitGroup is incremented by 1 for each task
		wg.Add(1)
		start := time.Now()
		resp, err := http.Get(url)
		timeElapsed := time.Since(start).Seconds()
		if err != nil {
			fmt.Println(err)
		}
		//Result is sent to results channel
		resultsCh <- Result{workerId: id, url: url, responseCode: resp.StatusCode, speed: timeElapsed}
		//waitGroup decremented by 1 as each task is completed
		wg.Done()
	}
}
