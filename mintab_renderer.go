package mintab

import (
	"io"
	"unicode"

	"github.com/mattn/go-runewidth"
)

// Render renders the table to the writer.
func (t *Table) Render() {
	if len(t.data) == 0 {
		return
	}
	if t.hasHeader && len(t.header) > 0 {
		switch t.format {
		case TextFormat, CompressedTextFormat:
			t.printBorder()
		}
		t.printHeader()
	}
	if t.format != BacklogFormat {
		t.printBorder()
	}
	t.printData()
	switch t.format {
	case TextFormat, CompressedTextFormat:
		t.printBorder()
	}
}

func (t *Table) printHeader() {
	t.b.Reset()
	t.b.Grow(t.tableWidth)
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
	for i, row := range t.data {
		if i > 0 {
			if t.format == TextFormat {
				t.printDataBorder(row)
			}
			if t.format == CompressedTextFormat && row[0] != "" {
				t.printBorder()
			}
		}
		for j := 0; j < t.lineHeights[i]; j++ {
			t.b.Reset()
			t.b.Grow(t.tableWidth * t.lineHeights[i])
			t.b.WriteString("|")
			for k, elems := range t.multilineData[i] {
				if j < len(elems) {
					t.writeField(elems[j], t.colWidths[k])
				} else {
					t.writeField("", t.colWidths[k])
				}
				t.b.WriteString("|")
			}
			t.b.WriteString("\n")
			t.print(t.b.String())
		}
	}
}

func (t *Table) printDataBorder(row []string) {
	t.b.Reset()
	t.b.Grow(t.tableWidth)
	sep := "+"
	for i, field := range row {
		t.b.WriteString(sep)
		v := " "
		if field != "" {
			v = "-"
		}
		for j := 0; j < t.colWidths[i]+t.marginWidth*2; j++ {
			t.b.WriteString(v)
		}
	}
	t.b.WriteString(sep)
	t.b.WriteString("\n")
	t.print(t.b.String())
}

func (t *Table) printBorder() {
	io.WriteString(t.w, t.border)
}

func (t *Table) print(s string) {
	io.WriteString(t.w, s)
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
