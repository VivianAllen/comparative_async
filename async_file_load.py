import asyncio
import aiofiles
import os
import pprint

DIR_PATH = "./files_to_load"

async def load_file_contents(filepath):
    async with aiofiles.open(filepath, 'r') as f:
        file_contents = await f.read()
    return file_contents

async def dir_files_to_dict(dir_path):
    filenames = os.listdir(dir_path)
    tasks = [asyncio.ensure_future(load_file_contents(os.path.join(dir_path, filename))) for filename in filenames]
    file_contents = await asyncio.gather(*tasks)
    pp = pprint.PrettyPrinter(indent=4)
    pp.pprint(dict(zip(filenames, file_contents)))

if __name__=="__main__":
    loop = asyncio.get_event_loop()
    loop.run_until_complete(dir_files_to_dict(DIR_PATH))

