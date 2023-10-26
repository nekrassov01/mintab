package mintab

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// Version of mintab.
const Version = "0.0.2"

// defaultEmptyFieldPlaceholder is a placeholder if the field has no value.
const defaultEmptyFieldPlaceholder = "N/A"

// The defaultWordDelimiter is the delimiter for multiple field values.
const defaultWordDelimitter = "<br>"

// TableFormat represents the format of the table.
type TableFormat string

// Markdown represents markdown table format.
const Markdown TableFormat = "markdown"

// Backlog represents backlog table format.
const Backlog TableFormat = "backlog"

// TableTheme represents the theme of the table.
type TableTheme string

// DarkTheme represents dark theme for the table.
const DarkTheme TableTheme = "dark"

// LightTheme represents light theme for the table.
const LightTheme TableTheme = "light"

// NoneTheme represents no theme for the table.
const NoneTheme TableTheme = "none"

// Table represents a table in a matrix of strings.
type Table struct {
	Data                  [][]string  // Data holds the table data in a matrix of strings.
	headers               []string    // headers holds the name of each field in the table header.
	format                TableFormat // format specifies the format of the table.
	theme                 TableTheme  // theme specifies the theme of the table.
	hasHeader             bool        // hasHeader indicates whether to enable header or not.
	emptyFieldPlaceholder string      // emptyFieldPlaceholder specifies the placeholder if the field is empty.
	wordDelimitter        string      // wordDelimitter specifies the word delimiter of the field.
	mergeFields           []int       // mergeFields holds indices of the field to be grouped.
	ignoreFields          []int       // ignoreFields holds indices of the fields to be ignored.
	colorFlags            []bool      // colorFlags holds flags indicating whether to color each row or not.
}

// New instantiates a table struct.
func New(input any, opts ...Option) *Table {
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Slice {
		return nil
	}
	if v.Len() == 0 {
		return nil
	}
	elem := v.Index(0).Interface()
	if reflect.TypeOf(elem).Kind() != reflect.Struct {
		return nil
	}
	t := &Table{
		Data:                  nil,
		headers:               nil,
		format:                Markdown,
		theme:                 NoneTheme,
		hasHeader:             true,
		emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
		wordDelimitter:        defaultWordDelimitter,
		mergeFields:           nil,
		ignoreFields:          nil,
		colorFlags:            getColorFlags(input),
	}
	for _, opt := range opts {
		opt(t)
	}
	t.headers = createHeaders(reflect.TypeOf(elem), t.ignoreFields)
	t.Data = createDataRows(input, t.headers, t.mergeFields, t.emptyFieldPlaceholder, t.wordDelimitter)
	return t
}

// Option is the type for passing options when instantiating Table.
type Option func(*Table)

// WithTableFormat specifies the table format.
func WithTableFormat(format TableFormat) Option {
	return func(t *Table) {
		t.format = format
	}
}

// WithTableTheme specifies the table theme.
func WithTableTheme(theme TableTheme) Option {
	return func(t *Table) {
		t.theme = theme
	}
}

// WithTableHeader enables/disables the header.
func WithTableHeader(has bool) Option {
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
func WithWordDelimitter(wordDelimitter string) Option {
	return func(t *Table) {
		t.wordDelimitter = wordDelimitter
	}
}

// WithMergeFields specifies the column numbers to be used for grouping.
// It depends on whether the first field value is the same as the previous row.
func WithMergeFields(mergeFields []int) Option {
	return func(t *Table) {
		t.mergeFields = mergeFields
	}
}

// WithIgnoreFields specifies the column numbers to be ignored.
// It depends on whether the first field value is the same as the previous row.
func WithIgnoreFields(ignoreFields []int) Option {
	return func(t *Table) {
		t.ignoreFields = ignoreFields
	}
}

// Out outputs the table as a string.
// It can be used as a markdown table or backlog table by copying and pasting.
func (t *Table) Out() string {
	if t == nil {
		return ""
	}
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	if t.hasHeader {
		table.SetHeader(t.headers)
		table.SetAutoFormatHeaders(false)
	}
	if t.format == Backlog {
		table.SetHeaderLine(false)
	}
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)
	table.AppendBulk(t.Data)
	table.Render()
	s := colorize(tableString.String(), t.colorFlags, t.format, t.theme, t.hasHeader)
	if t.format == Backlog {
		s = backlogify(s, t.hasHeader)
	}
	return s
}

// createHeaders uses reflect to extract field names and create headers.
func createHeaders(typ reflect.Type, ignoreFields []int) (headers []string) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if contains(ignoreFields, i) || field.PkgPath != "" {
			continue
		}
		headers = append(headers, field.Name)
	}
	return headers
}

// createDataRows checks whether what is passed using reflect is a slice of the structure,
// formats the field values, and converts them to a string table.
func createDataRows(input any, headers []string, mergeFields []int, emptyFieldPlaceholder string, wordDelimitter string) (data [][]string) {
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Slice {
		return nil
	}
	prev := make([]string, len(headers))
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		var values []string
		merge := true
		for i, header := range headers {
			value := formatValue(reflect.ValueOf(item).FieldByName(header), emptyFieldPlaceholder, wordDelimitter)
			if contains(mergeFields, i) {
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
	return data
}

// getColorFlags determines which rows to color.
func getColorFlags(input any) (colorFlags []bool) {
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Slice {
		return nil
	}
	var prev any
	colorFlag := true
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		firstField := reflect.ValueOf(item).Field(0).Interface()
		if prev != nil && prev != firstField {
			colorFlag = !colorFlag
		}
		colorFlags = append(colorFlags, colorFlag)
		prev = firstField
	}
	return colorFlags
}

// formatValue formats the value of each field.
func formatValue(v reflect.Value, emptyFieldPlaceholder string, wordDelimitter string) string {
	if isEmptyStr(v) {
		return emptyFieldPlaceholder
	}
	if v.Kind() == reflect.Slice {
		if v.IsNil() || v.Len() == 0 {
			return emptyFieldPlaceholder
		}
		var s []string
		for i := 0; i < v.Len(); i++ {
			val := v.Index(i)
			if isEmptyStr(val) {
				s = append(s, emptyFieldPlaceholder)
			} else {
				s = append(s, trim(v.Index(i)))
			}
		}
		return strings.Join(s, wordDelimitter)
	}
	return trim(v)
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

// backlogify converts to backlog tables format.
func backlogify(s string, hasHeader bool) string {
	if hasHeader {
		i := strings.Index(s, "\n")
		if i == -1 {
			return s
		}
		s = s[:i] + "h" + s[i:]
	}
	return strings.ReplaceAll(s, "<br>", "&br;")
}

// colorize adds color to the table.
func colorize(table string, colorFlags []bool, tableFormat TableFormat, tableTheme TableTheme, hasHeader bool) string {
	if tableTheme == NoneTheme {
		return table
	}
	offset := getOffset(tableFormat, hasHeader)
	var coloredLines []string
	lines := strings.Split(table, "\n")
	for i, line := range lines {
		if i < offset {
			coloredLines = append(coloredLines, line)
			continue
		}
		if i-offset < len(colorFlags) {
			if colorFlags[i-offset] {
				coloredLines = append(coloredLines, getColor(tableTheme).Sprint(line))
			} else {
				coloredLines = append(coloredLines, line)
			}
		} else {
			coloredLines = append(coloredLines, line)
		}
	}
	return strings.Join(coloredLines, "\n")
}

// getOffset determines the starting position for coloring.
func getOffset(tableFormat TableFormat, hasHeader bool) int {
	if !hasHeader {
		return 0
	}
	switch tableFormat {
	case Markdown:
		return 2
	case Backlog:
		return 1
	default:
		return 0
	}
}

// getColor sets the color theme.
func getColor(tableTheme TableTheme) *color.Color {
	switch tableTheme {
	case DarkTheme:
		return color.New(color.BgHiBlack, color.FgHiWhite)
	case LightTheme:
		return color.New(color.BgHiWhite, color.FgHiBlack)
	default:
		return color.New(color.Reset)
	}
}
