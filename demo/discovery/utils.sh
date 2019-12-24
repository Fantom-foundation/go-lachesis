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

# Bootnode
BOOTNODE=../../build/bootnode

######
EXEC=../../build/lachesis

# default ip using localhost
IP=127.0.0.1
# the actual ports are RPCPORT+1, RPCPORT+2, etc (4001, 4002, 4003, ... )
RPCPORT=4000
LOCALPORT=3000
WSPORT=3500