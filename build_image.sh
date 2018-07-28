#!/usr/bin/env bash
set -e
./build_app.sh

IMAGE_NAME=azatk/go-migrate:latest
docker build -t "$IMAGE_NAME" .
docker push "$IMAGE_NAME"