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
# the actual ports are RPCPORT+1, RPCPORT+2, etc (4001, 4002, 4003, ... )
RPCPORT=4000
LOCALPORT=3000
WSPORT=3500

# demo directory 
LACHESIS_BASE_DIR=/tmp/lachesis-dem

echo -e "\nStart $N nodes:"
for i in $(seq $N)
do
    rpcport=$((RPCPORT + i))
    localport=$((LOCALPORT + i))
    wsport=$((WSPORT + i))

    ${EXEC} \
	--bootnodes "${bootnode}" \
	--fakenet $i/$N \
	-port ${localport} --rpc --rpcapi "eth,debug,admin,web3,personal,net,txpool,ftm,sfc" --rpcport ${rpcport} \
	--ws --wsaddr="0.0.0.0" --wsport=${wsport} --wsorigins="*" --wsapi="eth,debug,admin,web3,personal,net,txpool,ftm,sfc" \
	--nousb --verbosity 3 \
	--datadir "${LACHESIS_BASE_DIR}/datadir/lach$i" &
    echo -e "Started lachesis client at ${IP}:${rpcport}, pid: $!"
done

echo
echo "Sleep for 10000 seconds..."
sleep 10000
