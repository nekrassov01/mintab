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

Notes
-----

- Only non-nested struct slices are accepted
- Using reflect

Usage
-----

[example](example_test.go)

Author
------

[nekrassov01](https://github.com/nekrassov01)

License
-------

[MIT](https://github.com/nekrassov01/mintab/blob/main/LICENSE)
