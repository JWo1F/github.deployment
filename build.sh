#!/bin/bash

function build {
  echo "Building for OS: $1 ARCH: $2"
  GOOS=$1 GOARCH=$2 go build -ldflags "-s -w" -o build/deployment-$1-$2 src/*.go
}

rm -rf build/*
build linux 386
build linux amd64
build linux arm
build linux arm64
build darwin arm64
build darwin amd64