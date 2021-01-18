#!/usr/bin/env bash
cd $(dirname $0)

set -ex

rm -fr dbs
mkdir dbs
go run ../lachesis -datadir=dbs &
sleep 10

for (( ; ; )); do
    curl -sK -v http://localhost:8080/debug/pprof/heap > heap_$(date '+%Y%m%d-%H%M').out
    sleep 120
done
