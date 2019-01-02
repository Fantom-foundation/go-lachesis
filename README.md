# Lachesis 

aBFT Consensus platform for distributed applications.

## Build Details

[![version](https://img.shields.io/github/tag/Fantom-foundation/go-lachesis.svg?style=flat-square&logo=github
)](https://github.com/Fantom-foundation/go-lachesis/releases/latest)  
[![Build Status](https://img.shields.io/travis/Fantom-foundation/go-lachesis.svg?style=flat-square&logo=travis)](https://travis-ci.org/Fantom-foundation/go-lachesis)  
[![appveyor](https://img.shields.io/appveyor/ci/andrecronje/go-lachesis.svg?style=flat-square&logo=appveyor)](https://ci.appveyor.com/project/andrecronje/go-lachesis)  
[![license](https://img.shields.io/github/license/Fantom-foundation/go-lachesis.svg?style=flat-square&logo=github)](LICENSE.md)  
[![libraries.io dependencies](https://img.shields.io/librariesio/github/Fantom-foundation/go-lachesis.svg?style=flat-square&logo=librariesio)](https://libraries.io/github/Fantom-foundation/go-lachesis)  

## Code Quality

[![Go Report Card](https://goreportcard.com/badge/github.com/Fantom-foundation/go-lachesis?style=flat-square&logo=goreportcard)](https://goreportcard.com/report/github.com/Fantom-foundation/go-lachesis)  
[![GolangCI](https://golangci.com/badges/github.com/Fantom-foundation/go-lachesis.svg?style=flat-square&logo=golangci)](https://golangci.com/r/github.com/Fantom-foundation/go-lachesis)   
[![Code Climate Maintainability Grade](https://img.shields.io/codeclimate/maintainability/Fantom-foundation/go-lachesis.svg?style=flat-square&logo=codeclimate)](https://codeclimate.com/github/Fantom-foundation/go-lachesis)  
[![Code Climate Maintainability](https://img.shields.io/codeclimate/maintainability-percentage/Fantom-foundation/go-lachesis.svg?style=flat-square&logo=codeclimate)](https://codeclimate.com/github/Fantom-foundation/go-lachesis)  
[![Code Climate Technical Dept](https://img.shields.io/codeclimate/tech-debt/Fantom-foundation/go-lachesis.svg?style=flat-square&logo=codeclimate)](https://codeclimate.com/github/Fantom-foundation/go-lachesis)  
[![Codacy code quality](https://img.shields.io/codacy/grade/c8c27910210f4b23bcbbe8c60338b1d5.svg?style=flat-square&logo=codacy)](https://app.codacy.com/project/andrecronje/go-lachesis/dashboard)  
[![cii best practices](https://img.shields.io/cii/level/2409.svg?style=flat-square&logo=cci)](https://bestpractices.coreinfrastructure.org/en/projects/2409)  
[![cii percentage](https://img.shields.io/cii/percentage/2409.svg?style=flat-square&logo=cci)](https://bestpractices.coreinfrastructure.org/en/projects/2409)  
  
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square&logo=godoc)](https://godoc.org/github.com/Fantom-foundation/go-lachesis)   

[Documentation](https://github.com/Fantom-foundation/fantom-documentation/wiki).  

[![Sonarcloud](https://sonarcloud.io/api/project_badges/quality_gate?project=Fantom-foundation_go-lachesis)](https://sonarcloud.io/dashboard?id=Fantom-foundation_go-lachesis)  
  
## GitHub


[![Commit Activity](https://img.shields.io/github/commit-activity/w/Fantom-foundation/go-lachesis.svg?style=flat-square&logo=github)](https://github.com/Fantom-foundation/go-lachesis/commits/master)  
[![Last Commit](https://img.shields.io/github/last-commit/Fantom-foundation/go-lachesis.svg?style=flat-square&logo=github)](https://github.com/Fantom-foundation/go-lachesis/commits/master)  
[![Contributors](https://img.shields.io/github/contributors/Fantom-foundation/go-lachesis.svg?style=flat-square&logo=github)](https://github.com/Fantom-foundation/go-lachesis/graphs/contributors)  
[![Issues][github-issues-image]][github-issues-url]  
[![LoC](https://tokei.rs/b1/github/Fantom-foundation/go-lachesis?category=lines)](https://github.com/Fantom-foundation/go-lachesis)  

[![Throughput Graph](https://graphs.waffle.io/Fantom-foundation/go-lachesis/throughput.svg)](https://waffle.io/Fantom-foundation/go-lachesis/metrics/throughput)  

## Social

[![](https://img.shields.io/gitter/room/nwjs/nw.js.svg?style=flat-square)](https://gitter.im/fantom-foundation)    
[![twitter][twitter-image]][twitter-url]  


[codecov-image]: https://codecov.io/gh/fantom-foundation/go-lachesis/branch/master/graph/badge.svg
[codecov-url]: https://codecov.io/gh/fantom-foundation/go-lachesis
[twitter-image]: https://img.shields.io/twitter/follow/FantomFDN.svg?style=social
[twitter-url]: https://twitter.com/intent/follow?screen_name=FantomFDN
[github-issues-image]: https://img.shields.io/github/issues/Fantom-foundation/go-lachesis.svg?style=flat-square&logo=github
[github-issues-url]: https://github.com/Fantom-foundation/go-lachesis/issues

# Features
- [x] k-node selection
- [x] k-parent EventBlock
- [x] EventBlock merge
- [ ] Lachesis consensus
    - [x] Dominators
    - [x] Self Dominators
    - [x] Atropos
    - [x] Clotho
    - [x] Frame
    - [x] Frame Received
    - [x] Dominated
    - [x] Lamport Timestamp
    - [x] Atropos Consensus Time
    - [x] Consensus Timestamp
    - [x] Ordering on same Consensus Timestamp (Lamport Timestamp)
    - [ ] Ordering on same Lamport Timestamp (Flag Table)
    - [x] Ordering on same Flag Table (Signature XOR)
    - [x] Transaction submit
    - [x] Consensus Transaction output
    - [x] Dynamic participants
        - [x] Peer add
        - [x] Peer Remove
- [x] Caching for performances
- [x] Sync
- [x] Event Signature
- [ ] Transaction validation
- [ ] Optimum Network pruning

## Dev

### Docker

Create an 3 node lachesis cluster with:

    n=3 BUILD_DIR="$PWD" ./scripts/docker/scale.bash

### Dependencies

  - [Docker](https://www.docker.com/get-started)
  - [jq](https://stedolan.github.io/jq)
  - [Bash](https://www.gnu.org/software/bash)
  - [git](https://git-scm.com)
  - [Go](https://golang.org)
  - [Glide](https://glide.sh)
  - [batch-ethkey](https://github.com/SamuelMarks/batch-ethkey) with: `go get -v github.com/SamuelMarks/batch-ethkey`
  - [protocol buffers 3](https://github.com/protocolbuffers/protobuf), with: installation of [a release]([here](https://github.com/protocolbuffers/protobuf/releases)) & `go get -u github.com/golang/protobuf/protoc-gen-go`

### Protobuffer 3

This project uses protobuffer 3 for the communication between posets.
To use it, you have to install both `protoc` and the plugin for go code
generation.

Once the stack is setup, you can compile the proto messages by
running this command:

```bash
make proto
```

### Lachesis and dependencies
Clone the [repository](https://github.com/Fantom-foundation/go-lachesis) in the appropriate
GOPATH subdirectory:

```bash
$ d="$GOPATH/src/github.com/Fantom-foundation"
$ mkdir -p "$d"
$ git clone https://github.com/Fantom-foundation/go-lachesis.git "$d"
```
Lachesis uses [Glide](http://github.com/Masterminds/glide) to manage dependencies.

```bash
$ curl https://glide.sh/get | sh
$ cd "$GOPATH/src/github.com/Fantom-foundation" && glide install
```
This will download all dependencies and put them in the **vendor** folder.

### Other requirements

Bash scripts used in this project assume the use of GNU versions of coreutils.
Please ensure you have GNU versions of these programs installed:-

example for macos:
```
# --with-default-names makes the `sed` and `awk` commands default to gnu sed and gnu awk respectively.
brew install gnu-sed gawk --with-default-names
```

### Testing

Lachesis has extensive unit-testing. Use the Go tool to run tests:
```bash
[...]/lachesis$ make test
```

If everything goes well, it should output something along these lines:
```
?   	github.com/Fantom-foundation/go-lachesis/cmd/dummy	[no test files]
?   	github.com/Fantom-foundation/go-lachesis/cmd/dummy/commands	[no test files]
?   	github.com/Fantom-foundation/go-lachesis/cmd/dummy_client	[no test files]
?   	github.com/Fantom-foundation/go-lachesis/cmd/lachesis	[no test files]
?   	github.com/Fantom-foundation/go-lachesis/cmd/lachesis/commands	[no test files]
?   	github.com/Fantom-foundation/go-lachesis/tester	[no test files]
ok  	github.com/Fantom-foundation/go-lachesis/src/common	(cached)
ok  	github.com/Fantom-foundation/go-lachesis/src/crypto	(cached)
ok  	github.com/Fantom-foundation/go-lachesis/src/difftool	(cached)
ok  	github.com/Fantom-foundation/go-lachesis/src/dummy	0.522s
?   	github.com/Fantom-foundation/go-lachesis/src/lachesis	[no test files]
?   	github.com/Fantom-foundation/go-lachesis/src/log	[no test files]
?   	github.com/Fantom-foundation/go-lachesis/src/mobile	[no test files]
ok  	github.com/Fantom-foundation/go-lachesis/src/net	(cached)
ok  	github.com/Fantom-foundation/go-lachesis/src/node	9.832s
?   	github.com/Fantom-foundation/go-lachesis/src/pb	[no test files]
ok  	github.com/Fantom-foundation/go-lachesis/src/peers	(cached)
ok  	github.com/Fantom-foundation/go-lachesis/src/poset	9.627s
ok  	github.com/Fantom-foundation/go-lachesis/src/proxy	1.019s
?   	github.com/Fantom-foundation/go-lachesis/src/proxy/internal	[no test files]
?   	github.com/Fantom-foundation/go-lachesis/src/proxy/proto	[no test files]
?   	github.com/Fantom-foundation/go-lachesis/src/service	[no test files]
?   	github.com/Fantom-foundation/go-lachesis/src/utils	[no test files]
?   	github.com/Fantom-foundation/go-lachesis/src/version	[no test files]
```

## Cross-build from source

The easiest way to build binaries is to do so in a hermetic Docker container.
Use this simple command:

```bash
[...]/lachesis$ make dist
```
This will launch the build in a Docker container and write all the artifacts in
the build/ folder.

```bash
[...]/lachesis$ tree --charset=nwildner build
build
|-- dist
|   |-- lachesis_0.4.3_SHA256SUMS
|   |-- lachesis_0.4.3_darwin_386.zip
|   |-- lachesis_0.4.3_darwin_amd64.zip
|   |-- lachesis_0.4.3_freebsd_386.zip
|   |-- lachesis_0.4.3_freebsd_arm.zip
|   |-- lachesis_0.4.3_linux_386.zip
|   |-- lachesis_0.4.3_linux_amd64.zip
|   |-- lachesis_0.4.3_linux_arm.zip
|   |-- lachesis_0.4.3_windows_386.zip
|   `-- lachesis_0.4.3_windows_amd64.zip
|-- lachesis
`-- pkg
    |-- darwin_386
    |   `-- lachesis
    |-- darwin_386.zip
    |-- darwin_amd64
    |   `-- lachesis
    |-- darwin_amd64.zip
    |-- freebsd_386
    |   `-- lachesis
    |-- freebsd_386.zip
    |-- freebsd_arm
    |   `-- lachesis
    |-- freebsd_arm.zip
    |-- linux_386
    |   `-- lachesis
    |-- linux_386.zip
    |-- linux_amd64
    |   `-- lachesis
    |-- linux_amd64.zip
    |-- linux_arm
    |   `-- lachesis
    |-- linux_arm.zip
    |-- windows_386
    |   `-- lachesis.exe
    |-- windows_386.zip
    |-- windows_amd64
    |   `-- lachesis.exe
    `-- windows_amd64.zip

11 directories, 29 files
```
