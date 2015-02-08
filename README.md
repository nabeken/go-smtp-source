# go-smtp-source

[![wercker status](https://app.wercker.com/status/0d12a3d5376d3b6488247867269f2302/m "wercker status")](https://app.wercker.com/project/bykey/0d12a3d5376d3b6488247867269f2302)

go-smtp-source is a simple drop replacement for smtp-source in Postfix distribution written in Go

## Motivation

I want to add some feature to smtp-source. I don't want to go with C because we have [Go](http://golang.org).

## Features (including TODOs)

go-smtp-source does not providea all features that smtp-source provided but it has some additional feature.

- STARTTLS support
- Precious time metrics support (smtp-source does not provide elasped time. We need to use `time` with smtp-source.)
- Clustering support for distributed load testing

See [smtp-source(1)](http://www.postfix.org/smtp-source.1.html) about original smtp-source.

## Performance

go-smtp-source should be performant.
I measured the performance for go-smtp-source and smtp-source against smtp-sink with sending 10000 messages.

```
$ ./bench.sh
Start sending 10000 messages...

Concurrency: 1
smtp-source:
6.99user 13.81system 0:30.88elapsed 67%CPU (0avgtext+0avgdata 2696maxresident)k
0inputs+0outputs (0major+128minor)pagefaults 0swaps

go-smtp-source:
4.34user 15.27system 0:18.00elapsed 108%CPU (0avgtext+0avgdata 6980maxresident)k
0inputs+0outputs (0major+335minor)pagefaults 0swaps
-------------------------
Concurrency: 100
smtp-source:
0.50user 2.07system 0:02.58elapsed 99%CPU (0avgtext+0avgdata 3164maxresident)k
0inputs+0outputs (0major+304minor)pagefaults 0swaps

go-smtp-source:
0.52user 5.51system 0:05.57elapsed 108%CPU (0avgtext+0avgdata 9820maxresident)k
0inputs+0outputs (0major+1036minor)pagefaults 0swaps
-------------------------
Concurrency: 1000
smtp-source:
0.60user 2.31system 0:02.93elapsed 99%CPU (0avgtext+0avgdata 3028maxresident)k
0inputs+0outputs (0major+216minor)pagefaults 0swaps

go-smtp-source:
0.58user 4.30system 0:04.45elapsed 109%CPU (0avgtext+0avgdata 31828maxresident)k
0inputs+0outputs (0major+7125minor)pagefaults 0swaps
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
