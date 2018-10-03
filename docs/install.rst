.. _install:

Install
=======

From Source
^^^^^^^^^^^

Clone the `repository <https://github.com/andrecronje/lachesis>`__ in the appropriate GOPATH subdirectory:

::

    $ mkdir -p $GOPATH/src/github.com/andrecronje/
    $ cd $GOPATH/src/github.com/andrecronje
    [...]/andrecronje$ git clone https://github.com/andrecronje/lachesis.git


The easiest way to build binaries is to do so in a hermetic Docker container. 
Use this simple command:  

::

	[...]/lachesis$ make dist

This will launch the build in a Docker container and write all the artifacts in  
the build/ folder.  

::
	
    [...]/lachesis$ tree build
    build/
    ├── dist
    │   ├── lachesis_0.1.0_darwin_386.zip
    │   ├── lachesis_0.1.0_darwin_amd64.zip
    │   ├── lachesis_0.1.0_freebsd_386.zip
    │   ├── lachesis_0.1.0_freebsd_arm.zip
    │   ├── lachesis_0.1.0_linux_386.zip
    │   ├── lachesis_0.1.0_linux_amd64.zip
    │   ├── lachesis_0.1.0_linux_arm.zip
    │   ├── lachesis_0.1.0_SHA256SUMS
    │   ├── lachesis_0.1.0_windows_386.zip
    │   └── lachesis_0.1.0_windows_amd64.zip
    └── pkg
        ├── darwin_386
        │   └── lachesis
        ├── darwin_amd64
        │   └── lachesis
        ├── freebsd_386
        │   └── lachesis
        ├── freebsd_arm
        │   └── lachesis
        ├── linux_386
        │   └── lachesis
        ├── linux_amd64
        │   └── lachesis
        ├── linux_arm
        │   └── lachesis
        ├── windows_386
        │   └── lachesis.exe
        └── windows_amd64
            └── lachesis.exe
    
Go Devs
^^^^^^^

Lachesis is written in `Golang <https://golang.org/>`__. Hence, the first step is 
to install **Go version 1.9 or above** which is both the programming language  
and a CLI tool for managing Go code. Go is very opinionated  and will require 
you to `define a workspace <https://golang.org/doc/code.html#Workspaces>`__ 
where all your go code will reside. 

Dependencies  
^^^^^^^^^^^^

Lachesis uses `Glide <http://github.com/Masterminds/glide>`__ to manage 
dependencies. For Ubuntu users:

::

    [...]/lachesis$ curl https://glide.sh/get | sh
    [...]/lachesis$ glide install

This will download all dependencies and put them in the **vendor** folder.

Testing
^^^^^^^

Lachesis has extensive unit-testing. Use the Go tool to run tests:  

::

    [...]/lachesis$ make test

If everything goes well, it should output something along these lines:  

::

    ok      github.com/andrecronje/lachesis/net      0.052s
    ok      github.com/andrecronje/lachesis/common   0.011s
    ?       github.com/andrecronje/lachesis/cmd      [no test files]
    ?       github.com/andrecronje/lachesis/cmd/dummy_client [no test files]
    ok      github.com/andrecronje/lachesis/poset        0.174s
    ok      github.com/andrecronje/lachesis/node     1.699s
    ok      github.com/andrecronje/lachesis/proxy    0.018s
    ok      github.com/andrecronje/lachesis/crypto   0.028s
