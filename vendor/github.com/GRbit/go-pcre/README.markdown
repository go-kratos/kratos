go-pcre
===============

This is a Go language package providing Perl-Compatible RegularExpression
support using libpcre or libpcre++.

## documentation

Use [godoc](https://godoc.org/github.com/GRbit/go-pcre).

## installation

1. install libpcre3-dev or libpcre++-dev

2. go get

```bash
sudo apt-get install libpcre3-dev
go get github.com/GRbit/go-pcre/
```

## usage

Go programs that depend on this package should import this package as
follows to allow automatic downloading:

```go
import (
  "github.com/GRbit/go-pcre/"
)
```

## LICENSE

This is a fork of [go-pcre](https://github.com/pantsing/go-pcre).
