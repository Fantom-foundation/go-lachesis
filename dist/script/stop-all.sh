#!/bin/bash

# kill all bootnode and lachesis processes
pkill "bootnode"
pkill "lachesis"
pkill "tx-storm"

# remove demo data
#rm -rf ./lachesis-demo-dir/
