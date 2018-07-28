#!/usr/bin/env bash
set -e
mkdir build

APP_PATH=github.com/akaumov/go-migrate
docker run --rm -v "$PWD":/go/src/"$APP_PATH" -w /go/src/"$APP_PATH" golang:1.10.3-alpine go build -x -v -o ./build/app