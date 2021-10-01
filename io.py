#!/usr/bin/env python
import asyncio
import aiofiles
import datetime
import os
import pprint
import sys

from concurrent.futures import (
    ThreadPoolExecutor,
    ProcessPoolExecutor
)
from itertools import repeat

from timing_decorators import (
    print_timing_and_results_async,
    print_timing_and_results_sync
)


DIR_PATH = "./files_to_load"


def full_filepaths_in_dir(dirpath):
    return [os.path.join(dirpath, filename) for filename in os.listdir(DIR_PATH)]


def get_result_string(filepath, file_contents_length):
    file_name = os.path.basename(filepath)
    processing_time = datetime.datetime.now().isoformat()
    return f"File {file_name} with {file_contents_length} characters processed at {processing_time}"


def get_file_contents_length_sync(filepath, results):
    """
    Load contents from filepath, print 'finished' message with filename, and append 'finished' message to shared
    results object.
    """
    with open(filepath, 'r') as f:
        file_contents = f.read()
    results.append(get_result_string(filepath, len(file_contents)))


async def get_file_contents_length_async(filepath, results):
    """
    Load contents from filepath asynchronously, print 'finished' message with filename, and append 'finished' message to
    shared results object.
    """
    # aiofiles provides an awaitable context manager to open files and await their contents (i.e. to allow the event
    # loop to pass on to execute other code while the file contents are retrieved from the disk).
    # https://pypi.org/project/aiofiles/
    # Note use of 'async' keyword in context manager construction. Things defined with 'async' are known as 'coroutines'
    # in python: https://docs.python.org/3/library/asyncio-task.html#id1
    async with aiofiles.open(filepath, 'r') as f:
        file_contents = await f.read()
    results.append(get_result_string(filepath, len(file_contents)))


@print_timing_and_results_sync
def io_heavy_sync():
    """
    Run io-heavy tasks with shared results object synchronously in loop and display results.
    """
    filepaths = full_filepaths_in_dir(DIR_PATH)
    results = []
    for filepath in filepaths:
        get_file_contents_length_sync(filepath, results)
    return results


@print_timing_and_results_async
async def io_heavy_async():
    """
    Run io-heavy tasks with shared results object asynchronously on event loop and display results.
    Note use of 'async' keyword in definition. Things defined with 'async' are known as 'coroutines' in python:
    https://docs.python.org/3/library/asyncio-task.html#id1
    """
    filepaths = full_filepaths_in_dir(DIR_PATH)
    results = []

    # asyncio.gather takes one or more function calls and schedules them as asynchronous tasks / coroutines to be
    # performed on an event loop. NB, gather expects to be given the tasks as
    # asyncio.gather(function_call_1, function_call_2, ... etc), so here we use * to unpack / exhaust an iterable
    # full of function calls. NB note use of await here to make sure we actually wait for all tasks to finish before
    # proceeding!
    # https://docs.python.org/3/library/asyncio-task.html#running-tasks-concurrently
    await asyncio.gather(*(get_file_contents_length_async(filepath, results) for filepath in filepaths))

    return results


@print_timing_and_results_sync
def io_heavy_multithread():
    """
    Run io-heavy tasks in a pool of worker threads with a shared results object and display results.
    """
    filepaths = full_filepaths_in_dir(DIR_PATH)
    results = []

    # ThreadPoolExecutor creates a 'pool' of threads in which to do tasks (i.e. a set you can use).
    # the max_workers argument (the number of threads in the pool) is optional, and will by default be set to something
    # sensibly scaled to the specs of the machine you are running the code on. Here we set it to one thread per
    # task in order to ensure the comparison with multiprocessing is fair.
    # https://docs.python.org/3/library/concurrent.futures.html#concurrent.futures.ThreadPoolExecutor
    with ThreadPoolExecutor(max_workers=4) as executor:

        # ThreadPoolExecutor.map takes a function and one or more iterables (lists, tuples, dicts, whatever) that
        # contain the arguments for your function (one iterable for each argument, if your function has multiple
        # arguments). ThreadPoolExecutor.map will then schedule the function call as a task to be executed using the
        # pool of worker threads. The size of the argument iterable(s) will determine the number of tasks scheduled.
        # NB - the iterables must therefore all be the same length. Here we use itertools.repeat to pass the same
        # argument in as many times as needed (repeat creates a generator object that always returns the same thing when
        # 'next' is called on it).
        # https://docs.python.org/3/library/concurrent.futures.html#concurrent.futures.Executor.map
        executor.map(get_file_contents_length_sync, filepaths, repeat(results))

    return results


@print_timing_and_results_sync
def io_heavy_multiproc_no_shared_memory():
    """
    Run io-heavy tasks in a pool of worker processes with a shared results object and display results.
    """
    filepaths = full_filepaths_in_dir(DIR_PATH)
    results = []

    # ProcessPoolExecutor creates a 'pool' of processes in which to do tasks (i.e. a set you can use).
    # the max_workers argument (the number of processes in the pool) is optional, and will by default be set to the
    # number of processors in the machine you are running the code on (). Here we set it to one process per
    # task in order to ensure the comparison with multithreading is fair.
    # https://docs.python.org/3/library/concurrent.futures.html#concurrent.futures.ProcessPoolExecutor
    with ProcessPoolExecutor(max_workers=4) as executor:

        # ProcessPoolExecutor.map takes a function and one or more iterables (lists, tuples, dicts, whatever) that
        # contain the arguments for your function (one iterable for each argument, if your function has multiple
        # arguments). ProcessPoolExecutor.map will then schedule the function call as a task to be executed using the
        # pool of worker processes. The size of the argument iterable(s) will determine the number of tasks scheduled.
        # NB - the iterables must therefore all be the same length. Here we use itertools.repeat to pass the same
        # argument in as many times as needed (repeat creates a generator object that always returns the same thing when
        # 'next' is called on it).
        # https://docs.python.org/3/library/concurrent.futures.html#concurrent.futures.Executor.map
        executor.map(get_file_contents_length_sync, filepaths, repeat(results))

    # NB - the results stored in the shared object given to each worker process will be blank! This is because by
    # default separate processes do not share memory with each other! You need to do something fancy to get that to
    # work: https://docs.python.org/3/library/multiprocessing.shared_memory.html
    return results


def main():
    io_heavy_sync()
    # here we use asyncio.run to start off an async task / coroutine and run it until it has finished. This allows you
    # to call async code as part of a sync function without relying on callbacks or some other signalling method to
    # execute your code in the expected order.
    # https://docs.python.org/3/library/asyncio-task.html#asyncio.run
    asyncio.run(io_heavy_async())
    io_heavy_multithread()
    io_heavy_multiproc_no_shared_memory()


if __name__=="__main__":
    main()

