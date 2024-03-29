#!/bin/bash

service_address="localhost:8999"

set -xe
./build-all.sh
./queue --address "$service_address" run &
./scenarios --ai-address "$service_address" health-probe
./scenarios --ai-address "$service_address" single-up
./scenarios --ai-address "$service_address" single-down
./scenarios --ai-address "$service_address" multiple-up-and-back
kill %1
