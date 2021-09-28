#!/usr/bin/env python
import pprint
import sys
import timeit


def cpu_heavy(x):
    count = 0
    for i in range(10**8):
        count += i


def run_cpu_heavy(n):
    for i in range(n):
        cpu_heavy(n)


def main():
    args = sys.argv[1:]

    if len(args) > 0:
        n = int(args[0])
    else:
        n = 4

    t0 = timeit.default_timer()
    file_contents = run_cpu_heavy(n)
    t1 = timeit.default_timer()

    print(f"Finished in {(t1 - t0) * 1000} milliseconds")


if __name__=="__main__":
    main()

