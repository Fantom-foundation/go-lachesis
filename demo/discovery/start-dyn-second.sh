#!/bin/bash

# This script will launch a cluster of N Lachesis nodes
# using 
# The parameter N = number of nodes to run

. ./utils.sh

bootnode=$( cat bootenode.txt )

echo -e "Bootnode=${bootnode}"

######
# demo directory 
LACHESIS_BASE_DIR=/tmp/lachesis-demo

echo -e "\nStart $M nodes:"
for j in $(seq $M)
do
	i=$(( M+j))
	
    port=$((PORT + i))
    localport=$((5050 + i))

    ${EXEC} \
	--bootnodes "${bootnode}" \
	--fakenet $i/$T \
	--port ${localport} --rpc --rpcapi "eth,debug,admin,web3" --rpcport ${port} --nousb --verbosity 3 \
	--datadir "${LACHESIS_BASE_DIR}/datadir/lach$i" &
    echo -e "Started lachesis client at ${IP}:${port}, pid: $!"
done


echo
echo "Sleep for 10000 seconds..."
sleep 10000
