#!/bin/bash

set -xe
./build-all.sh
./queue_grpc &
# Give RPC time to bind
sleep 1
./scenarios --ai-address localhost:9998 single-up
./scenarios --ai-address localhost:9998 single-down
kill %1
