#!/bin/bash

mkdir -p ./files_to_load
size_in_mb=${1-1}
size_in_b=$(( $size_in_mb *  1000000))

for n in {0..20}
do
  base64 /dev/urandom | head -c "$size_in_b" > files_to_load/file_"$n"
done
