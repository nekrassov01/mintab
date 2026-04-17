package mintab

import (
	"io"
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
)

// Render renders the table to the writer.
func (t *Table) Render() {
	if t.numRows == 0 {
		return
	}
	t.printHeader()
	t.printData()
}

func (t *Table) printHeader() {
	if !t.hasHeader || t.numColumns == 0 {
		return
	}
	b := bufPool.Get().(*strings.Builder)
	b.Reset()
	switch t.format {
	case TextFormat, CompressedTextFormat:
		b.Grow(t.tableWidth * 2)
		b.WriteString(t.border)
	case MarkdownFormat:
		b.Grow(t.tableWidth)
	case BacklogFormat:
		b.Grow(t.tableWidth + 1)
	}
	b.WriteString("|")
	for i, h := range t.header {
		t.writeField(b, h, t.colWidths[i])
		b.WriteString("|")
	}
	if t.format == BacklogFormat {
		b.WriteString("h")
	}
	b.WriteString("\n")
	s := b.String()
	b.Reset()
	bufPool.Put(b)
	t.print(s)
}

func (t *Table) printData() {
	if t.format == TextFormat || t.format == CompressedTextFormat {
		t.printBorder()
	}
	if t.format == MarkdownFormat {
		if t.hasHeader || t.numColumns > 0 {
			t.printBorder()
		}
	}
	for i, r := range t.data {
		b := bufPool.Get().(*strings.Builder)
		b.Reset()
		if i > 0 {
			switch t.format {
			case TextFormat:
				b.Grow(t.tableWidth * 2)
				t.writeDataBorder(b, r)
			case CompressedTextFormat:
				if r[0][0] == "" || len(t.mergedFields) == 0 {
					b.Grow(t.tableWidth)
				} else {
					b.Grow(t.tableWidth * 2)
					b.WriteString(t.border)
				}
			case MarkdownFormat, BacklogFormat:
				b.Grow(t.tableWidth)
			}
		} else {
			b.Grow(t.tableWidth)
		}
		t.writeRow(b, i)
		s := b.String()
		b.Reset()
		bufPool.Put(b)
		t.print(s)
	}
	if t.format == TextFormat || t.format == CompressedTextFormat {
		t.printBorder()
	}
}

func (t *Table) printBorder() {
	b := bufPool.Get().(*strings.Builder)
	b.Reset()
	b.Grow(t.tableWidth)
	b.WriteString(t.border)
	s := b.String()
	b.Reset()
	bufPool.Put(b)
	t.print(s)
}

func (t *Table) print(s string) {
	_, _ = io.WriteString(t.w, s)
}

func (t *Table) writeRow(b *strings.Builder, i int) {
	for j := 0; j < t.lineHeights[i]; j++ {
		b.WriteString("|")
		for k, elems := range t.data[i] {
			if j < len(elems) {
				t.writeField(b, elems[j], t.colWidths[k])
			} else {
				t.writeField(b, "", t.colWidths[k])
			}
			b.WriteString("|")
		}
		b.WriteString("\n")
	}
}

func (t *Table) writeDataBorder(b *strings.Builder, row [][]string) {
	sep := "+"
	for i, field := range row {
		b.WriteString(sep)
		v := " "
		if field[0] != "" {
			v = "-"
		}
		for j := 0; j < t.colWidths[i]+t.marginWidthBothSides; j++ {
			b.WriteString(v)
		}
	}
	b.WriteString(sep)
	b.WriteString("\n")
}

func (t *Table) writeField(b *strings.Builder, s string, w int) {
	b.WriteString(t.margin)
	isN := isNum(s)
	if !isN {
		b.WriteString(s)
	}
	pad := w - runewidth.StringWidth(s)
	if pad > 0 {
		for range pad {
			b.WriteByte(' ')
		}
	}
	if isN {
		b.WriteString(s)
	}
	b.WriteString(t.margin)
}

func isNum(s string) bool {
	if len(s) == 0 {
		return false
	}
	start := 0
	if s[0] == '-' {
		start = 1
		if len(s) == 1 {
			return false
		}
	}
	n := 0
	d := false
	for i := start; i < len(s); i++ {
		if s[i] == '.' {
			n++
			if n > 1 {
				return false
			}
			continue
		}
		if !unicode.IsDigit(rune(s[i])) {
			return false
		}
		d = true
	}
	return d
}
