FROM golang:1.16.2 AS build

WORKDIR /src
RUN git clone https://github.com/mdlayher/wireguard_exporter \
    && cd wireguard_exporter/cmd/wireguard_exporter/ \
    && go build .

FROM busybox:glibc
LABEL description="prometheus exporter for wireguard." \
    maintainer="Shenle Lu <lushenle@gmail.com>"
COPY --from=build /src/wireguard_exporter/cmd/wireguard_exporter/wireguard_exporter .
HEALTHCHECK --interval=1s --timeout=1s \
    CMD wget -q -O - http://localhost:9586/metrics || exit 1
EXPOSE 9586/tcp
ENTRYPOINT ["./wireguard_exporter"]
CMD ["$@"]
