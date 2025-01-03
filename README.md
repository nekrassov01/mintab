mintab
======

[![CI](https://github.com/nekrassov01/mintab/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/nekrassov01/mintab/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/nekrassov01/mintab/graph/badge.svg?token=RIV62CQILM)](https://codecov.io/gh/nekrassov01/mintab)
[![Go Reference](https://pkg.go.dev/badge/github.com/nekrassov01/mintab.svg)](https://pkg.go.dev/github.com/nekrassov01/mintab)
[![Go Report Card](https://goreportcard.com/badge/github.com/nekrassov01/mintab)](https://goreportcard.com/report/github.com/nekrassov01/mintab)

mintab is a minimum ASCII table utilities for golang

Motivation
----------

While [tablewriter](https://github.com/olekukonko/tablewriter) is useful, I wanted a smaller package with features such as backlog format support that tablewriter does not have.

Format
------

Text

```text
+----------+------+-----------+----------+------+-------+-------------+---------------+
| Instance | SG   | Direction | Protocol | From | To    | AddressType | CidrBlock     |
+----------+------+-----------+----------+------+-------+-------------+---------------+
| i-1      | sg-1 | Ingress   | tcp      |   22 |    22 | SG          | sg-10         |
+----------+------+-----------+----------+------+-------+-------------+---------------+
| i-1      | sg-1 | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
+----------+------+-----------+----------+------+-------+-------------+---------------+
| i-1      | sg-2 | Ingress   | tcp      |  443 |   443 | Ipv4        | 0.0.0.0/0     |
+----------+------+-----------+----------+------+-------+-------------+---------------+
| i-1      | sg-2 | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
+----------+------+-----------+----------+------+-------+-------------+---------------+
| i-2      | sg-3 | Ingress   | icmp     |   -1 |    -1 | SG          | sg-11         |
+----------+------+-----------+----------+------+-------+-------------+---------------+
| i-2      | sg-3 | Ingress   | tcp      | 3389 |  3389 | Ipv4        | 10.1.0.0/16   |
+----------+------+-----------+----------+------+-------+-------------+---------------+
| i-2      | sg-3 | Ingress   | tcp      |    0 | 65535 | PrefixList  | pl-id/pl-name |
+----------+------+-----------+----------+------+-------+-------------+---------------+
| i-2      | sg-3 | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
+----------+------+-----------+----------+------+-------+-------------+---------------+
```

Text merged

```text
+----------+------+-----------+----------+------+-------+-------------+---------------+
| Instance | SG   | Direction | Protocol | From | To    | AddressType | CidrBlock     |
+----------+------+-----------+----------+------+-------+-------------+---------------+
| i-1      | sg-1 | Ingress   | tcp      |   22 |    22 | SG          | sg-10         |
+          +      +-----------+----------+------+-------+-------------+---------------+
|          |      | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
+          +------+-----------+----------+------+-------+-------------+---------------+
|          | sg-2 | Ingress   | tcp      |  443 |   443 | Ipv4        | 0.0.0.0/0     |
+          +      +-----------+----------+------+-------+-------------+---------------+
|          |      | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
+----------+------+-----------+----------+------+-------+-------------+---------------+
| i-2      | sg-3 | Ingress   | icmp     |   -1 |    -1 | SG          | sg-11         |
+          +      +           +----------+------+-------+-------------+---------------+
|          |      |           | tcp      | 3389 |  3389 | Ipv4        | 10.1.0.0/16   |
+          +      +           +          +------+-------+-------------+---------------+
|          |      |           |          |    0 | 65535 | PrefixList  | pl-id/pl-name |
+          +      +-----------+----------+------+-------+-------------+---------------+
|          |      | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
+----------+------+-----------+----------+------+-------+-------------+---------------+
```

Compressed-text merged

```text
+----------+------+-----------+----------+------+-------+-------------+---------------+
| Instance | SG   | Direction | Protocol | From | To    | AddressType | CidrBlock     |
+----------+------+-----------+----------+------+-------+-------------+---------------+
| i-1      | sg-1 | Ingress   | tcp      |   22 |    22 | SG          | sg-10         |
|          |      | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
|          | sg-2 | Ingress   | tcp      |  443 |   443 | Ipv4        | 0.0.0.0/0     |
|          |      | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
+----------+------+-----------+----------+------+-------+-------------+---------------+
| i-2      | sg-3 | Ingress   | icmp     |   -1 |    -1 | SG          | sg-11         |
|          |      |           | tcp      | 3389 |  3389 | Ipv4        | 10.1.0.0/16   |
|          |      |           |          |    0 | 65535 | PrefixList  | pl-id/pl-name |
|          |      | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
+----------+------+-----------+----------+------+-------+-------------+---------------+
```

Markdown merged

```text
| Instance | SG   | Direction | Protocol | From | To    | AddressType | CidrBlock     |
| -------- | ---- | --------- | -------- | ---- | ----- | ----------- | ------------- |
| i-1      | sg-1 | Ingress   | tcp      | 22   | 22    | SG          | sg-10         |
|          |      | Egress    | -1       | 0    | 0     | Ipv4        | 0.0.0.0/0     |
|          | sg-2 | Ingress   | tcp      | 443  | 443   | Ipv4        | 0.0.0.0/0     |
|          |      | Egress    | -1       | 0    | 0     | Ipv4        | 0.0.0.0/0     |
| i-2      | sg-3 | Ingress   | icmp     | -1   | -1    | SG          | sg-11         |
|          |      |           | tcp      | 3389 | 3389  | Ipv4        | 10.1.0.0/16   |
|          |      |           |          | 0    | 65535 | PrefixList  | pl-id/pl-name |
|          |      | Egress    | -1       | 0    | 0     | Ipv4        | 0.0.0.0/0     |
```

Backlog merged

```text
| Instance | SG   | Direction | Protocol | From | To    | AddressType | CidrBlock     |h
| i-1      | sg-1 | Ingress   | tcp      |   22 |    22 | SG          | sg-10         |
|          |      | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
|          | sg-2 | Ingress   | tcp      |  443 |   443 | Ipv4        | 0.0.0.0/0     |
|          |      | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
| i-2      | sg-3 | Ingress   | icmp     |   -1 |    -1 | SG          | sg-11         |
|          |      |           | tcp      | 3389 |  3389 | Ipv4        | 10.1.0.0/16   |
|          |      |           |          |    0 | 65535 | PrefixList  | pl-id/pl-name |
|          |      | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
```

Support
-------

- Support markdown table format
- **Support [backlog](https://support-ja.backlog.com/hc/ja/articles/360035641594-%E3%83%86%E3%82%AD%E3%82%B9%E3%83%88%E6%95%B4%E5%BD%A2%E3%81%AE%E3%83%AB%E3%83%BC%E3%83%AB-Backlog%E8%A8%98%E6%B3%95#%E8%A1%A8) table format**
- Support multiple lines in a row
- **Support direct loading of struct slices**
- Support for column merging based on previous field values
- Support for column exclusion
- Support for HTML special character escapes (designed primarily for markdown)
- Support for string concatenation when the field is a slice of the primitive type values
- Support automatic string conversion of byte slices

Benchmark
---------

mintab is memory-efficient.

```text
$ make bench
go test -run=^$ -bench=. -benchmem -count 5 -cpuprofile=cpu.prof -memprofile=mem.prof
goos: darwin
goarch: arm64
pkg: benchmarks
cpu: Apple M2
BenchmarkMintabInput-8             38288             30660 ns/op            2488 B/op         41 allocs/op
BenchmarkMintabInput-8             39964             30058 ns/op            2488 B/op         41 allocs/op
BenchmarkMintabInput-8             38894             29954 ns/op            2488 B/op         41 allocs/op
BenchmarkMintabInput-8             39504             30917 ns/op            2488 B/op         41 allocs/op
BenchmarkMintabInput-8             39993             29877 ns/op            2488 B/op         41 allocs/op
BenchmarkMintabStruct-8            37072             31838 ns/op            2920 B/op         80 allocs/op
BenchmarkMintabStruct-8            38244             32118 ns/op            2920 B/op         80 allocs/op
BenchmarkMintabStruct-8            38517             31612 ns/op            2920 B/op         80 allocs/op
BenchmarkMintabStruct-8            37584             31750 ns/op            2920 B/op         80 allocs/op
BenchmarkMintabStruct-8            38352             31989 ns/op            2920 B/op         80 allocs/op
BenchmarkTableWriter-8             12877             95667 ns/op           11454 B/op        639 allocs/op
BenchmarkTableWriter-8             12712             95992 ns/op           11457 B/op        639 allocs/op
BenchmarkTableWriter-8             12834             95187 ns/op           11443 B/op        639 allocs/op
BenchmarkTableWriter-8             12716             95045 ns/op           11420 B/op        639 allocs/op
BenchmarkTableWriter-8             12680             94409 ns/op           11449 B/op        639 allocs/op
BenchmarkGoPrettyTable-8           65367             16973 ns/op            6540 B/op        192 allocs/op
BenchmarkGoPrettyTable-8           68565             17483 ns/op            6540 B/op        192 allocs/op
BenchmarkGoPrettyTable-8           70461             17172 ns/op            6540 B/op        192 allocs/op
BenchmarkGoPrettyTable-8           65881             17300 ns/op            6540 B/op        192 allocs/op
BenchmarkGoPrettyTable-8           70099             17335 ns/op            6540 B/op        192 allocs/op
PASS
ok      benchmarks      33.220s
```

Notes
-----

- Nested structs are not supported
- Using reflect

Usage
-----

[Example](examples/example_test.go)

Todo
----

- [ ] Add pre-loading support for streaming processing
- [ ] Add paging for large inputs
- [ ] Add minimal styling
- [ ] Add caption
- [ ] Add escape sequence support
- [ ] Add word wrapping with new line
- [ ] Improve performance and reduce memory allocations

Author
------

[nekrassov01](https://github.com/nekrassov01)

License
-------

[MIT](https://github.com/nekrassov01/mintab/blob/main/LICENSE)
