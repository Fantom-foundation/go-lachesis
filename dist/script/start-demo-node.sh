#!/bin/bash

# This script will start a new Lachesis node and connects it to a running bootnode

bootnode=$( cat bootenode.txt )
echo -e "Bootnode=${bootnode}"

######
EXEC=../dist/lachesis

# Usage: ./scripts/start-demo-node.sh 127.0.0.1 2 4001 3001 3501


# ip using localhost
IP=$1
fakenet=$2
# port RPCPORT such as (4001, 4002, 4003, ... )
rpcport=$3
# local port can be 3001
localport=$4
wsport=$5

# demo directory 
LACHESIS_BASE_DIR=/tmp/lachesis-demo

${EXEC} \
    -bootnodes "${bootnode}" \
	--fakenet ${fakenet}/100 \
	--port ${localport} --rpc --rpcapi "eth,debug,admin,web3,personal,net,txpool,ftm,sfc" --rpcport ${rpcport} \
	--ws --wsaddr="0.0.0.0" --wsport=${wsport} --wsorigins="*" --wsapi="eth,debug,admin,web3,personal,net,txpool,ftm,sfc" \
	--nousb --verbosity 3 \
	--datadir "${LACHESIS_BASE_DIR}/datadir/lach${fakenet}" &
echo -e "Started lachesis client at ${IP}:${rpcport}, pid: $!"
