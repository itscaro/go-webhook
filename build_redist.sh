#!/bin/sh

rm -rf ./dist && mkdir ./dist
docker build -t build . -f Dockerfile.redist
docker create --name build build
docker cp build:/root/webhook ./dist/
docker cp build:/root/webhook.exe ./dist/
docker cp build:/root/webhook-darwin ./dist/
docker cp build:/root/webhook-armhf ./dist/
docker rm -f build
