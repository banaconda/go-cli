#!/bin/bash

protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. \
--go-grpc_opt=paths=source_relative pkg/libnet/networker/networker.proto pkg/libvm/vmer/vmer.proto
go build
