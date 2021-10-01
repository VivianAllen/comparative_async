package main

import (
	"runtime"
	"testing"
)
var tensOfMillionsToCountTo=[]int{9, 5, 11, 7}
func BenchmarkRunCpuHeavySync(b *testing.B) {
	for n:=0;n<b.N;n++{
		runCpuHeavySync(tensOfMillionsToCountTo)
	}
}

func BenchmarkRunCpuHeavyGoroutinesOneCPU(b *testing.B) {
	runtime.GOMAXPROCS(1)
	for n:=0;n<b.N;n++{
		runCpuHeavyGoroutines(tensOfMillionsToCountTo)
	}
	
}

func BenchmarkRunCpuHeavyGoroutinesMultipleCPUS(b *testing.B) {
	runtime.GOMAXPROCS(8)
	for n:=0;n<b.N;n++{
		runCpuHeavyGoroutines(tensOfMillionsToCountTo)
	}
}
