#!/bin/bash

# This script will start a cluster of Lachesis nodes that connect to a bootnode

bootnode=$( cat bootenode.txt )
echo -e "Bootnode=${bootnode}"

# number of nodes N
N=$1

######
EXEC=../dist/lachesis

# default ip using localhost
IP=127.0.0.1
# default port PORT
# the actual ports are PORT+1, PORT+2, etc (18541, 18542, 18543, ... )
PORT=18540

# demo directory 
LACHESIS_BASE_DIR=./lachesis-demo-dir

echo -e "\nStart $N nodes:"
for i in $(seq $N)
do
    port=$((PORT + i))
    localport=$((5050 + i))

    ${EXEC} \
	--bootnodes "${bootnode}" \
	--fakenet $i/$N \
	--port ${localport} --rpc --rpcapi "eth,debug,admin,web3" --rpcport ${port} --nousb --verbosity 3 \
	--datadir "${LACHESIS_BASE_DIR}/datadir/lach$i" &
    echo -e "Started lachesis client at ${IP}:${port}, pid: $!"
done



echo
echo "Sleep for 10000 seconds..."
sleep 10000
