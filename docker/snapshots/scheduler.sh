#!/bin/bash

set -x

cd $(dirname $0)
mkdir -p /snapshots
touch /snapshots/part.snapshot
ln -s /snapshots /usr/share/nginx/html/snapshots
cp index.html /usr/share/nginx/html/index.html

while true
do
    lachesis --nousb &
    PID=$!

    sleep 2
    while [ $(( $(date +%s)-$(stat -c %Y /snapshots/part.snapshot) )) -lt 5 ]
    do
        sleep 5
    done

    EPOCH_WAS=$(cat /snapshots/was.epoch 2>/dev/null || echo 0)
    while
        EPOCH_NOW=$(lachesis attach --exec='ftm.currentEpoch()')
	[ ${EPOCH_NOW} -le ${EPOCH_WAS} ]
    do
        sleep 5
    done
    #'

    kill ${PID} && wait ${PID}
    lachesis export /snapshots/part.snapshot ${EPOCH_WAS} $((EPOCH_NOW-1))

    if [ -f ${EPOCH_WAS}.snapshot ]
    then
	cp /snapshots/${EPOCH_WAS}.snapshot /snapshots/${EPOCH_NOW}.snapshot
	tail -c +9 /snapshots/part.snapshot >> /snapshots/${EPOCH_NOW}.snapshot
    else
        mv /snapshots/part.snapshot /snapshots/${EPOCH_NOW}.snapshot
    fi
    echo ${EPOCH_NOW} > /snapshots/was.epoch
    
    sed "s/EPOCH/${EPOCH_NOW}/g" index.html > /usr/share/nginx/html/index.html
done
