#!/usr/bin/env python
import pprint
import sys
import timeit

from concurrent.futures import (
    ThreadPoolExecutor,
    ProcessPoolExecutor
)
from itertools import repeat

from timing_decorators import (
    run_timed_async,
    run_timed_sync
)


TENS_OF_MILLIONS_TO_COUNT_TO = [9, 7, 10, 8, 6]


def count_to_n_million(n, results):
    count = 0
    for i in range((10**6) * n):
        count += i
    print(n, end=', ')
    results.append(n)


def print_results(results):
    print('\norder as reported in shared object:')
    print(', '.join(str(x) for x in results))


@run_timed_sync
def run_cpu_heavy_sync():
    print('Synchronous')
    print('order as printed by tasks:')
    results = []
    for n in TENS_OF_MILLIONS_TO_COUNT_TO:
        count_to_n_million(n, results)
    print_results(results)


@run_timed_sync
def run_cpu_heavy_multithread():
    print('Multithreaded')
    print('order as printed by tasks:')
    results = []
    with ThreadPoolExecutor() as executor:
        executor.map(count_to_n_million, TENS_OF_MILLIONS_TO_COUNT_TO, repeat(results))
    print_results(results)


@run_timed_sync
def run_cpu_heavy_multiproc():
    print('Multiprocessed')
    print('order as printed by tasks:')
    results = []
    with ProcessPoolExecutor() as executor:
        executor.map(count_to_n_million, TENS_OF_MILLIONS_TO_COUNT_TO, repeat(results))
    print_results(results)


def main():
    run_cpu_heavy_sync()
    run_cpu_heavy_multithread()
    run_cpu_heavy_multiproc()


if __name__=="__main__":
    main()
