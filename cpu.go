package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	tensOfMillionsToCountTo := []int{9, 5, 11, 7}
	runCpuHeavySync(tensOfMillionsToCountTo)
	//sets number of available CPUs ==> 1
	runtime.GOMAXPROCS(1)
	runCpuHeavyGoroutines(tensOfMillionsToCountTo, 1)
	//sets number of available CPUs ==> 8
	runtime.GOMAXPROCS(8)
	runCpuHeavyGoroutines(tensOfMillionsToCountTo, 8)
}

func runCpuHeavySync(tensOfMillionsToCountTo []int) {
	t := time.Now()
	for _, v := range tensOfMillionsToCountTo {
		countToNMillion(v)
	}
	timeElapsed := time.Since(t)
	fmt.Printf("Synchronous function completed all counts in %v\n", timeElapsed)
}

func runCpuHeavyGoroutines(tensOfMillionsToCountTo []int, cpus int) {
	fmt.Printf("\nNumber of CPUs available to function: %v\n", cpus)
	c := make(chan string, len(tensOfMillionsToCountTo))
	wg := new(sync.WaitGroup)
	wg.Add(len(tensOfMillionsToCountTo))
	t := time.Now()
	for _, v := range tensOfMillionsToCountTo {
		go countToNMillionGoroutine(v, c, wg)
	}
	wg.Wait()
	timeElapsed := time.Since(t)
	close(c)
	for message := range c {
		fmt.Println(message)
	}
	fmt.Printf("Total calculation time: %v\n", timeElapsed)
}

func countToNMillion(n int) {
	count := 0
	for i := 0; i < 1000000*n; i++ {
		count += i
	}
}

func countToNMillionGoroutine(n int, c chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	t := time.Now()
	count := 0
	for i := 0; i < 1000000*n; i++ {
		count += i
	}
	timeElapsed := time.Since(t)
	c <- fmt.Sprintf("Goroutine counted to %v million in %v", n, timeElapsed)
}
