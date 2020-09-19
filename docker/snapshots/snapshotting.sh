#!/usr/bin/env bash

set -x

cd $(dirname $0)
mkdir -p files
touch ./files/part.snapshot
ln -s ${PWD}/files /usr/share/nginx/html/snapshots
rm -f /usr/share/nginx/html/index.html
ln -s ${PWD}/files/index.html /usr/share/nginx/html/index.html
PERIOD_SEC=$(( 60 ))

while true
do
    # start the node
    lachesis --nousb &
    PID=$!

    # wait for period
    while
	AGE_SEC=$(( $(date +%s)-$(stat -c %Y files/part.snapshot) ))
	[ ${AGE_SEC} -lt ${PERIOD_SEC} ]
    do
        sleep 5
    done

    # wait for new epoch
    EPOCH_WAS=$(cat files/was.epoch 2>/dev/null || echo 0)
    while
        EPOCH_NOW=$(lachesis attach --exec='ftm.currentEpoch()') #'
	[ ${EPOCH_NOW} -le ${EPOCH_WAS} ]
    do
        sleep 5
    done
    
    # stop the node and export new events snapshot
    kill ${PID} && wait ${PID}
    lachesis export files/part.snapshot ${EPOCH_WAS} $((EPOCH_NOW-1))

    if [ -f files/${EPOCH_WAS}.snapshot ]
    then
	# append prev file to the snapshot
	cp files/${EPOCH_WAS}.snapshot files/${EPOCH_NOW}.snapshot
	tail -c +9 files/part.snapshot >> files/${EPOCH_NOW}.snapshot
    else
	# snapshot file is the first
        cp files/part.snapshot files/${EPOCH_NOW}.snapshot
    fi
    echo ${EPOCH_NOW} > files/was.epoch
    
    # update link to the latest snapshot file
    sed "s/EPOCH/${EPOCH_NOW}/g" index.html > files/index.html

    # remove snapshots except the 5 last
    ls -t1 files/*.snapshot | tail -n +5 | while read F
    do
	rm -f ${F}
    done
done
