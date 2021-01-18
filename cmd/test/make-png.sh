#!/usr/bin/env bash

# rm empty
for f in heap_*.out; do
    if [ ! -s $f ]; then
	rm -f $f
    fi
done

# print first and last profile
go tool pprof -png $(ls -1 heap_*.out | head -n1)
go tool pprof -png $(ls -1 heap_*.out | tail -n1)

