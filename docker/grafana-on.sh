#!/usr/bin/env bash
cd $(dirname $0)
. ./_params.sh

docker run --rm \
    -p 3000:3000 \
    --net=lachesis \
    --name=grafana \
    grafana/grafana
