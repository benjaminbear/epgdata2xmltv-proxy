FROM golang:latest as builder

ENV GO111MODULE=on
ENV GOPATH=/root/go
RUN mkdir -p /root/go/src/epgdate2xmltv-proxy
COPY . /root/go/src/epgdate2xmltv-proxy/

RUN cd /root/go/src/epgdate2xmltv-proxy && go mod download && GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.Version=$DOCKER_TAG" -o /root/go/bin/epgdate2xmltv-proxy

FROM alpine:latest

RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
    && wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.32-r0/glibc-2.32-r0.apk \
    && apk add glibc-2.32-r0.apk \
    && rm glibc-2.32-r0.apk

WORKDIR /root
COPY --from=builder /root/go/bin/epgdate2xmltv-proxy /root/epgdate2xmltv-proxy
COPY epgdata_includes /root/epgdata_includes

EXPOSE 8080
CMD ["sh", "-c", "/root/epgdate2xmltv-proxy"]