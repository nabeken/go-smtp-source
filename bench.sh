#!/bin/bash
set -eo pipefail

M=10000

HOST=${HOST:-sink}
PORT=${PORT:-10025}

echo "Start sending $M messages... (GOMAXPROCS=${GOMAXPROCS:-default})"
echo

for s in 1 100 1000; do
  echo "Concurrency: $s"
  echo "smtp-source:"
  /usr/bin/time smtp-source -s $s -m $M -f from@example.com -t to@example.com -M smtp.example.com ${HOST}:${PORT}

  echo
  echo "go-smtp-source:"
  /usr/bin/time go-smtp-source -s $s -m $M -resolve-once ${HOST}:${PORT}

  echo
  echo "smtp-source (-d):"
  /usr/bin/time smtp-source -d -s $s -m $M -f from@example.com -t to@example.com -M smtp.example.com ${HOST}:${PORT}

  echo
  echo "go-smtp-source (-d):"
  /usr/bin/time go-smtp-source -d -s $s -m $M -resolve-once ${HOST}:${PORT}

  echo "-------------------------"
done
