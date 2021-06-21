#!/bin/sh
#
/usr/local/bin/ghfs -l 80 \
    -r /var/data \
    --global-auth \
    --upload /upload \
    --default-sort T \
    -H '.*' \
    --archive /download \
    "$@"
