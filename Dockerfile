FROM ubuntu:latest

ENV DEBIAN_FRONTEND noninteractive

RUN echo "root:root" | chpasswd

RUN apt update --fix-missing
RUN apt upgrade -y
RUN apt install -y wget git supervisor
RUN apt clean

COPY docker/supervisord-api.conf /etc/supervisor/conf.d/supervisord-api.conf

RUN mkdir -p /api
WORKDIR /api

RUN mkdir -p data/database
RUN mkdir -p data/images

COPY release/bus-stats-linux-amd64 bus-stats
COPY bin/surreal-v1.0.0-beta.8.linux-amd64 surreal
COPY config.toml /api/config.toml

RUN chmod +x ./bus-stats
RUN chmod +x ./surreal

COPY docker/startup.sh /startup.sh

EXPOSE 8080

CMD ["/api/bus-stats", "start"]