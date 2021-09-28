import functools
import timeit


def run_timed_sync(function):
    """
    Decorator to debug log execution time of decorated synchronous function.
    """
    @functools.wraps(function)
    def wrapper(*args, **kwargs):
        start_execution_time_s = timeit.default_timer()
        result = function(*args, **kwargs)
        execution_time_s = timeit.default_timer() - start_execution_time_s
        print(f'Function {function.__name__} executed in {execution_time_s}s.')
        return result
    return wrapper


def run_timed_async(function):
    """
    Decorator to debug log execution time of decorated asynchronous function.
    """
    @functools.wraps(function)
    async def wrapper(*args, **kwargs):
        start_execution_time_s = timeit.default_timer()
        result = await function(*args, **kwargs)
        execution_time_s = timeit.default_timer() - start_execution_time_s
        print(f'Function {function.__name__} executed in {execution_time_s}s.')
        return result
    return wrapper
