#!/usr/bin/env bash

cd $(dirname $0)
. ../_params.sh

cat ./prometheus.yml.template| PROMETHEUS_JOB=${PROMETHEUS_JOB} TARGET_HOST=${TARGET_HOST} ALERTMANAGER_HOST=${ALERTMANAGER_HOST} envsubst > ./prometheus.yml
