#!/bin/bash

check() {
  GOT=$(grep 'X-Rcpt-Args:' /sink/dump | wc -l)

  if [ "${GOT}" -ne "${M}" ]; then
    echo "NOT OK: wants '${M}' messages but got '${GOT}' messages" >&2
    exit 1
  else
    echo "OK: got '${GOT}' messages" >&2 
  fi
}

reset_dump() {
  rm -f "/sink/dump"
}

set -eo pipefail

M=${M:-10000}

HOST=${HOST:-sink_check}
PORT=${PORT:-10026}

reset_dump

echo "Start sending $M messages... (GOMAXPROCS=${GOMAXPROCS:-default})"
echo

for s in 1 100 1000; do
  echo "Concurrency: $s"
  echo "smtp-source:"
  /usr/bin/time smtp-source -s $s -m $M -f from@example.com -t to@example.com -M smtp.example.com ${HOST}:${PORT}
  check && reset_dump

  echo
  echo "go-smtp-source:"
  /usr/bin/time go-smtp-source -s $s -m $M -resolve-once ${HOST}:${PORT}
  check && reset_dump

  echo
  echo "smtp-source (-d):"
  /usr/bin/time smtp-source -d -s $s -m $M -f from@example.com -t to@example.com -M smtp.example.com ${HOST}:${PORT}
  check && reset_dump

  echo
  echo "go-smtp-source (-d):"
  /usr/bin/time go-smtp-source -d -s $s -m $M -resolve-once ${HOST}:${PORT}
  check && reset_dump

  echo "-------------------------"
done
