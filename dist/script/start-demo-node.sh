#!/bin/bash

# This script will start a new Lachesis node and connects it via an bootnode

bootnode=$( cat bootenode.txt )
echo -e "Bootnode=${bootnode}"

######
EXEC=../dist/lachesis

# Usage: ./scripts/start-demo-node.sh 127.0.0.1 2 18541 5051


# ip using localhost
IP=$1
fakenet=$2
# port PORT such as (18541, 18542, 18543, ... )
port=$3
# local port can be 5050
localport=$4

# demo directory 
LACHESIS_BASE_DIR=./lachesis-demo-dir

${EXEC} \
    -bootnodes "${bootnode}" \
	--fakenet ${fakenet}/100 \
	--port ${localport} --rpc --rpcapi "eth,debug,admin,web3" --rpcport ${port} --nousb --verbosity 3 \
	--datadir "${LACHESIS_BASE_DIR}/datadir/lach${fakenet}" &
echo -e "Started lachesis client at ${IP}:${port}, pid: $!"
