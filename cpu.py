#!/usr/bin/env python
import pprint
import sys
import timeit

from concurrent.futures import (
    ThreadPoolExecutor,
    ProcessPoolExecutor
)
from timing_decorators import (
    run_timed_async,
    run_timed_sync
)


CONCURRENCY = 4


def count_to_10_million(x):
    count = 0
    for i in range(10**7):
        count += i


@run_timed_sync
def run_cpu_heavy_sync():
    for i in range(CONCURRENCY):
        count_to_10_million(i)
    return 'done'


@run_timed_sync
def run_cpu_heavy_multithread():
    with ThreadPoolExecutor() as executor:
        results = executor.map(count_to_10_million, range(CONCURRENCY))


@run_timed_sync
def run_cpu_heavy_multiproc():
    with ProcessPoolExecutor() as executor:
        results = executor.map(count_to_10_million, range(CONCURRENCY))


def main():
    run_cpu_heavy_sync()
    run_cpu_heavy_multithread()
    run_cpu_heavy_multiproc()


if __name__=="__main__":
    main()
