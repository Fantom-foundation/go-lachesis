#!/usr/bin/env bash

set -euo pipefail
OPTIND=1         # Reset in case getopts has been used previously in the shell.

declare -r DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

. "$DIR/set_globals.bash"
. "$DIR/ncpus.bash"

# Config
declare -ri n="${n:-3}"
declare -r ip_start="${ip_start:-127.0.0.0}"
declare -r subnet="${subnet:-16}"
declare -r ip_range="$ip_start/$subnet"
declare -r node_addr="127.0.0.1"

# Create loopback aliases and cp json.peers per node datadir
declare -i node_num=0

declare -r digits="${#n}"

for ip in $(jq -rc '.[].NetAddr' "$PEERS_DIR/lachesis_data_dir/peers.json"); do
    ip="${ip%:*}";
    printf -v node_num_p "%0${digits}d" "$node_num"
    cp "$PEERS_DIR/lachesis_data_dir/peers.json" "$BUILD_DIR/lachesis_data_dir/$node_num_p/"

    echo "$ip"
    ((node_num++)) || true
    [ "$node_num" -gt "$n" ] && exit 0
done

# Run multi lachesis

declare debug=0
while getopts "d" opt; do
    case "$opt" in
    d)  debug=1
        ;;
    esac
done

shift $((OPTIND-1))

[ "${1:-}" = "--" ] && shift

if [ "$debug" == 0 ]; then
  GOMAXPROCS=$(($logicalCpuCount - 1)) "$BUILD_DIR/lachesis_$TARGET_OS" run --datadir "$BUILD_DIR/lachesis_data_dir" --store --listen="$node_addr":12000 --log=warn --heartbeat=5s -p "$node_addr":9000 --test --test_n=10 --test_delay=10
  rm -rf "$BUILD_DIR/lachesis_data_dir/"
else
  GOMAXPROCS=$(($logicalCpuCount - 1)) dlv --listen=localhost:37555 --headless=true --api-version=2 --backend=default exec "$BUILD_DIR/lachesis_$TARGET_OS" -- run --datadir "$BUILD_DIR/lachesis_data_dir" --store --listen="$node_addr":12000 --log=warn --heartbeat=5s -p "$node_addr":9000 --test --test_n=10 --test_delay=10
fi

declare -i rc=$?
exit "$rc"
