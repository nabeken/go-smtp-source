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

## Performance

go-smtp-source should be performant.
I measured the performance for go-smtp-source and smtp-source against smtp-sink with sending 10000 messages.

```
docker compose build
docker compose up -d sink
docker compose run bench
Start sending 10000 messages... (GOMAXPROCS=1)

Concurrency: 1
smtp-source:
1.32user 8.75system 0:22.78elapsed 44%CPU (0avgtext+0avgdata 3836maxresident)k
0inputs+0outputs (0major+230minor)pagefaults 0swaps

go-smtp-source:
1.34user 5.99system 0:11.59elapsed 63%CPU (0avgtext+0avgdata 10248maxresident)k
0inputs+0outputs (0major+1699minor)pagefaults 0swaps
-------------------------
Concurrency: 100
smtp-source:
0.21user 2.52system 0:02.76elapsed 98%CPU (0avgtext+0avgdata 4484maxresident)k
0inputs+0outputs (0major+451minor)pagefaults 0swaps

go-smtp-source:
0.68user 2.08system 0:03.03elapsed 90%CPU (0avgtext+0avgdata 10252maxresident)k
0inputs+0outputs (0major+1705minor)pagefaults 0swaps
-------------------------
Concurrency: 1000
smtp-source:
0.19user 2.61system 0:03.80elapsed 73%CPU (0avgtext+0avgdata 4344maxresident)k
0inputs+0outputs (0major+461minor)pagefaults 0swaps

go-smtp-source:
0.79user 1.79system 0:02.98elapsed 86%CPU (0avgtext+0avgdata 12236maxresident)k
0inputs+0outputs (0major+2234minor)pagefaults 0swaps
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
