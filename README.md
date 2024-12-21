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
go test -run=^$ -bench=. -benchmem -count 5 -cpuprofile=cpu.prof -memprofile=mem.prof
goos: darwin
goarch: arm64
pkg: benchmarks
cpu: Apple M2
BenchmarkMintabInput-8     	   39860	     29566 ns/op	    3848 B/op	      46 allocs/op
BenchmarkMintabInput-8     	   41198	     29765 ns/op	    3848 B/op	      46 allocs/op
BenchmarkMintabInput-8     	   40148	     29447 ns/op	    3848 B/op	      46 allocs/op
BenchmarkMintabInput-8     	   39190	     29306 ns/op	    3848 B/op	      46 allocs/op
BenchmarkMintabInput-8     	   40102	     29769 ns/op	    3848 B/op	      46 allocs/op
BenchmarkMintabStruct-8    	   38478	     31208 ns/op	    4281 B/op	      85 allocs/op
BenchmarkMintabStruct-8    	   37399	     31519 ns/op	    4280 B/op	      85 allocs/op
BenchmarkMintabStruct-8    	   38781	     31882 ns/op	    4280 B/op	      85 allocs/op
BenchmarkMintabStruct-8    	   38584	     33231 ns/op	    4280 B/op	      85 allocs/op
BenchmarkMintabStruct-8    	   37465	     31361 ns/op	    4280 B/op	      85 allocs/op
BenchmarkTableWriter-8     	   16680	     83253 ns/op	    9661 B/op	     474 allocs/op
BenchmarkTableWriter-8     	   15483	     74068 ns/op	    9695 B/op	     474 allocs/op
BenchmarkTableWriter-8     	   16646	     74946 ns/op	    9664 B/op	     474 allocs/op
BenchmarkTableWriter-8     	   16278	     78036 ns/op	    9663 B/op	     474 allocs/op
BenchmarkTableWriter-8     	   16555	     74486 ns/op	    9699 B/op	     474 allocs/op
BenchmarkGoPrettyTable-8   	  115596	     10444 ns/op	    4402 B/op	     134 allocs/op
BenchmarkGoPrettyTable-8   	  108080	     19097 ns/op	    4402 B/op	     134 allocs/op
BenchmarkGoPrettyTable-8   	  119058	     10530 ns/op	    4403 B/op	     134 allocs/op
BenchmarkGoPrettyTable-8   	  113080	     10464 ns/op	    4403 B/op	     134 allocs/op
BenchmarkGoPrettyTable-8   	  118528	     10465 ns/op	    4403 B/op	     134 allocs/op
PASS
ok  	benchmarks	32.812s
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
