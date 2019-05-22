#!/usr/bin/env bash
set -euo pipefail

export entry=main
export logshold=-1,-1,-1,85,0,0
export BUILD_DIR="$PWD"

#rm -rf "$BUILD_DIR/results" *.gv *.graph *.finality *.log
rm -rf *.gv *.graph *.finality *.log

for n in {3..65};
do
    export n
    ./scripts/multi.bash
    mkdir -p "$BUILD_DIR/results/$n"
    mv *.gv *.graph *.finality *.log "$BUILD_DIR/results/$n/"
done
