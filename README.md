# Comparative Async / Concurrent Code

NB Assumes unix OS w. bash/Zsh (MacOS / Zsh used for development)

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

- run benchmark tests for cpu-heavy code example with `go test -bench=.`