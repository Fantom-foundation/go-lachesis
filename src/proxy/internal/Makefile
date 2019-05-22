export GO?=go
export PROTOC?=protoc
export RM?=rm

.PHONY: build clean proto

clean:
	#$(RM) -rf *.pb.go

proto:
	#$(PROTOC) --go_out="plugins=grpc,Mgoogle/protobuf/empty.proto=$(GOPATH)/src/github.com/golang/protobuf/ptypes/empty:." *.proto
	#$(PROTOC) -I. -I"$(GOPATH)/src/github.com/gogo/protobuf/protobuf" \
	#--go_out="plugins=grpc,import_path=internal:." *.proto

build:

