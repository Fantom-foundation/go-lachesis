# Snapshot service

allows to download the network events as single file to sync your node faster, without p2p load.


## Docker image "lachesis-snapshots"

is an example how to organize snapshot service.
There are lachesis node which exports events snapshots and nginx which shares snapshots for downloading.
See details at [snapshotting.sh](./snapshots/snapshotting.sh).

Build it by `make snapshots` from [Dockerfile.snapshots](./Dockerfile.snapshots).

Run it by (volumes mapping is optional)
```sh
docker run --rm -d \
    -v ${PWD}/files:/snapshots/files \
    -v ${PWD}/node-db:/root/.lachesis \
    -p 127.0.0.1:8080:80 \
    lachesis-snapshots:latest
    
```
and see the download page at http://127.0.0.1:8080.
