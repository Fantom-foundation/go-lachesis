#!/usr/bin/env bash

go tool pprof -png $(ls -1 heap_*.out | head -n1)
go tool pprof -png $(ls -1 heap_*.out | tail -n1)

