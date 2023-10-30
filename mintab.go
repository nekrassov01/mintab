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
const defaultWordDelimiter = "<br>"

// markdownHeaderOffset is starting position of data rows in markdown table.
const markdownHeaderOffset = 2

// backlogHeaderOffset is starting position of data rows in backlog table.
const backlogHeaderOffset = 1

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
	wordDelimiter         string      // wordDelimiter specifies the word delimiter of the field.
	mergeFields           []int       // mergeFields holds indices of the field to be grouped.
	ignoreFields          []int       // ignoreFields holds indices of the fields to be ignored.
	colorFlags            []bool      // colorFlags holds flags indicating whether to color each row or not.
}

// New instantiates a table struct.
func New(input any, opts ...Option) (table *Table, err error) {
	table = &Table{}
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("cannot parse input: must be slice")
	}
	if v.Len() == 0 {
		return nil, fmt.Errorf("cannot parse input: no elements in slice")
	}
	elem := v.Index(0).Interface()
	e := reflect.TypeOf(elem)
	if e.Kind() != reflect.Struct {
		return nil, fmt.Errorf("cannot parse input: must be struct")
	}
	table.format = Markdown
	table.theme = NoneTheme
	table.hasHeader = true
	table.emptyFieldPlaceholder = defaultEmptyFieldPlaceholder
	table.wordDelimiter = defaultWordDelimiter
	table.colorFlags = createColorFlags(input)
	for _, opt := range opts {
		opt(table)
	}
	table.headers = createHeader(e, table.ignoreFields)
	table.Data, err = createData(input, table.headers, table.mergeFields, table.emptyFieldPlaceholder, table.wordDelimiter)
	if err != nil {
		return nil, err
	}
	return table, nil
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
func WithWordDelimiter(wordDelimiter string) Option {
	return func(t *Table) {
		t.wordDelimiter = wordDelimiter
	}
}

// WithMergeFields specifies the column numbers to be used for grouping.
// determined by whether the value in the first column is the same as in the previous row.
func WithMergeFields(mergeFields []int) Option {
	return func(t *Table) {
		t.mergeFields = mergeFields
	}
}

// WithIgnoreFields specifies the column numbers to be ignored.
func WithIgnoreFields(ignoreFields []int) Option {
	return func(t *Table) {
		t.ignoreFields = ignoreFields
	}
}

// Out outputs the table as a string.
// It can be used as a markdown table or backlog table by copying and pasting.
func (t *Table) Out() (string, error) {
	if t == nil || t.Data == nil {
		return "", fmt.Errorf("cannot parse table: empty data")
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
	s, err := colorize(tableString.String(), t.colorFlags, t.format, t.theme, t.hasHeader)
	if err != nil {
		return "", err
	}
	if t.format == Backlog {
		s = backlogify(s, t.hasHeader)
	}
	return s, nil
}

// createHeader uses reflect to extract field names and create headers.
func createHeader(t reflect.Type, ignoreFields []int) (headers []string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if contains(ignoreFields, i) || field.PkgPath != "" {
			continue
		}
		headers = append(headers, field.Name)
	}
	return headers
}

// createData checks whether what is passed using reflect is a slice of struct,
// formats the field values, and converts them to a string table.
func createData(input any, headers []string, mergeFields []int, emptyFieldPlaceholder string, wordDelimiter string) (data [][]string, err error) {
	v := reflect.ValueOf(input)
	prev := make([]string, len(headers))
	for i := 0; i < v.Len(); i++ {
		var values []string
		item := v.Index(i).Interface()
		merge := true
		for i, header := range headers {
			value, err := formatValue(reflect.ValueOf(item).FieldByName(header), emptyFieldPlaceholder, wordDelimiter)
			if err != nil {
				return nil, fmt.Errorf("cannot parse field: %w", err)
			}
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
	return data, nil
}

// createColorFlags determines which rows to color.
func createColorFlags(input any) (colorFlags []bool) {
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
func formatValue(v reflect.Value, emptyFieldPlaceholder string, wordDelimiter string) (string, error) {
	if isEmptyStr(v) {
		return emptyFieldPlaceholder, nil
	}
	switch v.Kind() {
	case reflect.Slice:
		if v.IsNil() || v.Len() == 0 {
			return emptyFieldPlaceholder, nil
		}
		var s []string
		for i := 0; i < v.Len(); i++ {
			value := v.Index(i)
			if value.Kind() == reflect.Slice || value.Kind() == reflect.Struct {
				return "", fmt.Errorf("elements of slice must not be nested")
			}
			if isEmptyStr(value) {
				s = append(s, emptyFieldPlaceholder)
			} else {
				s = append(s, trim(v.Index(i)))
			}
		}
		return strings.Join(s, wordDelimiter), nil
	case reflect.Struct:
		return "", fmt.Errorf("field must not be struct")
	}
	return trim(v), nil
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
func colorize(table string, colorFlags []bool, tableFormat TableFormat, tableTheme TableTheme, hasHeader bool) (string, error) {
	if tableTheme == NoneTheme {
		return table, nil
	}
	offset, err := getOffset(tableFormat, hasHeader)
	if err != nil {
		return "", fmt.Errorf("cannot parse table string: %w", err)
	}
	color, err := getColor(tableTheme)
	if err != nil {
		return "", fmt.Errorf("cannot parse table string: %w", err)
	}
	var clines []string
	lines := strings.Split(table, "\n")
	for i, line := range lines {
		if i < offset {
			clines = append(clines, line)
			continue
		}
		if i-offset < len(colorFlags) {
			if colorFlags[i-offset] {
				clines = append(clines, color.Sprint(line))
			} else {
				clines = append(clines, line)
			}
		} else {
			clines = append(clines, line)
		}
	}
	return strings.Join(clines, "\n"), nil
}

// getOffset determines the starting position for coloring.
func getOffset(tableFormat TableFormat, hasHeader bool) (int, error) {
	if !hasHeader {
		return 0, nil
	}
	switch tableFormat {
	case Markdown:
		return markdownHeaderOffset, nil
	case Backlog:
		return backlogHeaderOffset, nil
	default:
		return 0, fmt.Errorf("invalid table format detected")
	}
}

// getColor sets the color theme.
func getColor(tableTheme TableTheme) (*color.Color, error) {
	switch tableTheme {
	case DarkTheme:
		return color.New(color.BgHiBlack, color.FgHiWhite), nil
	case LightTheme:
		return color.New(color.BgHiWhite, color.FgHiBlack), nil
	case NoneTheme:
		return color.New(color.Reset), nil
	default:
		return nil, fmt.Errorf("invalid table theme detected")
	}
}
