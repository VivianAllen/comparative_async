#!/usr/bin/python

import os
import pprint
import sys
import timeit

DIR_PATH = "./files_to_load"

def load_file_contents(filepath):
    with open(filepath, 'r') as f:
        file_contents = f.read()
    return file_contents


def dir_files_to_dict(dir_path):
    filenames = os.listdir(dir_path)
    file_contents = (load_file_contents(os.path.join(dir_path, filename)) for filename in filenames)
    return dict(zip(filenames, file_contents))


def main():
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
    file_contents = dir_files_to_dict(dir_path)
    t1 = timeit.default_timer()

    if log_out:
        pp = pprint.PrettyPrinter(indent=4)
        pp.pprint(file_contents)

    print(f"All files processed in {(t1 - t0) * 1000} milliseconds")


if __name__=="__main__":
    main()

