# go-smtp-source

[![Build Status](https://travis-ci.org/nabeken/go-smtp-source.svg)](https://travis-ci.org/nabeken/go-smtp-source)

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
$ ./bench.sh
Start sending 10000 messages...

Concurrency: 1
smtp-source:
6.26user 15.14system 0:31.83elapsed 67%CPU (0avgtext+0avgdata 2536maxresident)k
0inputs+0outputs (0major+124minor)pagefaults 0swaps

go-smtp-source:
2.93user 14.69system 0:16.55elapsed 106%CPU (0avgtext+0avgdata 3932maxresident)k
0inputs+0outputs (0major+282minor)pagefaults 0swaps
-------------------------
Concurrency: 100
smtp-source:
0.41user 1.94system 0:02.36elapsed 99%CPU (0avgtext+0avgdata 2844maxresident)k
0inputs+0outputs (0major+181minor)pagefaults 0swaps

go-smtp-source:
0.23user 3.89system 0:03.75elapsed 109%CPU (0avgtext+0avgdata 5292maxresident)k
0inputs+0outputs (0major+631minor)pagefaults 0swaps
-------------------------
Concurrency: 1000
smtp-source:
0.50user 1.93system 0:02.44elapsed 99%CPU (0avgtext+0avgdata 2872maxresident)k
0inputs+0outputs (0major+194minor)pagefaults 0swaps

go-smtp-source:
0.37user 4.01system 0:03.91elapsed 111%CPU (0avgtext+0avgdata 12244maxresident)k
0inputs+0outputs (0major+2401minor)pagefaults 0swaps
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
