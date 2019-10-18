FROM golang:1.13-alpine3.10 AS builder

ENV GO111MODULE on
ENV GOPROXY https://goproxy.io

RUN apk upgrade \
    && apk add git \
    && go get github.com/shadowsocks/v2ray-plugin

FROM golang:1.13-alpine3.10 AS builder2

ENV GOPROXY direct

RUN apk upgrade \
    && apk add git \
    && go get github.com/lnnupet/middle-fish

FROM alpine:3.10 AS dist

RUN apk upgrade \
    && apk add tzdata \
    && rm -rf /var/cache/apk/*

COPY --from=builder /go/bin/v2ray-plugin /usr/bin/v2ray-plugin
COPY --from=builder2 /go/bin/middle-fish /usr/bin/shadowsocks
COPY config.json /usr/config/config.json

ENTRYPOINT ["shadowsocks", "-config-file-path","/usr/config/config.json"]