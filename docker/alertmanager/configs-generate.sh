#!/usr/bin/env bash

cd $(dirname $0)
. ../_params.sh

cat ./alertmanager.yml.template| SMTP_HOST=${SMTP_HOST} SMTP_FROM=${SMTP_FROM} SMTP_AUTH_USER=${SMTP_AUTH_USER} SMTP_AUTH_PASS=${SMTP_AUTH_PASS} EMAIL_ADMIN=${EMAIL_ADMIN} envsubst > ./alertmanager.yml
