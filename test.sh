#!/bin/bash

./build-all.sh
./queue_grpc &
# Give RPC time to bind
sleep 1
./bridge_player
kill %1
