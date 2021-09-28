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
from timing_decorators import (
    run_timed_async,
    run_timed_sync
)


DIR_PATH = "./files_to_load"


def full_filepaths_in_dir(dirpath):
    return [os.path.join(dirpath, filename) for filename in os.listdir(DIR_PATH)]


def load_file_contents_sync(filepath):
    with open(filepath, 'r') as f:
        file_contents = f.read()
    return file_contents


async def load_file_contents_async(filepath):
    async with aiofiles.open(filepath, 'r') as f:
        file_contents = await f.read()
    return file_contents


@run_timed_sync
def run_io_heavy_sync():
    filepaths = full_filepaths_in_dir(DIR_PATH)
    file_contents = (load_file_contents_sync(filepath) for filepath in filepaths)
    return dict(zip(filepaths, file_contents))


@run_timed_async
async def run_io_heavy_async():
    filepaths = full_filepaths_in_dir(DIR_PATH)
    file_contents = await asyncio.gather(*(load_file_contents_async(filepath) for filepath in filepaths))
    return dict(zip(filepaths, file_contents))


@run_timed_sync
def run_io_heavy_multithread():
    filepaths = full_filepaths_in_dir(DIR_PATH)
    with ThreadPoolExecutor() as executor:
        file_contents = executor.map(load_file_contents_sync, filepaths)
    return dict(zip(filepaths, file_contents))


@run_timed_sync
def run_io_heavy_multiproc():
    filepaths = full_filepaths_in_dir(DIR_PATH)
    with ProcessPoolExecutor() as executor:
        file_contents = executor.map(load_file_contents_sync, filepaths)
    return dict(zip(filepaths, file_contents))


def main():
    run_io_heavy_sync()
    asyncio.run(run_io_heavy_async())
    run_io_heavy_multithread()
    run_io_heavy_multiproc()


if __name__=="__main__":
    main()

