FROM golang:1.25

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
      postfix \
      time \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/go-smtp-source

COPY bench.sh /root/bench.sh
COPY bench_check.sh /root/bench_check.sh
