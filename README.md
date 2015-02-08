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
I measured the performance for go-smtp-source and smtp-source against smtp-sink.

(TBD)

## Installation

Download from [releases](https://github.com/nabeken/docker-cleaner/releases).

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
