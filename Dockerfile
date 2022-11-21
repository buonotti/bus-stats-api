FROM ubuntu:latest

ENV DEBIAN_FRONTEND noninteractive

RUN apt update
RUN apt upgrade -y --fix-missing 
RUN apt install -y curl

RUN curl -SsfL https://go.dev/dl/go1.19.3.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.19.3.linux-amd64.tar.gz
RUN export PATH=$PATH:/usr/local/go/bin

RUN mkdir -p /api
WORKDIR /api

RUN mkdir -p data/database
RUN mkdir -p data/images

RUN go build -o bus-stats-api

RUN curl -SsfL https://github.com/surrealdb/surrealdb/releases/download/v1.0.0-beta.8/surreal-v1.0.0-beta.8.linux-amd64.tgz --output surreal.tgz
RUN tar -xf surreal.tgz
COPY config.toml config.toml

RUN chmod +x ./bus-stats
RUN chmod +x ./surreal

EXPOSE 8080

CMD ["/api/bus-stats-api", "start"]