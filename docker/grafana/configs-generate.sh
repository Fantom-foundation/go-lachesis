#!/usr/bin/env bash

cd $(dirname $0)
. ../_params.sh

cat ./datasource.yaml.template| PROMETHEUS_URL=${PROMETHEUS_URL} envsubst > ./datasource.yaml
cat ./grafana-dashboard.json.template| PROMETHEUS_JOB=${PROMETHEUS_JOB} envsubst > ./lachesis-dashboard.json
