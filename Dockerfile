FROM ubuntu:latest

ENV DEBIAN_FRONTEND noninteractive

RUN apt update
RUN apt upgrade -y --fix-missing 
RUN apt install -y curl
RUN apt install -y git

RUN curl -SsfL https://go.dev/dl/go1.19.3.linux-amd64.tar.gz --output go1.19.3.linux-amd64.tar.gz
RUN rm -rf /usr/local/go
RUN tar -C /usr/local -xzf go1.19.3.linux-amd64.tar.gz
RUN export PATH=$PATH:/usr/local/go/bin

RUN git clone https://github.com/buonotti/bus-stats-api
WORKDIR /bus-stats-api

RUN curl -SsfL https://github.com/surrealdb/surrealdb/releases/download/v1.0.0-beta.8/surreal-v1.0.0-beta.8.linux-amd64.tgz --output surrealdb.tgz
RUN mkdir tmp
RUN tar -xf surrealdb.tgz -C tmp/
RUN mv tmp/surreal surrealdb

RUN mkdir -p data/database
RUN mkdir -p data/images

RUN /usr/local/go/bin/go build -o bus-stats-api

COPY config.toml config.toml

RUN chmod +x ./bus-stats-api
RUN chmod +x ./surrealdb

EXPOSE 8080

CMD ["/bus-stats-api/bus-stats-api", "start"]