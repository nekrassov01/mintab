package mintab

import (
	"fmt"
	"html"
	"reflect"
	"slices"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// Format defines the output format of the content.
type Format int

// Enumeration of possible output formats.
const (
	FormatText     Format = iota // FormatText represents plain text format.
	FormatMarkdown               // FormatMarkdown represents markdown format.
	FormatBacklog                // FormatBacklog represents Backlog-specific format.
)

// Formats holds the string representation of each Format constant.
var Formats = []string{
	"text",
	"markdown",
	"backlog",
}

// String returns the string representation of the Format.
// If the format is not within the range of predefined formats, an empty string is returned.
func (o Format) String() string {
	if o >= 0 && int(o) < len(Formats) {
		return Formats[o]
	}
	return ""
}

// Theme defines the visual theme preference.
type Theme int

// Enumeration of possible visual themes.
const (
	ThemeNone  Theme = iota // ThemeNone indicates no preference for a theme.
	ThemeDark               // ThemeDark indicates preference for a dark theme.
	ThemeLight              // ThemeLight indicates preference for a light theme.
)

// Themes holds the string representation of each Theme constant.
var Themes = []string{
	"none",
	"dark",
	"light",
}

// String returns the string representation of the Theme.
// If the theme is not within the range of predefined themes, an empty string is returned.
func (t Theme) String() string {
	if t >= 0 && int(t) < len(Themes) {
		return Themes[t]
	}
	return ""
}

// Dafault values
const (
	DefaultEmptyFieldPlaceholder         = "-"
	DefaultWordDelimiter                 = "\n"
	MarkdownDefaultEmptyFieldPlaceholder = "&#45;"
	MarkdownDefaultWordDelimiter         = "<br>"
	BacklogDefaultEmptyFieldPlaceholder  = "-"
	BacklogDefaultWordDelimiter          = "&br;"
)

// Table represents a table in a matrix of strings.
type Table struct {
	data                  [][]string // Data holds the table data in a matrix of strings.
	headers               []string   // headers holds the name of each field in the table header.
	format                Format     // format specifies the format of the table.
	theme                 Theme      // theme specifies the theme of the table.
	hasHeader             bool       // hasHeader indicates whether to enable header or not.
	emptyFieldPlaceholder string     // emptyFieldPlaceholder specifies the placeholder if the field is empty.
	wordDelimiter         string     // wordDelimiter specifies the word delimiter of the field.
	mergedFields          []int      // mergedFields holds indices of the field to be grouped.
	ignoredFields         []int      // ignoredFields holds indices of the fields to be ignored.
}

// NewTable instantiates a table struct.
func NewTable(opts ...Option) *Table {
	t := &Table{
		format:                FormatText,
		theme:                 ThemeNone,
		hasHeader:             true,
		emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
		wordDelimiter:         DefaultWordDelimiter,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Option is the type for passing options when instantiating Table.
type Option func(*Table)

// WithFormat specifies the table format.
func WithFormat(format Format) Option {
	return func(t *Table) {
		t.format = format
	}
}

// WithTheme specifies the table theme.
func WithTheme(theme Theme) Option {
	return func(t *Table) {
		t.theme = theme
	}
}

// WithHeader enables/disables the header.
func WithHeader(has bool) Option {
	return func(t *Table) {
		t.hasHeader = has
	}
}

// WithEmptyFieldPlaceholder specifies the value if the field is empty.
func WithEmptyFieldPlaceholder(emptyFieldPlaceholder string) Option {
	return func(t *Table) {
		t.emptyFieldPlaceholder = emptyFieldPlaceholder
	}
}

// WithWordDelimiter specifies the delimiter when a field has multiple values.
func WithWordDelimiter(wordDelimiter string) Option {
	return func(t *Table) {
		t.wordDelimiter = wordDelimiter
	}
}

// WithMergeFields specifies the column numbers to be used for grouping.
// determined by whether the value in the first column is the same as in the previous row.
func WithMergeFields(mergeFields []int) Option {
	return func(t *Table) {
		t.mergedFields = mergeFields
	}
}

// WithIgnoreFields specifies the column numbers to be ignored.
func WithIgnoreFields(ignoreFields []int) Option {
	return func(t *Table) {
		t.ignoredFields = ignoreFields
	}
}

// Load validates input and converts them to table data.
// Returns error if not struct slice.
func (t *Table) Load(input any) error {
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
	if _, ok := input.([]interface{}); ok {
		return fmt.Errorf("cannot parse input: elements of slice must not be empty interface")
	}
	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("cannot parse input: must be a slice or a pointer to a slice")
	}
	if v.Len() == 0 {
		return fmt.Errorf("cannot parse input: no data found")
	}
	e := v.Index(0)
	if e.Kind() == reflect.Ptr {
		e = e.Elem()
	}
	if e.Kind() != reflect.Struct {
		return fmt.Errorf("cannot parse input: elements of slice must be struct or pointer to struct")
	}
	t.setHeader(e.Type())
	return t.setData(v)
}

// Out outputs the table as a string.
// It can be used as a markdown table or backlog table by copying and pasting.
func (t *Table) Out() string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	if t.hasHeader {
		table.SetHeader(t.headers)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoFormatHeaders(false)
	}
	if t.format == FormatBacklog {
		table.SetHeaderLine(false)
	}
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)
	table.AppendBulk(t.data)
	table.Render()
	s := t.colorize(tableString.String())
	if t.format == FormatBacklog {
		s = t.backlogify(s)
	}
	return s
}

// setHeader uses reflect to extract field names and create headers.
func (t *Table) setHeader(typ reflect.Type) {
	if len(t.headers) > 0 {
		return
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if slices.Contains(t.ignoredFields, i) || field.PkgPath != "" {
			continue
		}
		t.headers = append(t.headers, field.Name)
	}
}

// setData converts input to a matrix of strings.
func (t *Table) setData(v reflect.Value) error {
	prev := make([]string, len(t.headers))
	t.data = make([][]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		values := make([]string, len(t.headers))
		itemValue := v.Index(i)
		if itemValue.Kind() == reflect.Ptr {
			itemValue = itemValue.Elem()
		}
		merge := true
		for j, header := range t.headers {
			field := itemValue.FieldByName(header)
			if !field.IsValid() {
				return fmt.Errorf("field \"%s\" does not exist", header)
			}
			value, err := t.formatValue(field)
			if err != nil {
				return fmt.Errorf("cannot format field \"%s\": %w", header, err)
			}
			if slices.Contains(t.mergedFields, j) {
				if value != prev[j] {
					merge = false
					prev[j] = value
				}
				if merge {
					value = ""
				}
			}
			values[j] = value
		}
		t.data[i] = values
	}
	return nil
}

// formatValue formats the value of each field.
// Perform multi-value delimiters and whitespace handling.
// Nested fields are not processed and an error is returned.
func (t *Table) formatValue(v reflect.Value) (string, error) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return t.emptyFieldPlaceholder, nil
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.String:
		s := v.String()
		if s == "" {
			return t.emptyFieldPlaceholder, nil
		}
		if t.format == FormatMarkdown {
			s = html.EscapeString(s)
			r := strings.NewReplacer(" ", "&nbsp;", "|", "&#124;", "*", "&#42;", "\\", "&#92;", "_", "&#095;")
			s = r.Replace(s)
		}
		return strings.TrimSpace(s), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprint(v.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprint(v.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return fmt.Sprint(v.Float()), nil
	case reflect.Slice:
		if v.Len() == 0 {
			return t.emptyFieldPlaceholder, nil
		}
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return string(v.Bytes()), nil
		}
		var s []string
		for i := 0; i < v.Len(); i++ {
			e := v.Index(i)
			if e.Kind() == reflect.Ptr {
				if e.IsNil() {
					s = append(s, t.emptyFieldPlaceholder)
					continue
				}
				e = e.Elem()
			}
			if e.Kind() == reflect.Slice || e.Kind() == reflect.Struct {
				return "", fmt.Errorf("elements of slice must not be nested")
			}
			fv, err := t.formatValue(e)
			if err != nil {
				return "", err
			}
			s = append(s, fv)
		}
		return strings.Join(s, t.wordDelimiter), nil
	default:
		return fmt.Sprint(v.Interface()), nil
	}
}

// backlogify converts to backlog table format.
func (t *Table) backlogify(s string) string {
	if t.hasHeader {
		i := strings.Index(s, "\n")
		if i == -1 {
			return s
		}
		s = s[:i] + "h" + s[i:]
	}
	return s
}

// colorize adds color to the table, row by row based on the first field value.
func (t *Table) colorize(table string) string {
	if t.theme == ThemeNone {
		return table
	}
	lines := strings.Split(table, "\n")
	var coloredLines []string
	var lastNonEmptyFirstFieldValue string
	var currentColorFlag bool
	for i, line := range lines {
		if i < t.getOffset() {
			coloredLines = append(coloredLines, line)
			continue
		}
		fields := strings.Split(line, "|")
		if len(fields) <= 1 {
			coloredLines = append(coloredLines, line)
			continue
		}
		firstFieldValue := strings.TrimSpace(fields[1])
		if firstFieldValue == "" {
			firstFieldValue = lastNonEmptyFirstFieldValue
		}
		if firstFieldValue != lastNonEmptyFirstFieldValue {
			currentColorFlag = !currentColorFlag
			lastNonEmptyFirstFieldValue = firstFieldValue
		}
		if currentColorFlag {
			line = t.getColor().Sprint(line)
		}
		coloredLines = append(coloredLines, line)
	}
	return strings.Join(coloredLines, "\n")
}

// getOffset determines the starting position for coloring.
func (t *Table) getOffset() int {
	if !t.hasHeader {
		return 0
	}
	switch t.format {
	case FormatText, FormatMarkdown:
		return 2
	case FormatBacklog:
		return 1
	default:
		return 0
	}
}

// getColor sets the color theme.
func (t *Table) getColor() *color.Color {
	switch t.theme {
	case ThemeNone:
		return &color.Color{}
	case ThemeDark:
		return color.New(color.BgHiBlack, color.FgHiWhite)
	case ThemeLight:
		return color.New(color.BgHiWhite, color.FgHiBlack)
	default:
		return &color.Color{}
	}
}
