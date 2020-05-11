#!/usr/bin/env bash
cd $(dirname $0)
. ./_params.sh

echo '--- ${NETWORK} ---' | NETWORK=${NETWORK} envsubst




exit

docker run --rm \
    -p 3000:3000 \
    --net=${NETWORK} \
    --name=grafana \
    grafana/grafana-lachesis
