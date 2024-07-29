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

- Support markdown table format
- **Support [backlog](https://support-ja.backlog.com/hc/ja/articles/360035641594-%E3%83%86%E3%82%AD%E3%82%B9%E3%83%88%E6%95%B4%E5%BD%A2%E3%81%AE%E3%83%AB%E3%83%BC%E3%83%AB-Backlog%E8%A8%98%E6%B3%95#%E8%A1%A8) table format**
- Support multiple lines in a row
- **Support direct loading of struct slices**
- Support for column merging based on previous field values
- Support for column exclusion
- Support for HTML special character escapes (designed primarily for markdown)
- Support for string concatenation when the field is a slice of the primitive type values
- Support automatic string conversion of byte slices

Notes
-----

- Nested structs are not supported
- Using reflect

Usage
-----

[Example](example_test.go)

Todo
----

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
