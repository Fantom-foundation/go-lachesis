#!/usr/bin/env bash

set -x

cd $(dirname $0)
mkdir -p files
touch ./files/part.snapshot
ln -s ./files /usr/share/nginx/html/snapshots
ln -s ./files/index.html /usr/share/nginx/html/index.html
PERIOD_SEC=$(( 60 ))

while true
do
    lachesis --nousb &
    PID=$!

    while [ $(( $(date +%s)-$(stat -c %Y files/part.snapshot) )) -lt ${PERIOD_SEC} ]
    do
        sleep 5
    done

    EPOCH_WAS=$(cat files/was.epoch 2>/dev/null || echo 0)
    while
        EPOCH_NOW=$(lachesis attach --exec='ftm.currentEpoch()')
	[ ${EPOCH_NOW} -le ${EPOCH_WAS} ]
    do
        sleep 5
    done
    
    #'
    kill ${PID} && wait ${PID}
    lachesis export files/part.snapshot ${EPOCH_WAS} $((EPOCH_NOW-1))

    if [ -f files/${EPOCH_WAS}.snapshot ]
    then
	cp files/${EPOCH_WAS}.snapshot files/${EPOCH_NOW}.snapshot
	tail -c +9 files/part.snapshot >> files/${EPOCH_NOW}.snapshot
    else
        cp files/part.snapshot files/${EPOCH_NOW}.snapshot
    fi
    echo ${EPOCH_NOW} > files/was.epoch
    
    sed "s/EPOCH/${EPOCH_NOW}/g" index.html > files/index.html

    find ./files/ -name "*.snapshot" -type f -mtime +1 -exec rm -f {} \;
done
