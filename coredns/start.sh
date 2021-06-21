#!/bin/bash

docker run -d \
    --network host \
    --name coredns \
    --restart always \
    -w /config \
    --log-opt max-size=10m \
    --log-opt max-file=3 \
    -v $PWD:/config \
    coredns/coredns
