package mintab

import (
	"fmt"
	"io"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
)

// Default values for table rendering.
const (
	TextDefaultEmptyFieldPlaceholder     = "-"
	TextDefaultWordDelimiter             = textNewLine
	MarkdownDefaultEmptyFieldPlaceholder = "\\" + TextDefaultEmptyFieldPlaceholder
	MarkdownDefaultWordDelimiter         = markdownNewLine
	BacklogDefaultEmptyFieldPlaceholder  = TextDefaultEmptyFieldPlaceholder
	BacklogDefaultWordDelimiter          = backlogNewLine

	textNewLine     = "\n"
	markdownNewLine = "<br>"
	backlogNewLine  = "&br;"
)

// A Format represents the output format.
type Format int

// Supported output formats.
const (
	TextFormat           Format = iota // Text table format.
	CompressedTextFormat               // Compressed text table format.
	MarkdownFormat                     // Markdown table format.
	BacklogFormat                      // Backlog-specific table format.
)

// Formats are string representations of output format.
var Formats = []string{
	"text",
	"compressed",
	"markdown",
	"backlog",
}

// String returns the string representation of a Format.
func (o Format) String() string {
	if o >= 0 && int(o) < len(Formats) {
		return Formats[o]
	}
	return ""
}

// Table represents a table structure for rendering data.
type Table struct {
	writer                io.Writer    // Destination for table output
	data                  [][]string   // Table data
	multilineData         [][][]string // Table data with strings divided by newlines
	header                []string     // Table header
	format                Format       // Output format
	newLine               string       // New line string: "\n"|"<br>"|"&br;"
	emptyFieldPlaceholder string       // Placeholder for empty fields
	wordDelimiter         string       // Delimiter for words within a field
	lineHeights           []int        // Height of lines consisting of fields containing line breaks
	columnWidths          []int        // Max widths of each columns
	border                string       // Border line based on column widths
	tableWidth            int          // Table full width
	marginWidth           int          // Margin size around the field
	margin                string       // Whitespaces around the field
	hasHeader             bool         // Whether header rendering
	hasEscape             bool         // Whether HTML escaping
	mergedFields          []int        // Indices of fields to be merged
	ignoredFields         []int        // Indices of fields to be ignored
}

// New instantiates a new Table with the writer and options.
func New(w io.Writer, opts ...Option) *Table {
	t := &Table{
		writer:                w,
		format:                TextFormat,
		newLine:               textNewLine,
		emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
		wordDelimiter:         TextDefaultWordDelimiter,
		marginWidth:           1,
		hasHeader:             true,
	}
	for _, opt := range opts {
		opt(t)
	}
	t.margin = strings.Repeat(" ", t.marginWidth)
	return t
}

// A Option sets an option on a Table.
type Option func(*Table)

// WithFormat sets the output format.
func WithFormat(format Format) Option {
	return func(t *Table) {
		t.format = format
	}
}

// WithHeader sets the table header.
func WithHeader(has bool) Option {
	return func(t *Table) {
		t.hasHeader = has
	}
}

// WithMargin sets the margin size around field values.
func WithMargin(width int) Option {
	return func(t *Table) {
		t.marginWidth = width
	}
}

// WithEmptyFieldPlaceholder sets the placeholder for empty fields.
func WithEmptyFieldPlaceholder(emptyFieldPlaceholder string) Option {
	return func(t *Table) {
		t.emptyFieldPlaceholder = emptyFieldPlaceholder
	}
}

// WithWordDelimiter sets the delimiter to split words in a field.
func WithWordDelimiter(wordDelimiter string) Option {
	return func(t *Table) {
		t.wordDelimiter = wordDelimiter
	}
}

// WithMergeFields sets column indices to be merged.
func WithMergeFields(mergeFields []int) Option {
	return func(t *Table) {
		t.mergedFields = mergeFields
	}
}

// WithIgnoreFields sets column indices to be ignored.
func WithIgnoreFields(ignoreFields []int) Option {
	return func(t *Table) {
		t.ignoredFields = ignoreFields
	}
}

// WithEscape enables or disables HTML escaping.
func WithEscape(has bool) Option {
	return func(t *Table) {
		t.hasEscape = has
	}
}

// Load validates the input and converts it into struct Table.
func (t *Table) Load(v any) error {
	if t.marginWidth < 0 {
		return fmt.Errorf("only unsigned integers are allowed in margin")
	}
	if _, ok := v.([]interface{}); ok {
		return fmt.Errorf("elements of slice must not be empty interface")
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Slice {
		if rv.IsZero() {
			return fmt.Errorf("no data found")
		}
		rv = reflect.Append(reflect.MakeSlice(reflect.SliceOf(rv.Type()), 0, 1), rv)
	}
	if rv.Len() == 0 {
		return fmt.Errorf("no data found")
	}
	e := rv.Index(0)
	if e.Kind() == reflect.Ptr {
		e = e.Elem()
	}
	if e.Kind() != reflect.Struct {
		return fmt.Errorf("elements of slice must be struct or pointer to struct")
	}
	t.setAttr()
	t.setHeader(e.Type())
	if err := t.setData(rv); err != nil {
		return err
	}
	t.setBorder()
	return nil
}

// Render renders the table to the writer.
func (t *Table) Render() {
	if t.hasHeader {
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
	var b strings.Builder
	b.Grow(t.tableWidth)
	b.WriteString("|")
	for i, h := range t.header {
		t.writeField(&b, h, t.columnWidths[i])
		b.WriteString("|")
	}
	if t.format == BacklogFormat {
		b.WriteString("h")
	}
	b.WriteString("\n")
	t.print(b.String())
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
			var b strings.Builder
			b.Grow(t.tableWidth)
			b.WriteString("|")
			for k, elems := range t.multilineData[i] {
				if j < len(elems) {
					t.writeField(&b, elems[j], t.columnWidths[k])
				} else {
					t.writeField(&b, "", t.columnWidths[k])
				}
				b.WriteString("|")
			}
			b.WriteString("\n")
			t.print(b.String())
		}
	}
}

func (t *Table) printDataBorder(row []string) {
	var b strings.Builder
	b.Grow(t.tableWidth)
	sep := "+"
	for i, field := range row {
		b.WriteString(sep)
		v := " "
		if field != "" {
			v = "-"
		}
		for j := 0; j < t.columnWidths[i]+t.marginWidth*2; j++ {
			b.WriteString(v)
		}
	}
	b.WriteString(sep)
	b.WriteString("\n")
	t.print(b.String())
}

func (t *Table) print(s string) {
	io.WriteString(t.writer, s)
}

func (t *Table) printBorder() {
	io.WriteString(t.writer, t.border)
}

func (t *Table) writeField(b *strings.Builder, s string, width int) {
	b.WriteString(t.margin)
	isN := isNum(s)
	if !isN {
		b.WriteString(s)
	}
	p := width - runewidth.StringWidth(s)
	if p > 0 {
		for i := 0; i < p; i++ {
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

func (t *Table) setAttr() {
	var p, d string
	switch t.format {
	case MarkdownFormat:
		p = MarkdownDefaultEmptyFieldPlaceholder
		d = MarkdownDefaultWordDelimiter
		t.newLine = markdownNewLine
	case BacklogFormat:
		p = BacklogDefaultEmptyFieldPlaceholder
		d = BacklogDefaultWordDelimiter
		t.newLine = backlogNewLine
	default:
		p = TextDefaultEmptyFieldPlaceholder
		d = TextDefaultWordDelimiter
	}
	if t.emptyFieldPlaceholder == TextDefaultEmptyFieldPlaceholder {
		t.emptyFieldPlaceholder = p
	}
	if t.wordDelimiter == TextDefaultWordDelimiter {
		t.wordDelimiter = d
	}
}

func (t *Table) setHeader(typ reflect.Type) {
	if len(t.header) > 0 {
		return
	}
	t.columnWidths = make([]int, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if slices.Contains(t.ignoredFields, i) || f.PkgPath != "" {
			continue
		}
		t.header = append(t.header, f.Name)
		t.columnWidths = append(t.columnWidths, runewidth.StringWidth(f.Name))
	}
}

func (t *Table) setData(rv reflect.Value) error {
	t.data = make([][]string, rv.Len())
	t.multilineData = make([][][]string, rv.Len())
	t.lineHeights = make([]int, rv.Len())
	prev := make([]string, len(t.header))
	for i := 0; i < rv.Len(); i++ {
		e := rv.Index(i)
		if e.Kind() == reflect.Ptr {
			e = e.Elem()
		}
		row := make([]string, len(t.header))
		multilineRow := make([][]string, len(t.header))
		isMerge := true
		n := 1
		for j, h := range t.header {
			field := e.FieldByName(h)
			if !field.IsValid() {
				return fmt.Errorf("invalid field detected: %s", h)
			}
			v, err := t.formatField(field)
			if err != nil {
				return fmt.Errorf("failed to format field \"%s\": %w", h, err)
			}
			if slices.Contains(t.mergedFields, j) {
				if v != prev[j] {
					isMerge = false
					prev[j] = v
				}
				if isMerge {
					v = ""
				}
			}
			row[j] = v
			elems := strings.Split(v, "\n")
			multilineRow[j] = elems
			for _, elem := range elems {
				width := runewidth.StringWidth(elem)
				if width > t.columnWidths[j] {
					t.columnWidths[j] = width
				}
			}
			if t.format == TextFormat {
				height := len(elems)
				if height > n {
					n = height
				}
			}
		}
		t.data[i] = row
		t.multilineData[i] = multilineRow
		t.lineHeights[i] = n
	}
	return nil
}

func (t *Table) setBorder() {
	var sep string
	switch t.format {
	case MarkdownFormat, BacklogFormat:
		sep = "|"
	default:
		sep = "+"
	}
	var b strings.Builder
	for _, width := range t.columnWidths {
		b.WriteString(sep)
		for i := 0; i < width+t.marginWidth*2; i++ {
			b.WriteByte('-')
		}
	}
	b.WriteString(sep)
	b.WriteString("\n")
	t.border = b.String()
	t.tableWidth = len(t.border)
}

func (t *Table) formatField(rv reflect.Value) (string, error) {
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return t.emptyFieldPlaceholder, nil
		}
		rv = rv.Elem()
	}
	v := stringer(rv)
	if v == "" {
		switch rv.Kind() {
		case reflect.String:
			v = rv.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v = strconv.FormatInt(rv.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v = strconv.FormatUint(rv.Uint(), 10)
		case reflect.Float32:
			v = strconv.FormatFloat(rv.Float(), 'f', -1, 32)
		case reflect.Float64:
			v = strconv.FormatFloat(rv.Float(), 'f', -1, 64)
		case reflect.Slice, reflect.Array:
			switch {
			case rv.Len() == 0:
				v = t.emptyFieldPlaceholder
			case rv.Type().Elem().Kind() == reflect.Uint8:
				v = string(rv.Bytes())
			default:
				var b strings.Builder
				for i := 0; i < rv.Len(); i++ {
					e := rv.Index(i)
					if i != 0 {
						b.WriteString(t.wordDelimiter)
					}
					if e.Kind() == reflect.Ptr {
						if e.IsNil() {
							b.WriteString(t.emptyFieldPlaceholder)
							continue
						}
						e = e.Elem()
					}
					if f := stringer(e); f != "" {
						b.WriteString(f)
						continue
					}
					if e.Kind() == reflect.Slice && e.Type().Elem().Kind() == reflect.Uint8 {
						b.WriteString(string(e.Bytes()))
						continue
					}
					if e.Kind() == reflect.Slice || e.Kind() == reflect.Array || e.Kind() == reflect.Struct {
						return "", fmt.Errorf("cannot represent nested fields")
					}
					switch e.Kind() {
					case reflect.String:
						v = e.String()
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						v = strconv.FormatInt(e.Int(), 10)
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						v = strconv.FormatUint(e.Uint(), 10)
					case reflect.Float32:
						v = strconv.FormatFloat(e.Float(), 'f', -1, 32)
					case reflect.Float64:
						v = strconv.FormatFloat(e.Float(), 'f', -1, 64)
					default:
						v = fmt.Sprint(e.Interface())
					}
					if v == "" {
						v = t.emptyFieldPlaceholder
					}
					b.WriteString(v)
				}
				v = b.String()
			}
		default:
			v = fmt.Sprint(rv.Interface())
		}
	}
	return strings.TrimSuffix(t.sanitize(v), "\n"), nil
}

func stringer(rv reflect.Value) string {
	if rv.CanInterface() {
		if s, ok := rv.Interface().(fmt.Stringer); ok {
			return s.String()
		}
	}
	return ""
}

func (t *Table) sanitize(s string) string {
	if s == "" {
		return t.emptyFieldPlaceholder
	}
	if t.hasEscape {
		s = t.escape(s)
	}
	if t.format == MarkdownFormat && strings.HasPrefix(s, "*") {
		s = "\\" + s
	}
	if t.format == TextFormat {
		return s
	}
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '\n':
			b.WriteString(t.newLine)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (t *Table) escape(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&lsquo;")
		case '&':
			b.WriteString("&amp;")
		case ' ':
			b.WriteString("&nbsp;")
		case '*':
			b.WriteString("&#42;")
		case '\\':
			b.WriteString("&#92;")
		case '_':
			b.WriteString("&#95;")
		case '|':
			b.WriteString("&#124;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
