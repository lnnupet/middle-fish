FROM golang:1.13-alpine3.10 AS builder

ENV GOPROXY https://goproxy.cn,direct

RUN apk upgrade \
    && apk add git \
    && go get github.com/lnnupet/middle-fish \
    && go get github.com/shadowsocks/v2ray-plugin

FROM alpine:3.10 AS dist

RUN apk upgrade \
    && apk add tzdata \
    && rm -rf /var/cache/apk/*

COPY config.json /usr/config/config.json
COPY --from=builder /go/bin/middle-fish /usr/bin/shadowsocks
COPY --from=builder /go/bin/v2ray-plugin /usr/bin/v2ray-plugin

ENTRYPOINT ["shadowsocks", "-config-file-path","/usr/config/config.json"]