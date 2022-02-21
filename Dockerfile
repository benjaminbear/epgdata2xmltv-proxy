FROM debian:11-slim

RUN apt-get update && apt-get install -y \
    wget \
    tzdata \
    && rm -rf /var/lib/apt/lists/* \
    mkdir /epg

WORKDIR /epg
COPY epgdata2xmltv-proxy /epg/epgdata2xmltv-proxy
COPY epgdata_includes /epg/epgdata_includes

EXPOSE 8080

CMD ["sh", "-c", "/epg/epgdata2xmltv-proxy"]
