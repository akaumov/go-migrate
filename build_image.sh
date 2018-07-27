#!/usr/bin/env bash

IMAGE_NAME=azatk/go-migrate:latest
docker build -t "$IMAGE_NAME" .
docker push "$IMAGE_NAME"