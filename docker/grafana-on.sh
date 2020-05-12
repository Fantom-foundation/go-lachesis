#!/usr/bin/env bash
cd $(dirname $0)
. ./_params.sh

docker volume create grafana-data

docker run -d --rm \
    -p 3000:3000 \
    --net=${NETWORK} \
    --name=grafana \
    --mount source=grafana-data,target=/var/lib/grafana \
    grafana-lachesis
