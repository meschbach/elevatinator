#!/bin/bash

./update-proto.sh
set -xe
go test ./...

go build .
go build -o queue ./cmd/queue
go build -o scenarios ./cmd/scenarios

for arch in arm64 amd64
do
  for os in linux darwin
  do
    echo "Building ${os}/${arch}"
    CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build -ldflags='-w -s -extldflags "-static"' -o "cmd/webservice/${arch}_${os}" ./cmd/webservice
  done
done
