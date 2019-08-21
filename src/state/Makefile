export PROTOC?=protoc
export RM?=rm

PROTOS:=$(wildcard *proto)
PBGOS:=$(PROTOS:.proto=.pb.go)
INTERWIRE:=..

.PHONY: build clean proto buildtests

clean:
	$(RM) -rf *.pb.go

proto: $(PROTOS)
	$(PROTOC) --proto_path=".:$(INTERWIRE):$(GOPATH)/src" --go_out=plugins=grpc:. $(PROTOS)

build:

buildtests:

