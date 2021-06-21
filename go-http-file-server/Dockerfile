FROM golang:1.16.2 AS build

WORKDIR /src
# https://github.com/mjpclab/go-http-file-server
ADD go-http-file-server .
RUN cd src \
    && GO111MODULE=off GOOS=linux go build -o ghfs -ldflags="-s -w"

FROM busybox:glibc
LABEL description="Golang HTTP File Server" \
    maintainer="Shenle Lu <lushenle@gmail.com>"
COPY --from=build /src/src/ghfs /usr/local/bin/
EXPOSE 80/tcp 443/tcp
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
HEALTHCHECK --interval=1s --timeout=1s \
    CMD nc -z 127.0.0.1 80 || exit 1
ENTRYPOINT ["/entrypoint.sh"]

