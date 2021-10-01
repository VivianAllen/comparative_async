package main

import (
	"runtime"
	"sync"
)

func main() {
	tensOfMillionsToCountTo := []int{9, 5, 11, 7}
	runCpuHeavySync(tensOfMillionsToCountTo)
	//sets number of available CPUs ==> 1
	runtime.GOMAXPROCS(1)
	runCpuHeavyGoroutines(tensOfMillionsToCountTo)
	//sets number of available CPUs ==> 8
	runtime.GOMAXPROCS(8)
	runCpuHeavyGoroutines(tensOfMillionsToCountTo)
}

func runCpuHeavySync(tensOfMillionsToCountTo []int) {
	for _, v := range tensOfMillionsToCountTo {
		countToNMillion(v)
	}
}

func runCpuHeavyGoroutines(tensOfMillionsToCountTo []int) {
	wg := new(sync.WaitGroup)
	wg.Add(len(tensOfMillionsToCountTo))
	for _, v := range tensOfMillionsToCountTo {
		go countToNMillionGoroutine(v, wg)
	}
	wg.Wait()
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
