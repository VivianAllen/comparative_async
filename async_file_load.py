#!/usr/bin/env python

import asyncio
import aiofiles
import os
import pprint
import sys
import timeit

DIR_PATH = "./files_to_load"

async def load_file_contents(filepath):
    async with aiofiles.open(filepath, 'r') as f:
        file_contents = await f.read()
    return file_contents


async def dir_files_to_dict(dir_path):
    filenames = os.listdir(dir_path)
    file_contents = await asyncio.gather(
        *(load_file_contents(os.path.join(dir_path, filename)) for filename in filenames)
    )
    return dict(zip(filenames, file_contents))


async def main():
    args = sys.argv[1:]

    if len(args) > 0:
        log_out = True if args[0].lower() == 'true' else False
    else:
        log_out = False

    if len(args) > 1:
        dir_path = args[1]
    else:
        dir_path = DIR_PATH

    t0 = timeit.default_timer()
    file_contents = await dir_files_to_dict(dir_path)
    t1 = timeit.default_timer()

    if log_out:
        pp = pprint.PrettyPrinter(indent=4)
        pp.pprint(file_contents)

    print(f"All files processed in {(t1 - t0) * 1000} milliseconds")


if __name__=="__main__":
    asyncio.run(main())

