#!/bin/bash

# This script will launch a cluster of N Lachesis nodes
# The parameter N = number of nodes to run

# number of nodes N
N=5

# Bootnode
BOOTNODE=../../build/bootnode

# Generate key
echo "Generate bootnode.key"
${BOOTNODE} -genkey bootnode.key

# bootnode's port
bport=3100

echo "Start bootnode with bootnode.key"
bootnode=$( "${BOOTNODE}" -nodekey bootnode.key --addr ":${bport}"  2>/dev/null | head -1 & )
bootnode=${bootnode/":"0"?"/":"${bport}"?"}

echo -e "Bootnode=${bootnode}"

######
EXEC=../../build/lachesis

# default ip using localhost
IP=127.0.0.1
# default port PORT
# the actual ports are PORT+1, PORT+2, etc (4001, 4002, 4003, ... )
PORT=4000
LOCALPORT=3000
WSPORT=3500

# demo directory 
LACHESIS_BASE_DIR=/tmp/lachesis-demo

echo -e "\nStart $N nodes:"
for i in $(seq $N)
do
    port=$((PORT + i))
    localport=$((LOCALPORT + i))
    wsport=$((WSPORT + i))

    ${EXEC} \
	--bootnodes "${bootnode}" \
	--fakenet $i/$N \
	--port ${localport} --rpc --rpcapi "eth,debug,admin,web3,personal,net,txpool,ftm,sfc" --rpcport ${port} \
	--ws --wsaddr="0.0.0.0" --wsport=${wsport} --wsorigins="*" --wsapi="eth,debug,admin,web3,personal,net,txpool,ftm,sfc" \
	--nousb --verbosity 3 \
	--datadir "${LACHESIS_BASE_DIR}/datadir/lach$i" &
    echo -e "Started lachesis client at ${IP}:${port}, pid: $!"
done

echo
echo "Sleep for 10000 seconds..."
sleep 10000
