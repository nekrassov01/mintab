package mintab

import (
	"io"
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
	t.b.Reset()
	switch t.format {
	case TextFormat, CompressedTextFormat:
		t.b.Grow(t.tableWidth * 2)
		t.b.WriteString(t.border)
	case MarkdownFormat:
		t.b.Grow(t.tableWidth)
	case BacklogFormat:
		t.b.Grow(t.tableWidth + 1)
	}
	t.b.WriteString("|")
	for i, h := range t.header {
		t.writeField(h, t.colWidths[i])
		t.b.WriteString("|")
	}
	if t.format == BacklogFormat {
		t.b.WriteString("h")
	}
	t.b.WriteString("\n")
	t.print(t.b.String())
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
		t.b.Reset()
		if i > 0 {
			switch t.format {
			case TextFormat:
				t.b.Grow(t.tableWidth * 2)
				t.writeDataBorder(r)
			case CompressedTextFormat:
				if r[0][0] == "" || len(t.mergedFields) == 0 {
					t.b.Grow(t.tableWidth)
				} else {
					t.b.Grow(t.tableWidth * 2)
					t.b.WriteString(t.border)
				}
			case MarkdownFormat, BacklogFormat:
				t.b.Grow(t.tableWidth)
			}
		} else {
			t.b.Grow(t.tableWidth)
		}
		t.writeRow(i)
		t.print(t.b.String())
	}
	if t.format == TextFormat || t.format == CompressedTextFormat {
		t.printBorder()
	}
}

func (t *Table) printBorder() {
	t.b.Reset()
	t.b.Grow(t.tableWidth)
	t.b.WriteString(t.border)
	t.print(t.b.String())
}

func (t *Table) print(s string) {
	io.WriteString(t.w, s)
}

func (t *Table) writeRow(i int) {
	for j := 0; j < t.lineHeights[i]; j++ {
		t.b.WriteString("|")
		for k, elems := range t.data[i] {
			if j < len(elems) {
				t.writeField(elems[j], t.colWidths[k])
			} else {
				t.writeField("", t.colWidths[k])
			}
			t.b.WriteString("|")
		}
		t.b.WriteString("\n")
	}
}

func (t *Table) writeDataBorder(row [][]string) {
	sep := "+"
	for i, field := range row {
		t.b.WriteString(sep)
		v := " "
		if field[0] != "" {
			v = "-"
		}
		for j := 0; j < t.colWidths[i]+t.marginWidthBothSides; j++ {
			t.b.WriteString(v)
		}
	}
	t.b.WriteString(sep)
	t.b.WriteString("\n")
}

func (t *Table) writeField(s string, w int) {
	t.b.WriteString(t.margin)
	isN := isNum(s)
	if !isN {
		t.b.WriteString(s)
	}
	p := w - runewidth.StringWidth(s)
	if p > 0 {
		for i := 0; i < p; i++ {
			t.b.WriteByte(' ')
		}
	}
	if isN {
		t.b.WriteString(s)
	}
	t.b.WriteString(t.margin)
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
