FROM golang:1.8

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
      postfix \
      time \
    && rm -rf /var/lib/apt/lists/*

RUN go get -u github.com/google/gops
RUN go get -u github.com/kardianos/osext

COPY . $GOPATH/src/github.com/nabeken/go-smtp-source/
RUN go get -d -v github.com/nabeken/go-smtp-source
RUN go install -v github.com/nabeken/go-smtp-source
COPY bench.sh /root/bench.sh
