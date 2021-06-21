FROM frolvlad/alpine-glibc:glibc-2.33
LABEL maintainer="Shenle Lu <lushenle@gmail.com>" \
    app="nerdctl" \
    version="0.8.3.m"
RUN apk add --no-cache ca-certificates \
    git \
    curl
COPY bin/ /usr/sbin/
