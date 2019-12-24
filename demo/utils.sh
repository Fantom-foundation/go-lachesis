#!/bin/bash

declare -a pids

# number of nodes N
N=3
# number of new nodes M
M=3
# total
T=$((N+M))

# base dir for running demo
LACHESIS_BASE_DIR=/tmp/lachesis-demo

#
PROG=lachesis
EXEC=../build/lachesis

# default ip using localhost
IP=127.0.0.1
# the actual ports are RPCPORT+1, RPCPORT+2, etc (4001, 4002, 4003, ... )
RPCPORT=4000
# actual local ports are LOCALPORT+1, LOCALPORT+2, ...
LOCALPORT=3000
WSPORT=3500

# Function to be called when user press ctrl-c.
ctrl_c() {
    echo "Shell terminating ..."
    for pid in "${pids[@]}"
    do
	echo "Killing ${pid}..."
	kill -9 ${pid}
	# Suppresses "[pid] Terminated ..." message.
	wait ${pid} &>/dev/null
    done
    exit;
}

start_node() {
	local i=$1
	local N=$2
	echo -e "start_node node $i:"

    rpcport=$((RPCPORT + i))
	localport=$((LOCALPORT + i))
	wsport=$((WSPORT + i))

	echo -e "port=${rpcport}, localport=${localport} "

    ${EXEC} \
	--fakenet $i/$N \
	--port ${localport} --rpc --rpcapi "eth,debug,admin,web3,personal,net,txpool,ftm,sfc" --rpcport ${rpcport} \
	--ws --wsaddr="0.0.0.0" --wsport=${wsport} --wsorigins="*" --wsapi="eth,debug,admin,web3,personal,net,txpool,ftm,sfc" \
	--nousb --verbosity 3 --metrics \
	--datadir "${LACHESIS_BASE_DIR}/datadir/lach$i" &
	pids+=($!)
    echo -e "Started lachesis client at ${IP}:${rpcport}, pid: $!"
    echo -e "\n"
}

start_nodes() {
	local i=$1
	local j=$2
	local N=$3
	echo -e "\nStart $N nodes:\n"
	for i in $(seq $i $j)
	do
	    start_node $i $N
	done
}

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

connect_pair() {
    local from=$1
    local to=$2

	echo " getting node-${to} address:"
	url=${IP}:$((RPCPORT + to))
	echo "    at url: ${url}"

    enode=$(attach_and_exec ${url} 'admin.nodeInfo.enode')
    echo "    p2p address = ${enode}"

    echo " connecting node-${from} to node-${to}:"
    url=${IP}:$((RPCPORT + from))
    echo "    at url: ${url}"

    res=$(attach_and_exec ${url} "admin.addPeer(${enode})")
    echo "    result = ${res}"

	return 0
}

connect_nodes() {
	local from=$1
	local to=$2

	echo -e "\nConnect nodes to ring:\n"
	for i in $(seq ${from} $((to-1)))
	do
	    j=$((i + 1))
	    conn=$(connect_pair $i $j)
	done
	conn=$(connect_pair ${to} ${from})
}


