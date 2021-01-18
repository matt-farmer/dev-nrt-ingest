#!/bin/bash

set -e

# delete all binary files
find . -type f -executable -exec sh -c "file -i '{}' | grep -q 'x-executable; charset=binary'" \; -print | xargs rm -f
for f in $(find ./ -name '*.log' -or -name '*.doc'); do rm $f; done

rm -rf ./file/*.json
rm -rf ./db/*
rm -f rrd.xml sif.xml