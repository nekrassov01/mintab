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
	writer                io.Writer  // Destination for table output.
	data                  [][]string // Data holds the content of the table.
	header                []string   // Names of each field in the table header.
	format                Format     // Output format of the table.
	border                string     // Pre-computed border based on column widths.
	tableWidth            int        //
	marginWidth           int        //
	margin                string     // Margin size around cell content.
	emptyFieldPlaceholder string     // Placeholder for empty fields.
	wordDelimiter         string     // Delimiter for words within a field.
	mergedFields          []int      // Indices of fields to be merged based on content.
	ignoredFields         []int      // Indices of fields to be ignored during rendering.
	columnWidths          []int      // Calculated max width of each column.
	hasHeader             bool       // Indicates if the header should be rendered.
	hasEscape             bool       // Indicates if escaping should be performed.
}

// New instantiates a new Table with the specified writer and options.
func New(w io.Writer, opts ...Option) *Table {
	t := &Table{
		writer:                w,
		format:                FormatText,
		marginWidth:           1,
		emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
		wordDelimiter:         DefaultWordDelimiter,
		hasHeader:             true,
	}
	for _, opt := range opts {
		opt(t)
	}
	t.margin = strings.Repeat(" ", t.marginWidth)
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
	if t.marginWidth < 0 {
		return fmt.Errorf("only unsigned integers are allowed in margin")
	}
	if _, ok := input.([]interface{}); ok {
		return fmt.Errorf("elements of slice must not be empty interface")
	}
	rv := reflect.ValueOf(input)
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
}

// printHeader renders the table header.
func (t *Table) printHeader() {
	var b strings.Builder
	b.Grow(t.tableWidth)
	b.WriteString("|")
	for i, h := range t.header {
		t.pad(&b, h, t.columnWidths[i])
		b.WriteString("|")
	}
	if t.format == FormatBacklog {
		b.WriteString("h")
	}
	fmt.Fprintln(t.writer, b.String())
}

// printData renders the table data with dynamic conditional borders.
func (t *Table) printData() {
	for i, row := range t.data {
		if i > 0 {
			if t.format == FormatText {
				t.printDataBorder(row)
			}
			if t.format == FormatCompressedText && row[0] != "" {
				t.printBorder()
			}
		}
		splited := make([][]string, len(row))
		n := 1
		for j, field := range row {
			splited[j] = strings.Split(field, "\n")
			if len(splited[j]) > n {
				n = len(splited[j])
			}
		}
		for k := 0; k < n; k++ {
			var b strings.Builder
			b.Grow(t.tableWidth)
			b.WriteString("|")
			for l, elem := range splited {
				if k < len(elem) {
					t.pad(&b, elem[k], t.columnWidths[l])
				} else {
					t.pad(&b, "", t.columnWidths[l])
				}
				b.WriteString("|")
			}
			fmt.Fprintln(t.writer, b.String())
		}
	}
}

// printDataBorder prints a conditional border based on the emptiness of fields in the current row,
// with continuity in border characters based on the emptiness of adjacent fields.
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
	fmt.Fprintln(t.writer, b.String())
}

// printBorder renders the table border based on column widths.
func (t *Table) printBorder() {
	fmt.Fprintln(t.writer, t.border)
}

// setAttr configures placeholders and delimiters based on the table format.
// It ensures consistency in the appearance and structure of table output.
func (t *Table) setAttr() {
	var p, d string
	switch t.format {
	case FormatMarkdown:
		p = MarkdownDefaultEmptyFieldPlaceholder
		d = MarkdownDefaultWordDelimiter
	case FormatBacklog:
		p = BacklogDefaultEmptyFieldPlaceholder
		d = BacklogDefaultWordDelimiter
	default:
		p = DefaultEmptyFieldPlaceholder
		d = DefaultWordDelimiter
	}
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
func (t *Table) setData(rv reflect.Value) error {
	t.data = make([][]string, rv.Len())
	prev := make([]string, len(t.header))
	for i := 0; i < rv.Len(); i++ {
		row := make([]string, len(t.header))
		field := rv.Index(i)
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		merge := true
		for j, h := range t.header {
			field := field.FieldByName(h)
			if !field.IsValid() {
				return fmt.Errorf("invalid field detected: %s", h)
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
			elems := strings.Split(f, "\n")
			for _, elem := range elems {
				lw := runewidth.StringWidth(elem)
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
	var b strings.Builder
	for _, width := range t.columnWidths {
		b.WriteString(sep)
		for i := 0; i < width+t.marginWidth*2; i++ {
			b.WriteByte('-')
		}
	}
	b.WriteString(sep)
	t.border = b.String()
	t.tableWidth = len(t.border)
}

// formatField formats a single field value based on its type and the table's configuration.
// It applies escaping if enabled and handles various data types, including slices and primitive types.
func (t *Table) formatField(rv reflect.Value) (string, error) {
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return t.emptyFieldPlaceholder, nil
		}
		rv = rv.Elem()
	}
	v := getString(rv)
	if v == "" {
		switch rv.Kind() {
		case reflect.String:
			v = rv.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v = strconv.FormatInt(rv.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v = strconv.FormatUint(rv.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			v = strconv.FormatFloat(rv.Float(), 'f', -1, 64)
		case reflect.Slice, reflect.Array:
			switch {
			case rv.Len() == 0:
				v = t.emptyFieldPlaceholder
			case rv.Type().Elem().Kind() == reflect.Uint8:
				v = string(rv.Bytes())
			default:
				ss := make([]string, rv.Len())
				for i := 0; i < rv.Len(); i++ {
					e := rv.Index(i)
					if i != 0 {
						ss[i] = t.wordDelimiter
					}
					if e.Kind() == reflect.Ptr {
						if e.IsNil() {
							ss[i] = t.emptyFieldPlaceholder
							continue
						}
						e = e.Elem()
					}
					if s := getString(e); s != "" {
						ss[i] = s
						continue
					}
					if e.Kind() == reflect.Slice && e.Type().Elem().Kind() == reflect.Uint8 {
						ss[i] = string(e.Bytes())
						continue
					}
					if e.Kind() == reflect.Slice || e.Kind() == reflect.Array || e.Kind() == reflect.Struct {
						return "", fmt.Errorf("cannot represent nested fields")
					}
					f, err := t.formatField(e)
					if err != nil {
						return "", err
					}
					ss[i] = f
				}
				v = strings.Join(ss, t.wordDelimiter)
			}
		default:
			v = fmt.Sprint(rv.Interface())
		}
	}
	if t.hasEscape {
		v = t.escape(v)
	}
	if t.format == FormatMarkdown && strings.HasPrefix(v, "*") {
		v = "\\" + v
	}
	if v == "" {
		v = t.emptyFieldPlaceholder
	}
	return strings.TrimSpace(t.replaceNL(v)), nil
}

func getString(v reflect.Value) string {
	if v.CanInterface() {
		if s, ok := v.Interface().(fmt.Stringer); ok {
			return s.String()
		}
	}
	return ""
}

func (t *Table) replaceNL(s string) string {
	if t.wordDelimiter == "\n" {
		return s
	}
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '\n':
			b.WriteString(t.wordDelimiter)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// escape applies HTML escaping to a string for safe rendering in Markdown and other formats.
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

// pad right-aligns numeric strings and left-aligns all other strings within a field of specified width.
func (t *Table) pad(b *strings.Builder, s string, width int) {
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

// isNum checks if a string represents a numeric value.
func isNum(s string) bool {
	if len(s) == 0 {
		return false
	}
	start := 0
	if s[0] == '-' || s[0] == '+' {
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
