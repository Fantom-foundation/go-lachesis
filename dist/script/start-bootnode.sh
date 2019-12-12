#!/bin/bash

# This script will launch a bootnode

BOOTNODE=../dist/bootnode

echo "Generate bootnode.key"
${BOOTNODE} -genkey bootnode.key

echo "Start bootnode with bootnode.key"
bootnode=$( "${BOOTNODE}" -nodekey bootnode.key 2>/dev/null | head -1 & )
echo ${bootnode} > bootenode.txt

echo -e "Bootnode=${bootnode}"