FROM nabeken/docker-smtp-sink:latest

USER root
COPY go-smtp-source-linux-amd64 /usr/local/bin/go-smtp-source
RUN chmod +x /usr/local/bin/go-smtp-source

COPY bench.sh /root/bench.sh
