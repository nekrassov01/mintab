mintab
======

[![CI](https://github.com/nekrassov01/mintab/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/nekrassov01/mintab/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/nekrassov01/mintab/graph/badge.svg?token=RIV62CQILM)](https://codecov.io/gh/nekrassov01/mintab)
[![Go Reference](https://pkg.go.dev/badge/github.com/nekrassov01/mintab.svg)](https://pkg.go.dev/github.com/nekrassov01/mintab)
[![Go Report Card](https://goreportcard.com/badge/github.com/nekrassov01/mintab)](https://goreportcard.com/report/github.com/nekrassov01/mintab)

mintab is a minimum ASCII table utilities written in Go

Motivation
----------

While [tablewriter](https://github.com/olekukonko/tablewriter) is useful, I wanted a smaller package with features such as backlog format support that tablewriter does not have.

Format
------

Text

![text](_assets/text.png)

Text merged

![text_merged](_assets/text_merged.png)

Compressed-text merged

![text_compressed](_assets/text_compressed.png)

Markdown merged

![markdown](_assets/markdown_merged.png)

Backlog merged

![backlog](_assets/backlog_merged.png)

Support
-------

- Text table format
- Markdown table format
- Backlog table format
- Group rows based on previous field value
- Ignore specified columns
- Escape HTML special characters
- Set multiple values to a field as a joined string
- Set byte slices as a string

Notes
-----

- Only non-nested struct slices are accepted
- Using reflect

Usage
-----

[Example](example_test.go)

Benchmark
---------

[A quick benchmark](benchmark.go)

This is only for reference as the functions are different, but for simple drawing, it has better performance than TableWriter.

```text
go test -run=^$ -bench=. -benchmem -count 5
goos: darwin
goarch: arm64
pkg: github.com/nekrassov01/mintab
BenchmarkMintab-8                  45578             25500 ns/op           20527 B/op        399 allocs/op
BenchmarkMintab-8                  46488             25449 ns/op           20527 B/op        399 allocs/op
BenchmarkMintab-8                  44702             26457 ns/op           20528 B/op        399 allocs/op
BenchmarkMintab-8                  42699             28344 ns/op           20527 B/op        399 allocs/op
BenchmarkMintab-8                  45213             31852 ns/op           20527 B/op        399 allocs/op
BenchmarkMintabSimple-8            55597             19234 ns/op           13033 B/op        242 allocs/op
BenchmarkMintabSimple-8            64444             18966 ns/op           13033 B/op        242 allocs/op
BenchmarkMintabSimple-8            53935             21939 ns/op           13034 B/op        242 allocs/op
BenchmarkMintabSimple-8            61573             18596 ns/op           13033 B/op        242 allocs/op
BenchmarkMintabSimple-8            64854             19147 ns/op           13033 B/op        242 allocs/op
BenchmarkTableWriter-8             21787             47804 ns/op           25421 B/op        701 allocs/op
BenchmarkTableWriter-8             26362             45354 ns/op           25365 B/op        701 allocs/op
BenchmarkTableWriter-8             26691             44275 ns/op           25332 B/op        701 allocs/op
BenchmarkTableWriter-8             26622             44199 ns/op           25360 B/op        701 allocs/op
BenchmarkTableWriter-8             27138             44492 ns/op           25297 B/op        701 allocs/op
PASS
ok      github.com/nekrassov01/mintab   24.097s
```

Author
------

[nekrassov01](https://github.com/nekrassov01)

License
-------

[MIT](https://github.com/nekrassov01/mintab/blob/main/LICENSE)
