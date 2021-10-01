import functools
import pprint
import timeit


def print_timing_and_results_sync(function):
    """
    Decorator to debug log execution time of decorated synchronous function. Invoked by placing @run_timed_sync
    above the target function, e.g:

    @run_timed_sync
    def my_function():
        do_stuff()
    """
    @functools.wraps(function)
    def wrapper(*args, **kwargs):
        start_execution_time_s = timeit.default_timer()
        results = function(*args, **kwargs)
        execution_time_s = timeit.default_timer() - start_execution_time_s
        print(f'Function {function.__name__} executed in {execution_time_s}s.')
        print('Results:')
        pp = pprint.PrettyPrinter(indent=2)
        pp.pprint(results)
        print('')
        return results
    return wrapper


def print_timing_and_results_async(function):
    """
    Decorator to debug log execution time of decorated asynchronous function. Invoked by placing @run_timed_async
    above the target coroutine / async function, e.g:

    @run_timed_async
    async def my_function():
        await do_stuff()
    """
    @functools.wraps(function)
    async def wrapper(*args, **kwargs):
        start_execution_time_s = timeit.default_timer()
        results = await function(*args, **kwargs)
        execution_time_s = timeit.default_timer() - start_execution_time_s
        print(f'Function {function.__name__} executed in {execution_time_s}s.')
        print('Results:')
        pp = pprint.PrettyPrinter(indent=2)
        pp.pprint(results)
        print('')
        return results
    return wrapper
