#!/usr/bin/env bash
cd $(dirname $0)
. ./_params.sh

docker run -d --rm \
    -p 9093:9093 \
    --net=${NETWORK} \
    --name=alertmanager \
    prom/alertmanager
