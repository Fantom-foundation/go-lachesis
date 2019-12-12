# Introduction

The `dist` folder contains a distribution of binaries and scripts that can be used to run Lachesis nodes. 
The folder includes binaries (`boonode`, `lachesis` and `tx-storm`) and a `script` folder. 


The `script` folder includes:

1. `start-bootnode.sh`: starts a bootnode 

2. `start-demo-nodes.sh`: starts a cluster of demo Lachesis nodes 

3. `start-demo-node.sh`: starts a single Lachesis node

2. `stop.sh`: stops all Lachesis nodes from running. The `lachesis-demo-dir` directory with node data is deleted.

3. `stop-all.sh` : stops all Lachesis nodes and bootnode nodes. The `lachesis-demo-dir` directory with node data is deleted.


Tx-generators 

1.  `./start-txstorm.sh`: Starts the tx-generators, the number of transaction  generated per second are specified on line 37:
`--rate=<number of transactions per second.`

2.  `./stop-txstorm.sh`: Destroys the tx-generators


# Build binaries

To generate the binaries, run `make dist` in the root directory of go-lachesis.

You should now have `lachesis`, `bootnode` and `tx-storm` under the `dist` folder


# How to run new Lachesis nodes with tx-generators

Under `dist` directory:

1. Run `./start-bootnode.sh` : start the bootnode

2. Run `./start-demo-nodes.sh 5`
(In this example, N is set to 5, which is the number of nodes to lauch).

3. Wait for the nodes to connect

3. Run `./start-txstorm.sh 5`

Transactions will be continuously send to the network, until stopped.

# Logs

Data are stored under `dist/lachesis-demo-dir` and logs are generated under `dist/txstorm_logs` for each node

# Stopping

Under `dist` directory

1. Run `stop-txstorm.sh`: Destroy the tx-generators first, if they are running.

2. Run `stop.sh` or `stop-all`: Destroys the network, and deletes the `./lachesis-demo-dir` directory
