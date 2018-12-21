#!/usr/bin/env bash

set -euo pipefail

declare -r DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
declare -r parent_dir="${DIR%/*}"
declare -r gparent_dir="${parent_dir%/*}"

. "$DIR/set_globals.bash"
"$DIR/docker/clean.bash"
. "$DIR/ncpus.bash"

# Config
declare -ri n="${n:-3}"
declare -r ip_start="${ip_start:-127.0.0.1}"
declare -r subnet="${subnet:-16}"
declare -r ip_range="$ip_start/$subnet"
declare -r entry="${entry:-main}" # you may use main_profile here to enable profiling
# e.g.
# n=3 entry=main_profile BUILD_DIR="$PWD" ./scripts/multi.bash

# Install deps
"$DIR/docker/install_deps.bash"

# Use -tags="netgo multi" in bgo build below to build multu lachesis version for testing
declare args="-X github.com/Fantom-foundation/go-lachesis/src/version.GitCommit=$(git rev-parse HEAD)"
if [ "$TARGET_OS" == "linux" ]; then
  args="$args -linkmode external -extldflags -static -s -w"
fi
env GOOS="$TARGET_OS" GOARCH=amd64 go build -tags="netgo multi" -ldflags "$args" -o lachesis_"$TARGET_OS" "$parent_dir/cmd/lachesis/$entry.go" || exit 1

# Create peers.json and lachesis_data_dir if needed
if [ ! -d "$DATAL_DIR/lachesis_data_dir" ]; then
    "$GOPATH/bin/batch-ethkey" -dir "$BUILD_DIR/nodes" -network "$ip_start" -inc-port -n "$n" > "$PEERS_DIR/peers.json"
    cat "$PEERS_DIR/peers.json"
    cp -rv "$BUILD_DIR/nodes" "$DATAL_DIR/lachesis_data_dir"
    cp -v "$PEERS_DIR/peers.json" "$DATAL_DIR/lachesis_data_dir/"
fi
