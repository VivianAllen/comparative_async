package main

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/bradhe/stopwatch"
)

func main() {
	tensOfMillionsToCountTo := []int{9, 5, 11, 7}
	runCpuHeavySync(tensOfMillionsToCountTo)
	runtime.GOMAXPROCS(1)
	runCpuHeavyGoroutines(tensOfMillionsToCountTo, 1)
	runtime.GOMAXPROCS(8)
	runCpuHeavyGoroutines(tensOfMillionsToCountTo, 8)
}

func runCpuHeavySync(tensOfMillionsToCountTo []int) {
	watch := stopwatch.Start()
	for _, v := range tensOfMillionsToCountTo {
		countToNMillion(v)
	}
	watch.Stop()
	fmt.Println("Synchronous function: time taken", watch.Milliseconds())
}

func runCpuHeavyGoroutines(tensOfMillionsToCountTo []int, numOfCPUs int) {
	wg := new(sync.WaitGroup)
	wg.Add(len(tensOfMillionsToCountTo))
	watch := stopwatch.Start()
	for _, v := range tensOfMillionsToCountTo {
		go countToNMillionGoroutine(v, wg)
	}
	wg.Wait()
	watch.Stop()
	fmt.Printf("Goroutines function: num of CPUs: %v, time taken: %v\n", numOfCPUs, watch.Milliseconds())
}

func countToNMillion(n int) {
	count := 0
	for i := 0; i < 1000000*n; i++ {
		count += i
	}
}

func countToNMillionGoroutine(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	count := 0
	for i := 0; i < 1000000*n; i++ {
		count += i
	}
}
