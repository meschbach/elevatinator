#!/bin/bash

./update-proto.sh
set -xe
go test ./...

go build .
go build -o queue ./cmd/queue
go build -o scenarios ./cmd/scenarios
