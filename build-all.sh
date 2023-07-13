#!/bin/bash

./update-proto.sh
set -xe
go test ./...

go build .
go build ./controllers/queue_grpc
go build -o scenarios ./cmd/scenarios
