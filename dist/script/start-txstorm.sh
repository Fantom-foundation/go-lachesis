#!/bin/bash

# Number of tx-storm agents
N=$1

# Program
PROG=../dist/tx-storm

# default values
TEST_ACCS_START=1000
TEST_ACCS_COUNT=100000

# default ip using localhost
IP=127.0.0.1
# default port PORT
# the actual ports are PORT+1, PORT+2, etc (18541, 18542, 18543, ... )
PORT=18540

TXLOGDIR=./txstorm_logs
mkdir -p ${TXLOGDIR}

# start N tx generators
echo -e "Start $N tx generators:"

for i in $(seq $N)
do
    port=$((PORT + i))
    echo -e "tx-storm $i at port ${port}:"
    f=${TXLOGDIR}/${i}
    
    ${PROG} \
	--num=$i/$N --rate=500 \
	--accs-start=${TEST_ACCS_START} --accs-count=${TEST_ACCS_COUNT} \
	--metrics --verbosity 5 \
	http://${IP}:${port} >${f}.log 2>${f}.err &
done
