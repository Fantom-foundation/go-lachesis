#!/bin/bash

# This script will launch a bootnode

BOOTNODE=../dist/bootnode

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