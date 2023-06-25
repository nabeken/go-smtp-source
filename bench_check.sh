#!/bin/bash
set -eo pipefail

check() {
  local r=$1
  WANT=$(expr ${M} \* ${r})
  GOT=$(grep 'X-Rcpt-Args:' /sink/dump | wc -l)

  if [ "${GOT}" -ne "${WANT}" ]; then
    echo "NOT OK: wants '${WANT}' messages but got '${GOT}' messages" >&2
    exit 1
  else
    echo "OK: got '${GOT}' messages" >&2 
  fi
}

reset_dump() {
  rm -f "/sink/dump"
}

M=${M:-10000}

HOST=${HOST:-sink_check}
PORT=${PORT:-10026}

reset_dump

echo "Start sending $M messages... (GOMAXPROCS=${GOMAXPROCS:-default})"
echo

for s in 1 100 1000; do
  for r in 1 3; do
    echo "Concurrency: $s / Recipients: $r"
    echo "smtp-source:"
    /usr/bin/time smtp-source -s $s -m $M -r $r -f from@example.com -t to@example.com -M smtp.example.com ${HOST}:${PORT}
    check $r && reset_dump

    echo
    echo "go-smtp-source:"
    /usr/bin/time go-smtp-source -s $s -m $M -r $r -resolve-once ${HOST}:${PORT}
    check $r && reset_dump

    echo
    echo "smtp-source (-d):"
    /usr/bin/time smtp-source -d -s $s -m $M -r $r -f from@example.com -t to@example.com -M smtp.example.com ${HOST}:${PORT}
    check $r && reset_dump

    echo
    echo "go-smtp-source (-d):"
    /usr/bin/time go-smtp-source -d -s $s -m $M -r $r -resolve-once ${HOST}:${PORT}
    check $r && reset_dump

    echo "-------------------------"
  done
done
