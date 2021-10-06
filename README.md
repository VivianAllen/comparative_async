# Comparative Async / Concurrent Code

NB Assumes unix OS w. bash/Zsh (MacOS / Zsh used for development)


## Intro

The idea here is to implement code for doing multiple tasks concurrently in a non-blocking way (aka asynchronous / concurrent code) in various languages. Ways of doing this:
- event-loop-based asynchronous code (i.e. JS Promises or Async/Await, Python Futures or Async/Await)
- multithreading
- multiprocessing
- things that combine aspects of the above in a way thats hard to categorise (looking at you, [go routines](https://tour.golang.org/concurrency/1))

### Ideas for things to implement
- run a single task in whatever flavours of non-blocking code exist in your target language, and return the results.
- run a set of concurrent tasks in whatever flavours of non-blocking code exist in your target language, wait for _all_ to finish, collect and return the results.
- create a long-running pool of workers that can be dynamically given tasks to perform by the main routine, and will communicate the results back to the main routine when done, and then will remain alive to wait for new tasks.

### Comparing performance for CPU and I/O bound operations

Most forms of async / concurrent code are better for applications involving a lot of I/O - if you are waiting on an external system (disk read/write, http call, etc) to return what you asked for, it makes sense that unblocking your code to allow more work to be done while it waits (or to wait for multiple things at once) is more efficient.

If your application involves lots of calculations to be done and processed in-memory, the advantage is less clear. If your code is all making use of the same set of resources (CPU / memory etc) then there may be no advantage in trying to do it all at once. Multiprocessing, in which each task is run as a separate instance that gets allocated its own resources, is in theory a better choice for concurrent CPU-heavy tasks (although the overhead of starting up new processes may remove that advantage for smaller tasks).

For Things That Don't Go In Pigeonholes ([Go routines](https://tour.golang.org/concurrency/1), [Nodes single-threaded event loop / multithreaded worker pool setup](https://sararavi14.medium.com/node-js-threading-model-badf7bb5fffa)), the above distinction is a bit more murky.

Why not give your implementations of async / concurrent task runners some I/O-heavy and CPU-heavy tasks to run, and collect data on their performance (execution time, memory footprint etc if possible) for comparison?

## Python

NB version 3.9.7 used for development.

### Setup
- setup virtual environment in project root with `python3 -m venv env` (or your preferred method)
- activate virtual environment with `source env/bin/activate` (or your preferred methods)
- install dependencies with `pip install -r requirements.txt`

### Use
- setup target files for io-heavy example with `./generate_files_to_load.sh`
- run io-heavy code example with `./io.py`
- run cpu-heavy code example with `./cpu.py`

### Cleanup
- target files can be remove with `./rm_files_to_load.sh`

## Go

Version 1.17 used for development.

### Use

- run cpu-heavy code example with `go run cpu.go`
- run benchmark tests for cpu-heavy code example with `go test -bench=.`
- run concurrent worker pool example with `go run worker_pool/worker_pool.go`
- run tests for worker pool example with `go test ./worker_pool/`

