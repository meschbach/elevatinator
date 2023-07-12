#!/bin/bash

set -xe
./build-all.sh
./queue_grpc &
# Give RPC time to bind
sleep 1
./cmd_scenarios --ai-address localhost:9998 single-up
kill %1
