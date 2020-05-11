#!/usr/bin/env bash

cd $(dirname $0)
. ../_params.sh

cat ./grafana/datasource.yaml.template| PROMETHEUS_URL=${PROMETHEUS_URL} envsubst > ./datasource.yaml
cat ./grafana/grafana-dashboard.json.template| PROMETHEUS_JOB=${PROMETHEUS_JOB} envsubst > ./lachesis-dashboard.json
