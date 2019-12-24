#!/bin/bash

# Number of tx-storm agents
N=5

# Program
PROG=../build/tx-storm

# default values
PORT_BASE=3000
RPCP_BASE=4000
WSP_BASE=4500

TEST_ACCS_START=1000
TEST_ACCS_COUNT=100000

# default ip using localhost
IP=127.0.0.1
# the actual ports are RPCPORT+1, RPCPORT+2, etc (4001, 4002, 4003, ... )
RPCPORT=4000

TXLOGDIR=/tmp/txstorm_logs
mkdir -p ${TXLOGDIR}

# start N tx generators
echo -e "Start $N tx generators:"

for i in $(seq $N)
do
    rpcport=$((RPCPORT + i))
    echo -e "tx-storm $i at port ${rpcport}:"
    f=${TXLOGDIR}/${i}
    
    ${PROG} \
	--num=$i/$N --rate=10 \
	--accs-start=${TEST_ACCS_START} --accs-count=${TEST_ACCS_COUNT} \
	--metrics --verbosity 5 \
	http://${IP}:${rpcport} >${f}.log 2>${f}.err &
done
