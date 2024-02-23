package mintab

import (
	"fmt"
	"io"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
)

// Default values for table rendering.
const (
	DefaultEmptyFieldPlaceholder         = "-"
	DefaultWordDelimiter                 = "\n"
	MarkdownDefaultEmptyFieldPlaceholder = "\\-"
	MarkdownDefaultWordDelimiter         = "<br>"
	BacklogDefaultEmptyFieldPlaceholder  = "-"
	BacklogDefaultWordDelimiter          = "&br;"
)

// Format defines the output format of the content.
type Format int

// Enumeration of supported output formats.
const (
	FormatText           Format = iota // Plain text format.
	FormatCompressedText               // Compressed plain text format.
	FormatMarkdown                     // Markdown format.
	FormatBacklog                      // Backlog-specific format.
)

// Formats holds the string representations of each format constant.
var Formats = []string{
	"text",
	"compressed",
	"markdown",
	"backlog",
}

// String returns the string representation of a Format.
// If the format is not within the predefined range, an empty string is returned.
func (o Format) String() string {
	if o >= 0 && int(o) < len(Formats) {
		return Formats[o]
	}
	return ""
}

// Table represents a table structure for rendering data in a matrix format.
type Table struct {
	writer                io.Writer       // Destination for table output.
	builder               strings.Builder // Builder for string concatenation.
	data                  [][]string      // Data holds the content of the table.
	header                []string        // Names of each field in the table header.
	format                Format          // Output format of the table.
	border                string          // Pre-computed border based on column widths.
	margin                int             // Margin size around cell content.
	emptyFieldPlaceholder string          // Placeholder for empty fields.
	wordDelimiter         string          // Delimiter for words within a field.
	mergedFields          []int           // Indices of fields to be merged based on content.
	ignoredFields         []int           // Indices of fields to be ignored during rendering.
	columnWidths          []int           // Calculated max width of each column.
	hasHeader             bool            // Indicates if the header should be rendered.
	hasEscape             bool            // Indicates if escaping should be performed.
}

// New instantiates a new Table with the specified writer and options.
func New(w io.Writer, opts ...Option) *Table {
	t := &Table{
		writer:                w,
		format:                FormatText,
		margin:                1,
		emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
		wordDelimiter:         DefaultWordDelimiter,
		hasHeader:             true,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Option defines a type for functional options used to configure a Table.
type Option func(*Table)

// WithFormat sets the output format of the table.
func WithFormat(format Format) Option {
	return func(t *Table) {
		t.format = format
	}
}

// WithHeader configures the rendering of the table header.
func WithHeader(has bool) Option {
	return func(t *Table) {
		t.hasHeader = has
	}
}

// WithMargin sets the margin size around cell content.
func WithMargin(margin int) Option {
	return func(t *Table) {
		t.margin = margin
	}
}

// WithEmptyFieldPlaceholder sets the placeholder for empty fields.
func WithEmptyFieldPlaceholder(emptyFieldPlaceholder string) Option {
	return func(t *Table) {
		t.emptyFieldPlaceholder = emptyFieldPlaceholder
	}
}

// WithWordDelimiter sets the delimiter for splitting words in a field.
func WithWordDelimiter(wordDelimiter string) Option {
	return func(t *Table) {
		t.wordDelimiter = wordDelimiter
	}
}

// WithMergeFields specifies columns for merging based on their content.
func WithMergeFields(mergeFields []int) Option {
	return func(t *Table) {
		t.mergedFields = mergeFields
	}
}

// WithIgnoreFields specifies columns to be ignored during rendering.
func WithIgnoreFields(ignoreFields []int) Option {
	return func(t *Table) {
		t.ignoredFields = ignoreFields
	}
}

// WithEscape enables or disables escaping of field content.
func WithEscape(has bool) Option {
	return func(t *Table) {
		t.hasEscape = has
	}
}

// Load validates the input and converts it into table data.
// Returns an error if the input is not a slice or if it's empty.
func (t *Table) Load(input any) error {
	if t.margin < 0 {
		return fmt.Errorf("only unsigned integers are allowed in margin")
	}
	if _, ok := input.([]interface{}); ok {
		return fmt.Errorf("elements of slice must not be empty interface")
	}
	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("input must be a slice or a pointer to a slice")
	}
	if v.Len() == 0 {
		return fmt.Errorf("no data found")
	}
	e := v.Index(0)
	if e.Kind() == reflect.Ptr {
		e = e.Elem()
	}
	if e.Kind() != reflect.Struct {
		return fmt.Errorf("elements of slice must be struct or pointer to struct")
	}
	t.setAttr()
	t.setHeader(e.Type())
	if err := t.setData(v); err != nil {
		return err
	}
	if t.format != FormatBacklog {
		t.setBorder()
	}
	return nil
}

// Out renders the table to the specified writer.
// It supports markdown and backlog formats for easy copying and pasting.
func (t *Table) Out() {
	if t.hasHeader {
		switch t.format {
		case FormatText, FormatCompressedText:
			t.printBorder()
		}
		t.printHeader()
	}
	if t.format != FormatBacklog {
		t.printBorder()
	}
	t.printData()
	switch t.format {
	case FormatText, FormatCompressedText:
		t.printBorder()
	}
	fmt.Fprintf(t.writer, "\n")
}

// printHeader renders the table header.
func (t *Table) printHeader() {
	t.builder.Reset()
	margin := t.getMargin()
	t.builder.WriteString("|")
	for i, h := range t.header {
		t.builder.WriteString(margin)
		t.builder.WriteString(pad(h, t.columnWidths[i]))
		t.builder.WriteString(margin)
		t.builder.WriteString("|")
	}
	if t.format == FormatBacklog {
		t.builder.WriteString("h")
	}
	fmt.Fprintln(t.writer, t.builder.String())
}

// printData renders the table data with dynamic conditional borders.
func (t *Table) printData() {
	var prev []string
	for ri, row := range t.data {
		lines := 1
		splitedCells := make([][]string, len(row))
		hasBorder := false
		for fi, field := range row {
			splitedCell := strings.Split(field, "\n")
			splitedCells[fi] = splitedCell
			if len(splitedCell) > lines {
				lines = len(splitedCell)
			}
			if ri == 0 {
				continue
			}
			if t.format == FormatCompressedText {
				if field != "" && (len(prev) <= fi || prev[fi] == "") || (row[0] != "") {
					hasBorder = true
				}
			}
			if t.format == FormatText {
				hasBorder = true
			}
		}
		if hasBorder {
			t.printDataBorder(row)
		}
		for line := 0; line < lines; line++ {
			t.builder.Reset()
			t.builder.WriteString("|")
			for sfi, splitedCell := range splitedCells {
				margin := t.getMargin()
				t.builder.WriteString(margin)
				if line < len(splitedCell) {
					t.builder.WriteString(pad(splitedCell[line], t.columnWidths[sfi]))
				} else {
					t.builder.WriteString(pad("", t.columnWidths[sfi]))
				}
				t.builder.WriteString(margin)
				t.builder.WriteString("|")
			}
			fmt.Fprintln(t.writer, t.builder.String())
		}
		prev = row
	}
}

// printDataBorder prints a conditional border based on the emptiness of fields in the current row,
// with continuity in border characters based on the emptiness of adjacent fields.
func (t *Table) printDataBorder(row []string) {
	t.builder.Reset()
	sep := "+"
	for i, field := range row {
		t.builder.WriteString(sep)
		v := " "
		if field != "" {
			v = "-"
		}
		segment := strings.Repeat(v, t.columnWidths[i]+t.margin*2)
		t.builder.WriteString(segment)
	}
	t.builder.WriteString(sep)
	fmt.Fprintln(t.writer, t.builder.String())
}

// printBorder renders the table border based on column widths.
func (t *Table) printBorder() {
	fmt.Fprintln(t.writer, t.border)
}

// setAttr configures placeholders and delimiters based on the table format.
// It ensures consistency in the appearance and structure of table output.
func (t *Table) setAttr() {
	p, d := t.getDefaultAttr()
	if t.emptyFieldPlaceholder == DefaultEmptyFieldPlaceholder {
		t.emptyFieldPlaceholder = p
	}
	if t.wordDelimiter == DefaultWordDelimiter {
		t.wordDelimiter = d
	}
}

// setHeader extracts field names from the struct type to create the table header.
// It also initializes column widths based on the header names.
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

// setData converts the input data to a matrix of strings and calculates column widths.
// It also handles field formatting based on the table format and whether fields are merged.
func (t *Table) setData(v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("input must be a slice or a pointer to a slice")
	}
	t.data = make([][]string, v.Len())
	prev := make([]string, len(t.header))
	for i := 0; i < v.Len(); i++ {
		row := make([]string, len(t.header))
		field := v.Index(i)
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		merge := true
		for j, h := range t.header {
			field := field.FieldByName(h)
			if !field.IsValid() {
				return fmt.Errorf("field \"%s\" does not exist", h)
			}
			f, err := t.formatField(field)
			if err != nil {
				return fmt.Errorf("failed to format field \"%s\": %w", h, err)
			}
			if slices.Contains(t.mergedFields, j) {
				if f != prev[j] {
					merge = false
					prev[j] = f
				}
				if merge {
					f = ""
				}
			}
			row[j] = f
			for _, line := range strings.Split(f, "\n") {
				lw := runewidth.StringWidth(line)
				if lw > t.columnWidths[j] {
					t.columnWidths[j] = lw
				}
			}
		}
		t.data[i] = row
	}
	return nil
}

// setBorder computes the table border string based on the calculated column widths.
func (t *Table) setBorder() {
	var sep string
	switch t.format {
	case FormatMarkdown, FormatBacklog:
		sep = "|"
	default:
		sep = "+"
	}
	t.builder.Reset()
	for _, width := range t.columnWidths {
		t.builder.WriteString(sep)
		t.builder.WriteString(strings.Repeat("-", width+t.margin*2))
	}
	t.builder.WriteString(sep)
	t.border = t.builder.String()
}

// getDefaultAttr returns the default empty field placeholder and word delimiter for the table's current format.
func (t *Table) getDefaultAttr() (emptyFieldPlaceholder string, wordDelimiter string) {
	switch t.format {
	case FormatMarkdown:
		return MarkdownDefaultEmptyFieldPlaceholder, MarkdownDefaultWordDelimiter
	case FormatBacklog:
		return BacklogDefaultEmptyFieldPlaceholder, BacklogDefaultWordDelimiter
	default:
		return DefaultEmptyFieldPlaceholder, DefaultWordDelimiter
	}
}

// getMargin returns a string consisting of spaces to be used as margin around cell content.
func (t *Table) getMargin() string {
	return strings.Repeat(" ", t.margin)
}

// formatField formats a single field value based on its type and the table's configuration.
// It applies escaping if enabled and handles various data types, including slices and primitive types.
func (t *Table) formatField(v reflect.Value) (string, error) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return t.emptyFieldPlaceholder, nil
		}
		v = v.Elem()
	}
	var ret string
	switch v.Kind() {
	case reflect.String:
		s := v.String()
		if t.hasEscape {
			s = t.escape(s)
		}
		if t.format == FormatMarkdown && strings.HasPrefix(s, "*") {
			s = "\\" + s
		}
		if s == "" {
			s = t.emptyFieldPlaceholder
		}
		ret = strings.TrimSpace(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ret = fmt.Sprint(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ret = fmt.Sprint(v.Uint())
	case reflect.Float32, reflect.Float64:
		ret = fmt.Sprint(v.Float())
	case reflect.Slice:
		switch {
		case v.Len() == 0:
			ret = t.emptyFieldPlaceholder
		case v.Type().Elem().Kind() == reflect.Uint8:
			ret = string(v.Bytes())
		default:
			var sl []string
			for i := 0; i < v.Len(); i++ {
				e := v.Index(i)
				if e.Kind() == reflect.Ptr {
					if e.IsNil() {
						sl = append(sl, t.emptyFieldPlaceholder)
						continue
					}
					e = e.Elem()
				}
				if e.Kind() == reflect.Slice || e.Kind() == reflect.Struct {
					return "", fmt.Errorf("field must not be nested")
				}
				f, err := t.formatField(e)
				if err != nil {
					return "", err
				}
				sl = append(sl, f)
			}
			ret = strings.Join(sl, t.wordDelimiter)
		}
	default:
		ret = fmt.Sprint(v.Interface())
	}
	if v.Kind() == reflect.String || v.Kind() == reflect.Slice {
		ret = t.replaceNL(ret)
	}
	return strings.TrimSpace(ret), nil
}

// getDefaultAttr returns the default empty field placeholder and word delimiter for the table's current format.
func (t *Table) replaceNL(s string) string {
	_, d := t.getDefaultAttr()
	if d == "\n" {
		return s
	}
	t.builder.Reset()
	for _, r := range s {
		switch r {
		case '\n':
			t.builder.WriteString(d)
		default:
			t.builder.WriteRune(r)
		}
	}
	return t.builder.String()
}

// escape applies HTML escaping to a string for safe rendering in Markdown and other formats.
func (t *Table) escape(s string) string {
	t.builder.Reset()
	for _, r := range s {
		switch r {
		case '<':
			t.builder.WriteString("&lt;")
		case '>':
			t.builder.WriteString("&gt;")
		case '"':
			t.builder.WriteString("&quot;")
		case '\'':
			t.builder.WriteString("&lsquo;")
		case '&':
			t.builder.WriteString("&amp;")
		case ' ':
			t.builder.WriteString("&nbsp;")
		case '*':
			t.builder.WriteString("&#42;")
		case '\\':
			t.builder.WriteString("&#92;")
		case '_':
			t.builder.WriteString("&#95;")
		case '|':
			t.builder.WriteString("&#124;")
		default:
			t.builder.WriteRune(r)
		}
	}
	return t.builder.String()
}

// pad right-aligns numeric strings and left-aligns all other strings within a field of specified width.
func pad(s string, w int) string {
	if isNum(s) {
		return padL(s, w)
	}
	return padR(s, w)
}

// padR left-aligns a string within a field of specified width.
func padR(s string, w int) string {
	p := w - runewidth.StringWidth(s)
	if p > 0 {
		return s + strings.Repeat(" ", p)
	}
	return s
}

// padL right-aligns a string within a field of specified width, primarily for numeric data.
func padL(s string, w int) string {
	p := w - runewidth.StringWidth(s)
	if p > 0 {
		return strings.Repeat(" ", p) + s
	}
	return s
}

// isNum checks if a string represents a numeric value.
func isNum(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}
