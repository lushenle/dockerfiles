# stage 1: build src code to binary
FROM golang:1.16-buster as builder

COPY *.go /app/

RUN cd /app && CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o main main.go

# stage 2: use alpine as base image
FROM alpine-glibc:alpine-3.14_glibc-2.33
LABEL maintainer="Shnele Lu <lushenle@gmail.com>" \
    app=hellogo \
    version=v1.0

RUN apk update && \
    apk --no-cache add tzdata ca-certificates && \
    cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*

HEALTHCHECK --interval=15s --timeout=5s \
    CMD curl -fs http://localhost:3000/ || exit 1

COPY --from=builder /app/main .
EXPOSE 3000/tcp
CMD ["./main"] 

