#!/bin/bash

# This script will launch a cluster of Lachesis nodes
# The parameter N = number of nodes to run

# number of nodes N
N=5

#
PROG=lachesis
EXEC=../build/lachesis

# default ip using localhost
IP=127.0.0.1
# the actual ports are RPCPORT+1, RPCPORT+2, etc (4001, 4002, 4003, ... )
RPCPORT=4000
LOCALPORT=3000
WSPORT=3500

LACHESIS_BASE_DIR=/tmp/lachesis-demo


echo -e "\nStart $N nodes:"
for i in $(seq $N)
do
    rpcport=$((RPCPORT + i))
    localport=$((LOCALPORT + i))
    wsport=$((WSPORT + i))
    proport=$((9000 + i))
    
    ${EXEC} \
	--fakenet $i/$N \
	--port ${localport} --rpc --rpcapi "eth,debug,admin,web3,personal,net,txpool,ftm,sfc" --rpcport ${rpcport} \
	--ws --wsaddr="0.0.0.0" --wsport=${wsport} --wsorigins="*" --wsapi="eth,debug,admin,web3,personal,net,txpool,ftm,sfc" \
	--nousb --verbosity=3 --metrics \
	--metrics.prometheus.endpoint ":${proport}"\
	--datadir "${LACHESIS_BASE_DIR}/datadir/lach$i" &
    echo -e "Started lachesis client at ${IP}:${rpcport}"
done



attach_and_exec() {
    local URL=$1
    local CMD=$2

    for attempt in $(seq 20)
    do
        if (( attempt > 5 ));
        then
            echo "  - attempt ${attempt}: " >&2
        fi;

        res=$("${EXEC}" --exec "${CMD}" attach http://${URL} 2> /dev/null)
        if [ $? -eq 0 ]
        then
            #echo "success" >&2
            echo $res
            return 0
        else
            #echo "wait" >&2
            sleep 1
        fi
    done
    echo "failed RPC connection to ${NAME}" >&2
    return 1
}


echo -e "\nConnect nodes to ring:\n"
for i in $(seq $N)
do
    j=$((i % N + 1))

    echo " getting node-$j address:"
	url=${IP}:$((RPCPORT + j))
	echo "    at url: ${url}"

    enode=$(attach_and_exec ${url} 'admin.nodeInfo.enode')
    echo "    p2p address = ${enode}"

    echo " connecting node-$i to node-$j:"
    url=${IP}:$((RPCPORT + i))
    echo "    at url: ${url}"

    res=$(attach_and_exec ${url} "admin.addPeer(${enode})")
    echo "    result = ${res}"
done


echo
echo "Sleep for 10000 seconds..."
sleep 10000
