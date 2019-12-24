#!/bin/bash

# kill all lachesis processes
pkill "lachesis"

# remove demo data
rm -rf /tmp/lachesis-demo/datadir/
