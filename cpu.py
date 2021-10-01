#!/usr/bin/env python
import datetime
import pprint
import sys
import timeit

from concurrent.futures import (
    ThreadPoolExecutor,
    ProcessPoolExecutor
)
from itertools import repeat

from timing_decorators import (
    print_timing_and_results_sync
)


TENS_OF_MILLIONS_TO_COUNT_TO = [9, 5, 11, 7]


def count_to_n_million(n):
    count = 0
    for i in range((10**6) * n):
        count += 1
    return {
        'counted_to_n_million': n,
        'final_count': count,
        'processing_time': datetime.datetime.now().isoformat()
    }


@print_timing_and_results_sync
def cpu_heavy_sync():
    """
    Run cpu-heavy tasks synchronously in loop and return results sorted by processing time.
    """
    results = [count_to_n_million(n) for n in TENS_OF_MILLIONS_TO_COUNT_TO]
    return sorted(results, key=lambda x: x['processing_time'])


@print_timing_and_results_sync
def cpu_heavy_multithread():
    """
    Run cpu-heavy tasks in a pool of worker threads and return results sorted by processing time.
    """
    # ThreadPoolExecutor creates a 'pool' of threads in which to do tasks (i.e. a set you can use).
    # the max_workers argument (the number of threads in the pool) is optional, and will by default be set to something
    # sensibly scaled to the specs of the machine you are running the code on. Here we set it to one thread per
    # task in order to ensure the comparison with multiprocessing is fair.
    # https://docs.python.org/3/library/concurrent.futures.html#concurrent.futures.ThreadPoolExecutor
    with ThreadPoolExecutor(max_workers=len(TENS_OF_MILLIONS_TO_COUNT_TO)) as executor:
        # ThreadPoolExecutor.map takes a function and one or more iterables (lists, tuples, dicts, whatever) that
        # contain the arguments for your function (one iterable for each argument, if your function has multiple
        # arguments). ThreadPoolExecutor.map will then schedule the function call as a task to be executed using the
        # pool of worker threads. The size of the argument iterable(s) will determine the number of tasks scheduled.
        # https://docs.python.org/3/library/concurrent.futures.html#concurrent.futures.Executor.map
        results = executor.map(count_to_n_million, TENS_OF_MILLIONS_TO_COUNT_TO)
    return sorted(results, key=lambda x: x['processing_time'])


@print_timing_and_results_sync
def cpu_heavy_multiproc():
    """
    Run cpu-heavy tasks in a pool of worker processes and return results sorted by processing time.
    """
    # ProcessPoolExecutor creates a 'pool' of processes in which to do tasks (i.e. a set you can use).
    # the max_workers argument (the number of processes in the pool) is optional, and will by default be set to the
    # number of processors in the machine you are running the code on (). Here we set it to one process per
    # task in order to ensure the comparison with multithreading is fair.
    # https://docs.python.org/3/library/concurrent.futures.html#concurrent.futures.ProcessPoolExecutor
    with ProcessPoolExecutor(max_workers=len(TENS_OF_MILLIONS_TO_COUNT_TO)) as executor:
        # ProcessPoolExecutor.map takes a function and one or more iterables (lists, tuples, dicts, whatever) that
        # contain the arguments for your function (one iterable for each argument, if your function has multiple
        # arguments). ProcessPoolExecutor.map will then schedule the function call as a task to be executed using the
        # pool of worker processes. The size of the argument iterable(s) will determine the number of tasks scheduled.
        # https://docs.python.org/3/library/concurrent.futures.html#concurrent.futures.Executor.map
        results = executor.map(count_to_n_million, TENS_OF_MILLIONS_TO_COUNT_TO)
    return sorted(results, key=lambda x: x['processing_time'])


def main():
    cpu_heavy_sync()
    cpu_heavy_multithread()
    cpu_heavy_multiproc()


if __name__=="__main__":
    main()
