#!/bin/bash

# This script will launch a bootnode

BOOTNODE=../dist/bootnode

echo "Generate bootnode.key"
${BOOTNODE} -genkey bootnode.key

bport=3001

echo "Start bootnode with bootnode.key"
bootnode=$( "${BOOTNODE}" -nodekey bootnode.key --addr :${bport} 2>/dev/null | head -1 & )
bootnode=${bootnode/":"0"?"/":"${bport}"?"}

echo ${bootnode} > bootenode.txt

echo -e "Bootnode=${bootnode}"