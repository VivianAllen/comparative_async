package main

import (
	"fmt"

	"github.com/bradhe/stopwatch"
)

func main() {
	tensOfMillionsToCountTo := []int{9, 5, 11, 7}
	runCpuHeavySync(tensOfMillionsToCountTo)
}

func runCpuHeavySync(tensOfMillionsToCountTo []int) {
	watch := stopwatch.Start()
	for _, v := range tensOfMillionsToCountTo {
		countToNMillion(v)
	}
	watch.Stop()
	fmt.Println("Synchronous function: time taken", watch.Milliseconds())
}

func countToNMillion(n int) {
	count := 0
	for i := 0; i < 1000000*n; i++ {
		count += i
	}
}
