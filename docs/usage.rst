.. _usage:

Usage
=====

In this section we will guide you through deploying an application on top of
Lachesis. Lachesis comes with the Dummy application which is used in this
demonstration. It is a simple chat application where participants write
messages on a channel and Lachesis guarantees that everyone sees the same messages
in the same order.

Docker
------

We have provided a series of scripts to bootstrap a demo. Let us first use the
easy method to view the demo and then we will take a closer look at what is
happening behind the scenes.

Make sure you have `Docker <https://docker.com>`__ installed.

The demo will pull Docker images from our `official public Docker registry
<https://hub.docker.com/u/Fantom-foundation/>`__

::

    [...]/lachesis$ cd demo
    [...]/lachesis/demo$ make


Once the testnet is started, a script is automatically launched to monitor
consensus figures:

::

    consensus_events:180 consensus_transactions:40 events_per_second:0.00 id:1 last_block_index:3 last_consensus_round:17 num_peers:3 round_events:7 rounds_per_second:0.00 state:Babbling sync_rate:1.00 transaction_pool:0 undetermined_events:18
    consensus_events:180 consensus_transactions:40 events_per_second:0.00 id:3 last_block_index:3 last_consensus_round:17 num_peers:3 round_events:7 rounds_per_second:0.00 state:Babbling sync_rate:1.00 transaction_pool:0 undetermined_events:20
    consensus_events:180 consensus_transactions:40 events_per_second:0.00 id:2 last_block_index:3 last_consensus_round:17 num_peers:3 round_events:7 rounds_per_second:0.00 state:Babbling sync_rate:1.00 transaction_pool:0 undetermined_events:21
    consensus_events:180 consensus_transactions:40 events_per_second:0.00 id:0 last_block_index:3 last_consensus_round:17 num_peers:3 round_events:7 rounds_per_second:0.00 state:Babbling sync_rate:1.00 transaction_pool:0 undetermined_events:20

Running ``docker ps -a`` will show you that 9 docker containers have been launched:

::

    [...]/lachesis/demo$ docker ps -a
    CONTAINER ID        IMAGE                    COMMAND                  CREATED             STATUS              PORTS                   NAMES
    ba80ef275f22        Fantom-foundation/watcher   "/watch.sh"              48 seconds ago      Up 7 seconds                                watcher
    4620ed62a67d        Fantom-foundation/dummy     "dummy '--name=client"   49 seconds ago      Up 48 seconds       1339/tcp                client4
    847ea77bd7fc        Fantom-foundation/go-lachesis    "lachesis run --cache_s"   50 seconds ago      Up 49 seconds       80/tcp, 1337-1338/tcp   node4
    11df03bf9690        Fantom-foundation/dummy     "dummy '--name=client"   51 seconds ago      Up 50 seconds       1339/tcp                client3
    00af002747ca        Fantom-foundation/go-lachesis    "lachesis run --cache_s"   52 seconds ago      Up 50 seconds       80/tcp, 1337-1338/tcp   node3
    b2011d3d65bb        Fantom-foundation/dummy     "dummy '--name=client"   53 seconds ago      Up 51 seconds       1339/tcp                client2
    e953b50bc1db        Fantom-foundation/go-lachesis    "lachesis run --cache_s"   53 seconds ago      Up 52 seconds       80/tcp, 1337-1338/tcp   node2
    0c9dd65de193        Fantom-foundation/dummy     "dummy '--name=client"   54 seconds ago      Up 53 seconds       1339/tcp                client1
    d1f4e5008d4d        Fantom-foundation/go-lachesis    "lachesis run --cache_s"   55 seconds ago      Up 54 seconds       80/tcp, 1337-1338/tcp   node1


Indeed, each node is comprised of an App and a Lachesis node (cf Design section).
The ``watcher`` container monitors consensus figures.

Run the ``demo`` script to play with the ``Dummy App`` which is a simple chat application
powered by the Lachesis consensus platform:

::

    [...]/lachesis/demo$ make demo

.. image:: assets/demo.png

Finally, stop the testnet:

::

    [...]/lachesis/demo$ make stop

Manual Setup
------------

The above scripts hide a lot of the complications involved in setting up a
Lachesis network. They generate the configuration files automatically, copy them
to the right places and launch the nodes in Docker containers. We recommend
looking at these scripts closely to understand how the demo works. Here, we will
attempt to explain the individual steps that take place behind the scenes.

Configuration
~~~~~~~~~~~~~

Lachesis reads configuration from the directory specified by the ``datadir`` flag
which defaults to ``~/.lachesis`` on linux/osx. This directory must contain two
files:

 - ``peers.json``  : Lists all the participants in the network.
 - ``priv_key.pem``: Contains the private key of the validator runnning the node.

Every participant has a cryptographic key-pair that is used to encrypt, sign and
verify messages. The private key is secret but the public key is used by other
nodes to verify messages signed with the private key. The encryption scheme used
by Lachesis is ECDSA with the P256 curve.

To run a Lachesis network, it is necessary to predefine who the participants are
going to be. Each participant will generate a key-pair and decide which network
address it is going to be using for the Lachesis protocol. Someone, or some
process, then needs to aggregate the public keys and network addresses of all
participants into a single file - the peers.json file. Every participant uses a
copy of the same peers.json file. Lachesis will read that file to identify the
participants in the network, communicate with them and verify their
cryptographic signatures.

To generate key-pairs in a format usable by Lachesis, we have created the
``keygen`` command. It is left to the user to derive a scheme to produce the
configuration files but the docker demo scripts are a good place to start.

So let us say I want to participate in a Lachesis network. I am going to start by
running ``lachesis keygen`` to create a key-pair:

::

  lachesis keygen
  Your private key has been saved to: /home/martin/.lachesis/priv_key.pem
  Your public key has been saved to: /home/martin/.lachesis/key.pub

The private key looks something like this:

::

  -----BEGIN EC PRIVATE KEY-----
  MHcCAQEEIJ3orqofiSXu07mD+f46gZFK3EKSTqhXsbLVmA/aLmyqoAoGCCqGSM49
  AwEHoUQDQgAEXgNNc8hJdWrntlFcpg2WpakRsTpNi0W8DgsC7bRQCd9szAdO6298
  Z5V0D5k2ZO3ulw+KcXyJNE+EN/QSvfDRfA==
  -----END EC PRIVATE KEY-----

and the corresponding public key looks like this:

::

  0x045E034D73C849756AE7B6515CA60D96A5A911B13A4D8B45BC0E0B02EDB45009DF6CCC074EEB6F7C6795740F993664EDEE970F8A717C89344F8437F412BDF0D17C

**DO NOT REUSE THESE KEYS**

Next, I am going to copy the public key (key.pub) and communicate it to whoever
is responsible for producing the peers.json file. At the same time, I will tell
them that I am going to be listening on 172.77.5.2:1337.

Suppose three other people do the same thing. The resulting peers.json file
could look something like this:

::

    [
	{
		"NetAddr":"172.77.5.1:1337",
		"PubKeyHex":"0x0471AEE3CAE4E8442D37C9F5481FB32C4531511988652DF923B79ED4ED992021183D31E0F6FBFE96D89B6D03D7250292DFECD4FC414D83A5C38FA3FAD0D8572864"
	},
	{
		"NetAddr":"172.77.5.2:1337",
		"PubKeyHex":"0x045E034D73C849756AE7B6515CA60D96A5A911B13A4D8B45BC0E0B02EDB45009DF6CCC074EEB6F7C6795740F993664EDEE970F8A717C89344F8437F412BDF0D17C"
	},
	{
		"NetAddr":"172.77.5.3:1337",
		"PubKeyHex":"0x047CCCD40D90B331C64CE27911D3A31AF7DC16C1EA6D570FDC2120920663E0A678D7B5D0C19B6A77FEA829F8198F4F487B68206B93B7AD17D7C49CA7E0164D0033"
	},
	{
		"NetAddr":"172.77.5.4:1337",
		"PubKeyHex":"0x0406CB5043E7337700E3B154993C872B1C61A84B1A739528C4A10135A3D64939C094B4A999BD21C3D5E9E9ECF15B202414F073795C9483B2F51ADA7EE59EB5EAC4"
	}
    ]

Now everyone is going to take a copy of this peers.json file and put it in a
folder together with the priv_key.pem file they generated in the previous step.
That is the folder that they need to specify as the datadir when they run
Lachesis.

Lachesis Executable
-----------------

Let us take a look at the help provided by the Lachesis CLI:

::

  Run node

  Usage:
    lachesis run [flags]

  Flags:
        --cache-size int          Number of items in caches (default 500)
    -c, --client-connect string   IP:Port to connect to client (default "127.0.0.1:1339")
        --datadir string          Top-level directory for configuration and data (default "/home/martin/.lachesis")
        --heartbeat duration      Time between gossips (default 1s)
    -h, --help                    help for run
    -l, --listen string           Listen IP:Port for lachesis node (default ":1337")
        --log string              debug, info, warn, error, fatal, panic
        --max-pool int            Connection pool size max (default 2)
    -p, --proxy-listen string     Listen IP:Port for lachesis proxy (default "127.0.0.1:1338")
    -s, --service-listen string   Listen IP:Port for HTTP service
        --standalone              Do not create a proxy
        --store                   Use badgerDB instead of in-mem DB
        --sync-limit int          Max number of events for sync (default 100)
    -t, --timeout duration        TCP Timeout (default 1s)


So we have just seen what the ``datadir`` flag does. The ``listen`` flag
corresponds to the NetAddr in the peers.json file; that is the endpoint that
Lachesis uses to communicate with other Lachesis nodes.

As we explained in the architecture section, each Lachesis node works in
conjunction with an application for which it orders transactions. When Lachesis
and the application are connected by a TCP interface, we specify two other
endpoints:

 - ``proxy-listen``  : where Lachesis listens for transactions from the App
 - ``client-connect`` : where the App listens for transactions from Lachesis

We can also specify where Lachesis exposes its HTTP API providing information on
the Poset and Blockchain data store. This is controlled by the optional
``service-listen`` flag.

Finally, we can choose to run Lachesis with a database backend or only with an
in-memory cache. With the ``store`` flag set, Lachesis will look for a database
file in ``datadir``/babdger_db. If the file exists, the node will load the
database and bootstrap itself to a state consistent with the database and it
will be able to proceed with the consensus algorithm from there. If the file
does not exist yet, it will be created and the node will start from a clean
state.

Here is how the Docker demo starts Lachesis nodes together wth the Dummy
application:

::

    for i in $(seq 1 $N)
    do
        docker run -d --name=client$i --net=lachesisnet --ip=172.77.5.$(($N+$i)) -it Fantom-foundation/dummy:0.4.0 \
        --name="client $i" \
        --client-listen="172.77.5.$(($N+$i)):1339" \
        --proxy-connect="172.77.5.$i:1338" \
        --discard \
        --log="debug"
    done

    for i in $(seq 1 $N)
    do
        docker create --name=node$i --net=lachesisnet --ip=172.77.5.$i Fantom-foundation/go-lachesis:0.4.0 run \
        --cache-size=50000 \
        --timeout=200ms \
        --heartbeat=10ms \
        --listen="172.77.5.$i:1337" \
        --proxy-listen="172.77.5.$i:1338" \
        --client-connect="172.77.5.$(($N+$i)):1339" \
        --service-listen="172.77.5.$i:80" \
        --sync-limit=1000 \
        --store \
        --log="debug"

        docker cp $MPWD/conf/node$i node$i:/.lachesis
        docker start node$i
    done

Stats, blocks and Logs
----------------------

Once a node is up and running, we can call the ``stats`` endpoint exposed by the
HTTP service:

::

    curl -s http://172.77.5.1:80/stats

or request to see a specific block:

::

    curl -s http://172.77.5.1:80/block/1

Or we can look at the logs produced by Lachesis:

::

    docker logs node1
