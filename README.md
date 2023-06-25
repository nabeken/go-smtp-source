# go-smtp-source

[![Go](https://github.com/nabeken/go-smtp-source/actions/workflows/go.yml/badge.svg)](https://github.com/nabeken/go-smtp-source/actions/workflows/go.yml)

go-smtp-source is a simple drop-in replacement for smtp-source in Postfix distribution written in Go

## Motivation

I want to add some feature to smtp-source. I don't want to go with C because we have [Go](http://golang.org).

## Features (including TODOs)

go-smtp-source does not providea all features that smtp-source provided but it has some additional feature.

- :heavy_check_mark: STARTTLS support
- :construction: Precious time metrics support (smtp-source does not provide elasped time. We need to use `time` with smtp-source.)
- :construction: Clustering support for distributed load testing

See [smtp-source(1)](http://www.postfix.org/smtp-source.1.html) about original smtp-source.

## Bench

go-smtp-source should be performant.
I measured the performance for go-smtp-source and smtp-source against smtp-sink with sending 10000 messages.

```sh
docker compose build
docker compose up -d sink
docker compose run --rm bench
Start sending 10000 messages... (GOMAXPROCS=default)

Concurrency: 1
smtp-source:
0.53user 6.06system 0:14.88elapsed 44%CPU (0avgtext+0avgdata 5880maxresident)k
0inputs+0outputs (0major+310minor)pagefaults 0swaps

go-smtp-source:
0.87user 2.51system 0:03.31elapsed 102%CPU (0avgtext+0avgdata 11408maxresident)k
0inputs+8outputs (0major+4612minor)pagefaults 0swaps

smtp-source (-d):
0.33user 2.86system 0:07.96elapsed 40%CPU (0avgtext+0avgdata 5876maxresident)k
0inputs+0outputs (0major+308minor)pagefaults 0swaps

go-smtp-source (-d):
1.13user 4.92system 0:08.46elapsed 71%CPU (0avgtext+0avgdata 10460maxresident)k
0inputs+8outputs (0major+1697minor)pagefaults 0swaps
-------------------------
Concurrency: 100
smtp-source:
0.31user 1.59system 0:02.16elapsed 88%CPU (0avgtext+0avgdata 6644maxresident)k
0inputs+0outputs (0major+1098minor)pagefaults 0swaps

go-smtp-source:
0.85user 2.55system 0:03.35elapsed 101%CPU (0avgtext+0avgdata 11496maxresident)k
0inputs+8outputs (0major+5212minor)pagefaults 0swaps

smtp-source (-d):
0.20user 0.57system 0:01.07elapsed 72%CPU (0avgtext+0avgdata 6532maxresident)k
0inputs+0outputs (0major+604minor)pagefaults 0swaps

go-smtp-source (-d):
0.29user 0.65system 0:00.91elapsed 104%CPU (0avgtext+0avgdata 12080maxresident)k
0inputs+8outputs (0major+2140minor)pagefaults 0swaps
-------------------------
Concurrency: 1000
smtp-source:
0.40user 2.08system 0:44.61elapsed 5%CPU (0avgtext+0avgdata 6324maxresident)k
0inputs+0outputs (0major+548minor)pagefaults 0swaps

go-smtp-source:
0.89user 2.51system 0:03.35elapsed 101%CPU (0avgtext+0avgdata 11448maxresident)k
0inputs+8outputs (0major+5538minor)pagefaults 0swaps

smtp-source (-d):
0.18user 0.71system 0:04.36elapsed 20%CPU (0avgtext+0avgdata 6792maxresident)k
0inputs+0outputs (0major+625minor)pagefaults 0swaps

go-smtp-source (-d):
0.29user 0.64system 0:00.89elapsed 104%CPU (0avgtext+0avgdata 11748maxresident)k
0inputs+8outputs (0major+2175minor)pagefaults 0swaps
-------------------------
```

## Test

To confirm the delivery result, you can run `bench_check.sh` instead.

```sh
docker compose build
docker compose up -d sink
docker compose run --rm bench /root/bench_check.sh
Start sending 10000 messages... (GOMAXPROCS=default)

Concurrency: 1
smtp-source:
0.54user 5.93system 0:15.28elapsed 42%CPU (0avgtext+0avgdata 5852maxresident)k
0inputs+0outputs (0major+304minor)pagefaults 0swaps
OK: got '10000' messages

go-smtp-source:
0.99user 2.94system 0:04.31elapsed 91%CPU (0avgtext+0avgdata 11436maxresident)k
0inputs+8outputs (0major+5532minor)pagefaults 0swaps
OK: got '10000' messages

smtp-source (-d):
0.29user 2.80system 0:08.33elapsed 37%CPU (0avgtext+0avgdata 5804maxresident)k
0inputs+0outputs (0major+309minor)pagefaults 0swaps
OK: got '10000' messages

go-smtp-source (-d):
1.05user 5.16system 0:09.32elapsed 66%CPU (0avgtext+0avgdata 9964maxresident)k
0inputs+8outputs (0major+1565minor)pagefaults 0swaps
OK: got '10000' messages
-------------------------
Concurrency: 100
smtp-source:
0.43user 2.28system 0:04.28elapsed 63%CPU (0avgtext+0avgdata 6196maxresident)k
0inputs+0outputs (0major+433minor)pagefaults 0swaps
OK: got '10000' messages

go-smtp-source:
0.95user 3.02system 0:04.32elapsed 91%CPU (0avgtext+0avgdata 11360maxresident)k
0inputs+8outputs (0major+5624minor)pagefaults 0swaps
OK: got '10000' messages

smtp-source (-d):
0.33user 1.11system 0:03.05elapsed 47%CPU (0avgtext+0avgdata 6544maxresident)k
0inputs+0outputs (0major+602minor)pagefaults 0swaps
OK: got '10000' messages

go-smtp-source (-d):
0.62user 1.64system 0:02.70elapsed 83%CPU (0avgtext+0avgdata 11556maxresident)k
0inputs+8outputs (0major+2047minor)pagefaults 0swaps
OK: got '10000' messages
-------------------------
Concurrency: 1000
smtp-source:
0.51user 2.96system 1:12.16elapsed 4%CPU (0avgtext+0avgdata 6324maxresident)k
0inputs+0outputs (0major+719minor)pagefaults 0swaps
OK: got '10000' messages

go-smtp-source:
0.99user 2.96system 0:04.32elapsed 91%CPU (0avgtext+0avgdata 11424maxresident)k
0inputs+8outputs (0major+5805minor)pagefaults 0swaps
OK: got '10000' messages

smtp-source (-d):
0.30user 1.33system 0:09.50elapsed 17%CPU (0avgtext+0avgdata 7216maxresident)k
0inputs+0outputs (0major+634minor)pagefaults 0swaps
OK: got '10000' messages

go-smtp-source (-d):
0.60user 1.41system 0:02.48elapsed 81%CPU (0avgtext+0avgdata 11480maxresident)k
0inputs+8outputs (0major+2045minor)pagefaults 0swaps
OK: got '10000' messages
-------------------------
```

## Installation

Download from [releases](https://github.com/nabeken/go-smtp-source/releases).

Or

```sh
go get -u github.com/nabeken/go-smtp-source
```

## Usage

Send 100 messages in 10 concurrency to SMTP server running on 127.0.0.1:10025 over TLS.

```sh
go-smtp-source -s 10 -m 100 -tls 127.0.0.1:10025
```

## smtp-sink

[smtp-sink(1)](http://www.postfix.org/smtp-sink.1.html) is a good friend for benchmarking {go-,}smtp-source.
