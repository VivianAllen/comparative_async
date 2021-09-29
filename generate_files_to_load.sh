#!/bin/bash

mkdir -p ./files_to_load
size_in_mb=${1-100}

for n in {0..20}
do
  this_file_size_in_mb=$(($size_in_mb - $n))
  size_in_b=$(($this_file_size_in_mb *  1000000))
  base64 /dev/urandom | head -c "$size_in_b" > files_to_load/"$this_file_size_in_mb"
done
