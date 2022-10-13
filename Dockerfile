FROM ubuntu:latest

ENV DEBIAN_FRONTEND noninteractive

RUN mkdir -p /api
WORKDIR /api

RUN mkdir -p data/database
RUN mkdir -p data/images
RUN mkdir -p bin/

COPY release/bus-stats-linux-amd64 bus-stats
COPY bin/surreal-v1.0.0-beta.8.linux-amd64 bin/surreal-v1.0.0-beta.8.linux-amd64
COPY config.toml /api/config.toml

RUN chmod +x ./bus-stats
RUN chmod +x ./bin/surreal-v1.0.0-beta.8.linux-amd64

EXPOSE 8080

CMD ["/api/bus-stats", "start"]