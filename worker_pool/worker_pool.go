package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	tasks := make(chan string, 10)
	results := make(chan Result)

	//sets up worker pool
	for w := 1; w <= 3; w++ {
		go workerTask(w, tasks, results)
	}

	for _, url := range urls {
		tasks <- url
	}
	//closing the channel communicates to worker pool that no more values will be sent to urlsCh - see https://tour.golang.org/concurrency/4
	close(tasks)

	for i := 1; i <= len(urls); i++ {
		result := <-results
		fmt.Println(result)
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
func workerTask(id int, urls <-chan string, results chan<- Result) {
	for url := range urls {
		start := time.Now()
		resp, err := http.Get(url)
		timeElapsed := time.Since(start).Seconds()
		if err != nil {
			fmt.Println(err)
		}
		results <- Result{workerId: id, url: url, responseCode: resp.StatusCode, speed: timeElapsed}
	}
}
