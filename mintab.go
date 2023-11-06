package mintab

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// Version
const Version = "0.0.7"

// Table format
const (
	MarkdownFormat = iota
	BacklogFormat
)

// Table theme
const (
	NoneTheme = iota
	DarkTheme
	LightTheme
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

// NewTable instantiates a table struct.
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
		return fmt.Errorf("cannot parse input: no elements in slice")
	}
	e := v.Index(0)
	if e.Kind() == reflect.Ptr {
		e = e.Elem()
	}
	if e.Kind() != reflect.Struct {
		return fmt.Errorf("cannot parse input: elements of slice must be struct or pointer to struct")
	}
	t.setColorFlags(v)
	t.setHeader(e.Type())
	if err = t.setData(v); err != nil {
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
func (t *Table) setHeader(typ reflect.Type) {
	if len(t.headers) > 0 {
		return
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if contains(t.ignoredFields, i) || field.PkgPath != "" {
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
			if contains(t.mergedFields, j) {
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

// setColorFlags determines which rows to color.
func (t *Table) setColorFlags(v reflect.Value) {
	if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Slice {
		v = v.Elem()
	}
	t.colorFlags = make([]bool, v.Len())
	colorFlag := true
	var prev any
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		firstField := item.Field(0).Interface()
		if i != 0 && !reflect.DeepEqual(prev, firstField) {
			colorFlag = !colorFlag
		}
		t.colorFlags[i] = colorFlag
		prev = firstField
	}
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
		if v.String() == "" {
			return t.emptyFieldPlaceholder, nil
		}
		return strings.TrimSpace(v.String()), nil
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
	return strings.ReplaceAll(s, "<br>", "&br;")
}

// colorize adds color to the table.
func (t *Table) colorize(table string) string {
	if t.theme == NoneTheme {
		return table
	}
	offset := t.getOffset()
	color := t.getColor()
	lines := strings.Split(table, "\n")
	m := make(map[int]struct{})
	for i, colorFlag := range t.colorFlags {
		if colorFlag {
			m[offset+i] = struct{}{}
		}
	}
	var clines []string
	for i, line := range lines {
		if _, ok := m[i]; ok {
			line = color.Sprint(line)
		}
		clines = append(clines, line)
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
	case NoneTheme:
		return &color.Color{}
	case DarkTheme:
		return color.New(color.BgHiBlack, color.FgHiWhite)
	case LightTheme:
		return color.New(color.BgHiWhite, color.FgHiBlack)
	default:
		return &color.Color{}
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
