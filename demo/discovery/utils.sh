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

FVM=${GOPATH}/src/github.com/Fantom-foundation/go-ethereum
BOOTNODE=${FVM}/build/bin/bootnode

######
EXEC=../../build/lachesis

# default ip using localhost
IP=127.0.0.1
# default port PORT
# the actual ports are PORT+1, PORT+2, etc (18541, 18542, 18543, ... )
PORT=18540

# demo directory 
LACHESIS_BASE_DIR=/tmp/lachesis-demo