FROM golang:1.15 as builder

ENV GO111MODULE=on
ENV GOPATH=/root/go
RUN mkdir -p /root/go/src/epgdata2xmltv-proxy
COPY . /root/go/src/epgdata2xmltv-proxy/

RUN TAG=$(curl --silent "https://api.github.com/repos/benjaminbear/epgdata2xmltv-proxy/releases/latest" | grep -Po '"tag_name": "\K.*?(?=")') \
    && cd /root/go/src/epgdata2xmltv-proxy && go mod download \
    && GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.Version=${TAG}" -o /root/go/bin/epgdata2xmltv-proxy

FROM ubuntu:20.04

RUN apt-get update && apt-get install -y \
    wget \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /epg
COPY --from=builder /root/go/bin/epgdata2xmltv-proxy /epg/epgdata2xmltv-proxy
COPY epgdata_includes /epg/epgdata_includes

EXPOSE 8080

CMD ["sh", "-c", "/epg/epgdata2xmltv-proxy"]
