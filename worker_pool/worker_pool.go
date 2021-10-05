package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	var urls = []string{"https://www.google.com", "https://www.methods.co.uk", "https://www.github.com", "https://www.stackoverflow.com", "https://go.dev", "https://www.youtube.com", "https://www.ons.gov.uk", "https://coronavirusresources.phe.gov.uk/", "https://campaignresources.phe.gov.uk/resources", "https://www.twitter.com", "https://www.facebook.com"}

	workerPool(urls)
}

type Result struct {
	workerId     int
	url          string
	responseCode int
	speed        float64
}

func workerPool(urls []string) {
	tasks := make(chan string, 10)
	results := make(chan Result)
	//sets up wait group to track task completion
	wg := new(sync.WaitGroup)

	//sets up worker pool
	for w := 1; w <= 3; w++ {
		go workerTask(w, tasks, results)
	}

	for _, url := range urls {
		tasks <- url
		//for each url added to the tasks channel, wait group counter is incremented by one
		wg.Add(1)
	}

	go func(results <-chan Result) {
		for result := range results {
			fmt.Println(result)
			//for each result printed, wait group counter is decremented
			wg.Done()
		}
	}(results)
	//function will wait for wait group counter to reach zero before exiting
	wg.Wait()
}

//function set up to receive from the urls channel and send to the results channel
func workerTask(id int, tasks <-chan string, results chan<- Result) {
	//each time a url from the tasks channel is used inside the workerTask function, it is removed from the tasks
	//when there are no tasks remaining in the channel, workerTask goroutine will exit
	for url := range tasks {
		start := time.Now()
		resp, err := http.Get(url)
		timeElapsed := time.Since(start).Seconds()
		if err != nil {
			fmt.Println(err)
		}
		//Result is sent to results channel
		results <- Result{workerId: id, url: url, responseCode: resp.StatusCode, speed: timeElapsed}
	}
}
