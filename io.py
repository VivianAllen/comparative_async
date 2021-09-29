#!/usr/bin/env python
import asyncio
import aiofiles
import os
import pprint
import sys

from concurrent.futures import (
    ThreadPoolExecutor,
    ProcessPoolExecutor
)
from itertools import repeat

from timing_decorators import (
    run_timed_async,
    run_timed_sync
)


DIR_PATH = "./files_to_load"


def full_filepaths_in_dir(dirpath):
    return [os.path.join(dirpath, filename) for filename in os.listdir(DIR_PATH)]


def load_file_contents_sync(filepath, results):
    with open(filepath, 'r') as f:
        file_contents = f.read()
    result_str = os.path.basename(filepath)
    print(result_str, end=', ')
    results.append(result_str)

async def load_file_contents_async(filepath, results):
    async with aiofiles.open(filepath, 'r') as f:
        file_contents = await f.read()
    result_str = os.path.basename(filepath)
    print(result_str, end=', ')
    results.append(result_str)


def print_results(results):
    print('\norder as reported in shared object:')
    print(', '.join(str(x) for x in results))

@run_timed_sync
def run_io_heavy_sync():
    print('Synchronous')
    print('order as printed by tasks:')
    filepaths = full_filepaths_in_dir(DIR_PATH)
    results = []
    for filepath in filepaths:
        load_file_contents_sync(filepath, results)
    print_results(results)


@run_timed_async
async def run_io_heavy_async():
    print('Asynchronous')
    print('order as printed by tasks:')
    filepaths = full_filepaths_in_dir(DIR_PATH)
    results = []
    await asyncio.gather(*(load_file_contents_async(filepath, results) for filepath in filepaths))
    print_results(results)


@run_timed_sync
def run_io_heavy_multithread():
    print('Multithreaded')
    print('order as printed by tasks:')
    filepaths = full_filepaths_in_dir(DIR_PATH)
    results = []
    with ThreadPoolExecutor() as executor:
        executor.map(load_file_contents_sync, filepaths, repeat(results))
    print_results(results)


@run_timed_sync
def run_io_heavy_multiproc():
    print('Multiprocessed')
    print('order as printed by tasks:')
    filepaths = full_filepaths_in_dir(DIR_PATH)
    results = []
    with ProcessPoolExecutor() as executor:
        executor.map(load_file_contents_sync, filepaths, repeat(results))
    print_results(results)


def main():
    run_io_heavy_sync()
    asyncio.run(run_io_heavy_async())
    run_io_heavy_multithread()
    run_io_heavy_multiproc()


if __name__=="__main__":
    main()

