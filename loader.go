package mintab

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
)

// Load validates v and converts it to a struct Table. v must be passed in one of the following two ways:
//
// 1. Struct `Input`
//   - The number of columns in all rows must be the same.
//   - Header is allowd to be nil.
//
// 2. Any struct slices
//   - If a struct is passed, it is converted to a slice with one element.
//   - If the field is a slice with primitive data type or a slice of byte slice, it is converted to a string.
//   - If the field is struct, an error is returned (nested structs are not supported)
func (t *Table) Load(v any) error {
	if _, ok := v.([]any); ok {
		return fmt.Errorf("cannot load input: elements of slice must not be any")
	}
	switch tv := v.(type) {
	case nil:
		return nil
	case Input:
		if err := t.loadInput(tv); err != nil {
			return err
		}
	case *Input:
		if err := t.loadInput(*tv); err != nil {
			return err
		}
	default:
		if err := t.loadStruct(tv); err != nil {
			return err
		}
	}
	return nil
}

func (t *Table) loadInput(v Input) error {
	t.numRows = len(v.Data)
	if t.numRows == 0 {
		return nil
	}
	t.setFormat()
	if err := t.setInputHeader(v); err != nil {
		return err
	}
	if err := t.setInputData(v); err != nil {
		return err
	}
	t.setBorder()
	return nil
}

func (t *Table) loadStruct(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		if rv.IsZero() {
			return nil
		}
		rv = reflect.Append(reflect.MakeSlice(reflect.SliceOf(rv.Type()), 0, 1), rv)
	}
	t.numRows = rv.Len()
	if t.numRows == 0 {
		return nil
	}
	t.setFormat()
	if err := t.setStructHeader(rv); err != nil {
		return err
	}
	if err := t.setStructData(rv); err != nil {
		return err
	}
	t.setBorder()
	return nil
}

func (t *Table) setFormat() {
	var p, d string
	switch t.format {
	case MarkdownFormat:
		p = MarkdownDefaultPlaceholder
		d = MarkdownDefaultWordDelimiter
		t.newLine = markdownNewLine
	case BacklogFormat:
		p = BacklogDefaultPlaceholder
		d = BacklogDefaultWordDelimiter
		t.newLine = backlogNewLine
	default:
		p = TextDefaultPlaceholder
		d = TextDefaultWordDelimiter
	}
	if t.placeholder == TextDefaultPlaceholder {
		t.placeholder = p
	}
	if t.wordDelimiter == TextDefaultWordDelimiter {
		t.wordDelimiter = d
	}
}

func (t *Table) setInputHeader(v Input) error {
	t.numColumns = len(v.Header)
	firstRow := v.Data[0]
	t.numColumnsFirstRow = len(firstRow)
	if t.numColumns > 0 {
		if t.numColumns != t.numColumnsFirstRow {
			return fmt.Errorf("cannot load input: number of columns must be the same as header")
		}
		t.header = make([]string, 0, t.numColumns)
	}
	t.colWidths = make([]int, 0, t.numColumns)
	for i, h := range v.Header {
		if !slices.Contains(t.ignoredFields, i) {
			t.header = append(t.header, h)
			t.colWidths = append(t.colWidths, runewidth.StringWidth(h))
		}
	}
	t.numColumns = len(t.colWidths)
	if t.numColumns == 0 {
		for i := range firstRow {
			if !slices.Contains(t.ignoredFields, i) {
				t.colWidths = append(t.colWidths, 0)
			}
		}
	}
	return nil
}

func (t *Table) setStructHeader(rv reflect.Value) error {
	e := rv.Index(0)
	if e.Kind() == reflect.Ptr {
		e = e.Elem()
	}
	if e.Kind() != reflect.Struct {
		return fmt.Errorf("cannot load input: elements of slice must be struct or pointer to struct")
	}
	typ := e.Type()
	t.numColumns = typ.NumField()
	t.header = make([]string, 0, t.numColumns)
	t.colWidths = make([]int, 0, t.numColumns)
	for i := 0; i < t.numColumns; i++ {
		field := typ.Field(i)
		if !slices.Contains(t.ignoredFields, i) && field.PkgPath == "" {
			t.header = append(t.header, field.Name)
			t.colWidths = append(t.colWidths, runewidth.StringWidth(field.Name))
		}
	}
	t.numColumns = len(t.colWidths)
	if t.numColumns == 0 {
		return fmt.Errorf("cannot load input: at least one exported field is required")
	}
	return nil
}

func (t *Table) setInputData(v Input) error {
	t.data = make([][][]string, t.numRows)
	t.lineHeights = make([]int, t.numRows)
	n := t.numColumns
	if n == 0 {
		n = t.numColumnsFirstRow
	}
	t.prevRow = make([]string, n)
	for i, r := range v.Data {
		if i > 0 && len(r) != t.numColumnsFirstRow {
			return fmt.Errorf("cannot load input: number of columns must be the same for all rows")
		}
		row := make([][]string, n)
		t.isMerge = true
		t.lineHeights[i] = 1
		for j, field := range r {
			if slices.Contains(t.ignoredFields, j) {
				continue
			}
			s, err := t.formatField(reflect.ValueOf(field))
			if err != nil {
				return err
			}
			s = t.merge(s, j)
			elems := strings.Split(s, "\n")
			row[j] = elems
			t.updateColWidths(elems, j)
			t.getLineHeight(elems, i)
		}
		t.data[i] = row
	}
	return nil
}

func (t *Table) setStructData(rv reflect.Value) error {
	t.data = make([][][]string, t.numRows)
	t.lineHeights = make([]int, t.numRows)
	t.prevRow = make([]string, t.numColumns)
	for i := 0; i < t.numRows; i++ {
		e := rv.Index(i)
		if e.Kind() == reflect.Ptr {
			e = e.Elem()
		}
		row := make([][]string, t.numColumns)
		t.isMerge = true
		t.lineHeights[i] = 1
		for j, h := range t.header {
			field := e.FieldByName(h)
			if !field.IsValid() {
				return fmt.Errorf("cannot load input: invalid field detected: %s", h)
			}
			s, err := t.formatField(field)
			if err != nil {
				return err
			}
			s = t.merge(s, j)
			elems := strings.Split(s, "\n")
			row[j] = elems
			t.updateColWidths(elems, j)
			t.getLineHeight(elems, i)
		}
		t.data[i] = row
	}
	return nil
}

func (t *Table) merge(s string, i int) string {
	if slices.Contains(t.mergedFields, i) {
		if s != t.prevRow[i] {
			t.isMerge = false
			t.prevRow[i] = s
		}
		if t.isMerge {
			s = ""
		}
	}
	return s
}

func (t *Table) updateColWidths(elems []string, i int) {
	for _, elem := range elems {
		w := runewidth.StringWidth(elem)
		if w > t.colWidths[i] {
			t.colWidths[i] = w
		}
	}
}

func (t *Table) getLineHeight(elems []string, i int) {
	switch t.format {
	case TextFormat, CompressedTextFormat:
		height := len(elems)
		if height > t.lineHeights[i] {
			t.lineHeights[i] = height
		}
	}
}

func (t *Table) setBorder() {
	var sep string
	switch t.format {
	case MarkdownFormat, BacklogFormat:
		sep = "|"
	default:
		sep = "+"
	}
	t.b.Reset()
	t.b.Grow(128)
	for _, w := range t.colWidths {
		t.b.WriteString(sep)
		for i := 0; i < w+t.marginWidthBothSides; i++ {
			t.b.WriteByte('-')
		}
	}
	t.b.WriteString(sep)
	t.b.WriteString("\n")
	t.border = t.b.String()
	t.tableWidth = len(t.border)
}

func (t *Table) formatField(rv reflect.Value) (string, error) {
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return t.placeholder, nil
		}
		rv = rv.Elem()
	}
	if s := getStringer(rv); s != "" {
		return t.sanitize(s), nil
	}
	switch rv.Kind() {
	case reflect.String:
		return t.sanitize(rv.String()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10), nil
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 64), nil
	case reflect.Struct:
		return "", fmt.Errorf("cannot load input: nested fields not supported")
	case reflect.Slice, reflect.Array:
		s, err := t.formatSlice(rv)
		if err != nil {
			return "", err
		}
		return t.sanitize(s), nil
	default:
		return t.sanitize(fmt.Sprint(rv.Interface())), nil
	}
}

func (t *Table) formatSlice(rv reflect.Value) (string, error) {
	l := rv.Len()
	switch {
	case l == 0:
		return t.placeholder, nil
	case rv.Type().Elem().Kind() == reflect.Uint8:
		return string(rv.Bytes()), nil
	default:
		t.b.Reset()
		for i := 0; i < l; i++ {
			e := rv.Index(i)
			if i != 0 {
				t.b.WriteString(t.wordDelimiter)
			}
			if e.Kind() == reflect.Ptr {
				if e.IsNil() {
					t.b.WriteString(t.placeholder)
					continue
				}
				e = e.Elem()
			}
			if s := getStringer(e); s != "" {
				t.b.WriteString(s)
				continue
			}
			if e.Kind() == reflect.Slice && e.Type().Elem().Kind() == reflect.Uint8 {
				t.b.WriteString(string(e.Bytes()))
				continue
			}
			if e.Kind() == reflect.Slice || e.Kind() == reflect.Array || e.Kind() == reflect.Struct {
				return "", fmt.Errorf("cannot load input: nested fields not supported")
			}
			switch v := e.Interface().(type) {
			case string:
				if v == "" {
					t.b.WriteString(t.placeholder)
				} else {
					t.b.WriteString(v)
				}
			case int, int8, int16, int32, int64:
				t.b.WriteString(strconv.FormatInt(v.(int64), 10))
			case uint, uint8, uint16, uint32, uint64:
				t.b.WriteString(strconv.FormatUint(v.(uint64), 10))
			case float32:
				t.b.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 32))
			case float64:
				t.b.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
			default:
				t.b.WriteString(fmt.Sprint(v))
			}
		}
		return t.b.String(), nil
	}
}

func getStringer(rv reflect.Value) string {
	if s, ok := rv.Interface().(fmt.Stringer); ok {
		return s.String()
	}
	return ""
}

func (t *Table) sanitize(s string) string {
	if s == "" {
		return t.placeholder
	}
	if t.isEscape {
		s = t.escape(s)
	}
	if t.format == MarkdownFormat && strings.HasPrefix(s, "*") {
		s = "\\" + s
	}
	if t.format == TextFormat {
		return s
	}
	t.b.Reset()
	t.b.Grow(len(s) + len(t.newLine)*strings.Count(s, "\n"))
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			t.b.WriteString(s[start:i])
			t.b.WriteString(t.newLine)
			start = i + 1
		}
	}
	t.b.WriteString(s[start:])
	return strings.TrimSpace(t.b.String())
}

func (t *Table) escape(s string) string {
	t.b.Reset()
	for _, r := range s {
		switch r {
		case '<':
			t.b.WriteString("&lt;")
		case '>':
			t.b.WriteString("&gt;")
		case '"':
			t.b.WriteString("&quot;")
		case '\'':
			t.b.WriteString("&lsquo;")
		case '&':
			t.b.WriteString("&amp;")
		case ' ':
			t.b.WriteString("&nbsp;")
		case '*':
			t.b.WriteString("&#42;")
		case '\\':
			t.b.WriteString("&#92;")
		case '_':
			t.b.WriteString("&#95;")
		case '|':
			t.b.WriteString("&#124;")
		default:
			t.b.WriteRune(r)
		}
	}
	return t.b.String()
}
