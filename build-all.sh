#!/bin/bash

./update-proto.sh
go build .
go build ./scenarios/bridge_player
go build ./controllers/queue_grpc
