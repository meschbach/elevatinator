#!/bin/bash

./update-proto.sh
set -xe
go test ./...

go build .
go build ./scenarios/bridge_player
go build ./controllers/queue_grpc
go build -o cmd_scenarios ./cmd/scenarios
