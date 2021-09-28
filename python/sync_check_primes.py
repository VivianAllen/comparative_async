#!/usr/bin/env python
import pprint
import sys
import timeit


def is_prime(number):
    if number == 1:
        return f"{number} is not prime"
    for i in range(2, int(number/2)+1):
        if (number % i) == 0:
            return f"{number} is not prime"
    return f"{number} is prime"


def check_primes_up_to(primes_up_to):
    return [is_prime(number) for number in range(1, primes_up_to+1)]


def main():
    args = sys.argv[1:]

    if len(args) > 0:
        primes_up_to = int(args[0])
    else:
        primes_up_to = 10

    if len(args) > 1:
         log_out = True if args[1].lower() == 'true' else False
    else:
        log_out = False


    t0 = timeit.default_timer()
    file_contents = check_primes_up_to(primes_up_to)
    t1 = timeit.default_timer()

    if log_out:
        pp = pprint.PrettyPrinter(indent=4)
        pp.pprint(file_contents)
    print(f"All numbers processed in {(t1 - t0) * 1000} milliseconds")


if __name__=="__main__":
    main()

