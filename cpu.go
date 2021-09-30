package main

import (
	"fmt"

	"github.com/bradhe/stopwatch"
)

func main() {
	tensOfMillionsToCountTo := []int{9, 5, 11, 7}
	m := make(map[int]string)
	for _, v := range tensOfMillionsToCountTo {
		countToNMillion(v, m)
	}
	fmt.Println(m)
}

func countToNMillion(n int, m map[int]string) {
	count := 0
	watch := stopwatch.Start()
	for i := 0; i < 1000000*n; i++ {
		count += i
	}
	watch.Stop()
	m[n] = watch.Milliseconds().String()
}
