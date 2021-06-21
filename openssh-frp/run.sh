#!/bin/bash
# 
# sh -p 6001 app@139.198.15.135
#
docker container run --restart=always \
    -d -v ~/logs:/logs \
    --name sshd \
    --hostname logserver \
    -p 2222:22 \
    -e APP_PASSWORD="manunkind"  \
    -e SSH_GATEWAY_PORTS=yes \
    openssh:tmux-3.0
