package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	tasks := make(chan string, 10)
	results := make(chan Result, 11)

	//sets up worker pool
	for w := 1; w <= 3; w++ {
		go workerTask(w, tasks, results)
	}

	for _, url := range urls {
		tasks <- url
	}
	//closing the channel communicates to worker pool that no more values will be sent to urlsCh - see https://tour.golang.org/concurrency/4
	close(tasks)

	for result := range results {
		fmt.Println(result, len(results))
	}

}

var urls = []string{"https://www.google.com", "https://www.methods.co.uk", "https://www.github.com", "https://www.stackoverflow.com", "https://go.dev", "https://www.youtube.com", "https://www.ons.gov.uk", "https://coronavirusresources.phe.gov.uk/", "https://campaignresources.phe.gov.uk/resources", "https://www.twitter.com", "https://www.facebook.com"}

type Result struct {
	workerId     int
	url          string
	responseCode int
	speed        float64
}

//function set up to receive from the urls channel and send to the results channel
func workerTask(id int, tasks <-chan string, results chan<- Result) {
	//each time a url from the tasks channel is used inside the workerTask function, it is removed from the tasks
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
