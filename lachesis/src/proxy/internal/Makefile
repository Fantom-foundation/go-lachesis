export PROTOC?=protoc
export RM?=rm

PROTOS:=$(wildcard *proto)
PBGOS:=$(PROTOS:.proto=.pb.go)
INTERWIRE:=../..
GOEMPTY:=$(GOPATH)/src/github.com/golang/protobuf/ptypes/empty

.PHONY: build clean proto

clean:
	$(RM) -rf *.pb.go

proto: $(PROTOS)
	$(PROTOC) --proto_path=".:$(INTERWIRE):$(GOPATH)/src:$(GOEMPTY)" --go_out=plugins=grpc,import_path=internal:. $(PROTOS)

build:
