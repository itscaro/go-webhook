#!/bin/sh

#--build-arg http_proxy=$http_proxy --build-arg https_proxy=$https_proxy

docker build -t build . -f Dockerfile.redist && \
 docker create --name build build && \
 docker cp build:/root/webhook . && \
 docker cp build:/root/webhook-darwin . && \
 docker cp build:/root/webhook-armhf . && \
 docker rm -f build
