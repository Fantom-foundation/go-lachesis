#!/bin/bash

# This script will launch a cluster of N Lachesis nodes
# The parameter N = number of nodes to run

. ./utils.sh

# Generate bootnode.key
echo "Generate bootnode.key"
${BOOTNODE} -genkey bootnode.key

# bootnode's port
bport=3100

# Start bootnode
echo "Start bootnode with bootnode.key"
bootnode=$( "${BOOTNODE}" -nodekey bootnode.key --addr :${bport} 2>/dev/null | head -1 & )
bootnode=${bootnode/":"0"?"/":"${bport}"?"}

echo ${bootnode} > bootenode.txt

echo -e "Bootnode=${bootnode}"

######
echo -e "\nStart $N nodes:"
for i in $(seq $N)
do
    rpcport=$((RPCPORT + i))
    localport=$((LOCALPORT + i))
    wsport=$((WSPORT + i))

    ${EXEC} \
	--bootnodes "${bootnode}" \
	--fakenet $i/$T \
	--port ${localport} --rpc --rpcapi "eth,debug,admin,web3,personal,net,txpool,ftm,sfc" --rpcport ${rpcport} \
	--ws --wsaddr="0.0.0.0" --wsport=${wsport} --wsorigins="*" --wsapi="eth,debug,admin,web3,personal,net,txpool,ftm,sfc" \
	--nousb --verbosity 3 --metrics \
	--datadir "${LACHESIS_BASE_DIR}/datadir/lach$i" &
    echo -e "Started lachesis client at ${IP}:${rpcport}, pid: $!"
done

echo
echo "Sleep for 10000 seconds..."
sleep 10000
