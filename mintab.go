package mintab

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// Version of mintab.
const Version = "0.0.5"

const (
	MarkdownFormat = iota
	BacklogFormat
)

const (
	DarkTheme = iota
	LightTheme
	NoneTheme
)

// Table represents a table in a matrix of strings.
type Table struct {
	data                  [][]string // Data holds the table data in a matrix of strings.
	headers               []string   // headers holds the name of each field in the table header.
	format                int        // format specifies the format of the table.
	theme                 int        // theme specifies the theme of the table.
	hasHeader             bool       // hasHeader indicates whether to enable header or not.
	emptyFieldPlaceholder string     // emptyFieldPlaceholder specifies the placeholder if the field is empty.
	wordDelimiter         string     // wordDelimiter specifies the word delimiter of the field.
	mergedFields          []int      // mergedFields holds indices of the field to be grouped.
	ignoredFields         []int      // ignoredFields holds indices of the fields to be ignored.
	colorFlags            []bool     // colorFlags holds flags indicating whether to color each row or not.
}

// New instantiates a table struct.
func NewTable(opts ...Option) (t *Table) {
	t = &Table{}
	t.format = MarkdownFormat
	t.theme = NoneTheme
	t.hasHeader = true
	t.emptyFieldPlaceholder = "N/A"
	t.wordDelimiter = "<br>"
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Option is the type for passing options when instantiating Table.
type Option func(*Table)

// WithFormat specifies the table format.
func WithFormat(format int) Option {
	return func(t *Table) {
		t.format = format
	}
}

// WithTheme specifies the table theme.
func WithTheme(theme int) Option {
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
func (t *Table) Load(input any) (err error) {
	if _, ok := input.([]interface{}); ok {
		return fmt.Errorf("cannot parse input: must not be slice of empty interface")
	}
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("cannot parse input: must be slice")
	}
	if v.Len() == 0 {
		return fmt.Errorf("cannot parse input: no elements in slice")
	}
	e := v.Index(0)
	if e.Kind() != reflect.Struct {
		return fmt.Errorf("cannot parse input: must be struct")
	}
	t.colorFlags = getColorFlags(input)
	t.headers = t.setHeader(e.Type())
	t.data, err = t.setData(input)
	if err != nil {
		return err
	}
	return nil
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
	if t.format == BacklogFormat {
		table.SetHeaderLine(false)
	}
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)
	table.AppendBulk(t.data)
	table.Render()
	s := t.colorize(tableString.String())
	if t.format == BacklogFormat {
		s = t.backlogify(s)
	}
	return s
}

// setHeader uses reflect to extract field names and create headers.
func (t *Table) setHeader(typ reflect.Type) (headers []string) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if contains(t.ignoredFields, i) || field.PkgPath != "" {
			continue
		}
		headers = append(headers, field.Name)
	}
	return headers
}

// setData checks whether what is passed using reflect is a slice of struct,
// formats the field values, and converts them to a string table.
func (t *Table) setData(input any) (data [][]string, err error) {
	v := reflect.ValueOf(input)
	prev := make([]string, len(t.headers))
	for i := 0; i < v.Len(); i++ {
		var values []string
		item := v.Index(i).Interface()
		merge := true
		for i, header := range t.headers {
			value, err := t.formatValue(reflect.ValueOf(item).FieldByName(header))
			if err != nil {
				return nil, fmt.Errorf("cannot parse field: %w", err)
			}
			if contains(t.mergedFields, i) {
				if value != prev[i] {
					merge = false
					prev[i] = value
				}
				if merge {
					value = ""
				}
			}
			values = append(values, value)
		}
		data = append(data, values)
	}
	return data, nil
}

// getColorFlags determines which rows to color.
func getColorFlags(input any) (colorFlags []bool) {
	v := reflect.ValueOf(input)
	var prev any
	colorFlag := true
	for i := 0; i < v.Len(); i++ {
		firstField := reflect.ValueOf(v.Index(i).Interface()).Field(0).Interface()
		if prev != nil && prev != firstField {
			colorFlag = !colorFlag
		}
		colorFlags = append(colorFlags, colorFlag)
		prev = firstField
	}
	return colorFlags
}

// formatValue formats the value of each field.
// Perform multi-value delimiters and whitespace handling.
// Nested fields are not processed and an error is returned.
func (t *Table) formatValue(v reflect.Value) (string, error) {
	if isEmptyStr(v) {
		return t.emptyFieldPlaceholder, nil
	}
	switch v.Kind() {
	case reflect.Slice:
		if v.IsNil() || v.Len() == 0 {
			return t.emptyFieldPlaceholder, nil
		}
		var s []string
		for i := 0; i < v.Len(); i++ {
			value := v.Index(i)
			if value.Kind() == reflect.Slice || value.Kind() == reflect.Struct {
				return "", fmt.Errorf("elements of slice must not be nested")
			}
			if isEmptyStr(value) {
				s = append(s, t.emptyFieldPlaceholder)
			} else {
				s = append(s, trim(v.Index(i)))
			}
		}
		return strings.Join(s, t.wordDelimiter), nil
	case reflect.Struct:
		return "", fmt.Errorf("field must not be struct")
	}
	return trim(v), nil
}

// backlogify converts to backlog tables format.
func (t *Table) backlogify(s string) string {
	if t.hasHeader {
		i := strings.Index(s, "\n")
		if i == -1 {
			return s
		}
		s = s[:i] + "h" + s[i:]
	}
	return strings.ReplaceAll(s, "<br>", "&br;")
}

// colorize adds color to the table.
func (t *Table) colorize(table string) string {
	if t.theme == NoneTheme {
		return table
	}
	offset := t.getOffset()
	color := t.getColor()
	var clines []string
	lines := strings.Split(table, "\n")
	for i, line := range lines {
		if i < offset {
			clines = append(clines, line)
			continue
		}
		if i-offset < len(t.colorFlags) {
			if t.colorFlags[i-offset] {
				clines = append(clines, color.Sprint(line))
			} else {
				clines = append(clines, line)
			}
		} else {
			clines = append(clines, line)
		}
	}
	return strings.Join(clines, "\n")
}

// getOffset determines the starting position for coloring.
func (t *Table) getOffset() int {
	if !t.hasHeader {
		return 0
	}
	switch t.format {
	case MarkdownFormat:
		return 2
	case BacklogFormat:
		return 1
	default:
		return 0
	}
}

// getColor sets the color theme.
func (t *Table) getColor() *color.Color {
	switch t.theme {
	case DarkTheme:
		return color.New(color.BgHiBlack, color.FgHiWhite)
	case LightTheme:
		return color.New(color.BgHiWhite, color.FgHiBlack)
	default:
		return color.New(color.Reset)
	}
}

// contains is a helper function used to determine grouping.
func contains(values []int, target int) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

// trim is a helper function used in formatting field values.
func trim(v reflect.Value) string {
	return strings.TrimSpace(fmt.Sprint(v))
}

// isEmptyStr is a helper function that checks if a field value is empty.
func isEmptyStr(v reflect.Value) bool {
	return v.Kind() == reflect.String && v.String() == ""
}
