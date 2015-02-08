#!/bin/sh
M=10000

echo "Start sending $M messages..."
echo

for s in 1 100 1000; do
  echo "Concurrency: $s"
  echo "smtp-source:"
  /usr/bin/time smtp-source -s $s -m $M -f from@example.com -t to@example.com -M smtp.example.com ${SINK_PORT_10025_TCP_ADDR}:${SINK_PORT_10025_TCP_PORT}

  echo
  echo "go-smtp-source:"
  /usr/bin/time ./go-smtp-source -s $s -m $M ${SINK_PORT_10025_TCP_ADDR}:${SINK_PORT_10025_TCP_PORT}

  echo "-------------------------"
done
