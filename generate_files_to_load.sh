#!/bin/bash

# Generate set of 20 random text files from 80mb to 100MB in size in 'files_to_load' directory created in directory
# this script is run from. Files will be named for their size in MB.

mkdir -p ./files_to_load
size_in_mb=100

for n in {0..20}
do
  this_file_size_in_mb=$(($size_in_mb - $n))
  size_in_b=$(($this_file_size_in_mb *  1000000))
  base64 /dev/urandom | head -c "$size_in_b" > files_to_load/"$this_file_size_in_mb"
done
