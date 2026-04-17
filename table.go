package mintab

import (
	"io"
	"strings"
)

const (
	// TextDefaultPlaceholder is the default placeholder when a field is empty in text table format.
	TextDefaultPlaceholder = "-"

	// TextDefaultWordDelimiter is the default word delimiter in text table format.
	TextDefaultWordDelimiter = textNewLine

	// MarkdownDefaultPlaceholder is the default placeholder when a field is empty in markdown table format.
	MarkdownDefaultPlaceholder = "\\" + TextDefaultPlaceholder

	// MarkdownDefaultWordDelimiter is the default word delimiter in markdown table format.
	MarkdownDefaultWordDelimiter = markdownNewLine

	// BacklogDefaultPlaceholder is the default placeholder when a field is empty in backlog table format.
	BacklogDefaultPlaceholder = TextDefaultPlaceholder

	// BacklogDefaultWordDelimiter is the default word delimiter in backlog table format.
	BacklogDefaultWordDelimiter = backlogNewLine

	textNewLine     = "\n"
	markdownNewLine = "<br>"
	backlogNewLine  = "&br;"
)

// Input is a struct for loading values into Table.
type Input struct {
	Header []string // table header
	Data   [][]any  // matrix with any types
}

// Table represents a table structure for rendering data.
type Table struct {
	w                    io.Writer         // Destination for table output
	b                    strings.Builder   // Internal string builder
	r                    *strings.Replacer // Replacer for new lines in fields
	format               Format            // Output table format: text|compressed-text|markdown|backlog
	header               []string          // Table header after parsing
	data                 [][][]string      // Matrix after parsing with each field divided by new lines
	newLine              string            // New line string representation: "\n"|"<br>"|"&br;"
	placeholder          string            // Placeholder for empty fields
	wordDelimiter        string            // Delimiter for words within a field
	colWidths            []int             // Max widths of each columns
	lineHeights          []int             // Heights of lines with fields containing line breaks
	numColumns           int               // Number of columns
	numColumnsFirstRow   int               // Number of columns of the first data row
	numRows              int               // Number of rows
	border               string            // Border line based on column widths
	tableWidth           int               // Table full width
	marginWidth          int               // Margin size around the field
	marginWidthBothSides int               // Twice of margin size
	margin               string            // Whitespaces around the field
	hasHeader            bool              // Whether header rendering
	isEscape             bool              // Whether HTML escaping (mainly designed for markdown)
	isMerge              bool              // Track whether to merge fields
	isBytesToString      bool              // Whether []uint8 should be treated as string
	prevRow              []string          // Retain previous row
	mergedFields         []int             // Indices of columns to merge
	ignoredFields        []int             // Indices of columns to ignore
}

// New instantiates a new Table with the writer and options.
func New(w io.Writer, opts ...Option) *Table {
	t := &Table{
		w:                    w,
		format:               TextFormat,
		newLine:              textNewLine,
		placeholder:          TextDefaultPlaceholder,
		wordDelimiter:        TextDefaultWordDelimiter,
		marginWidth:          1,
		marginWidthBothSides: 2,
		hasHeader:            true,
		isBytesToString:      true,
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
	if width < 0 {
		width = 0
	}
	return func(t *Table) {
		t.marginWidth = width
		t.marginWidthBothSides = width * 2
	}
}

// WithPlaceholder sets the placeholder for empty fields.
func WithPlaceholder(placeholder string) Option {
	if placeholder == "" {
		placeholder = " "
	}
	return func(t *Table) {
		t.placeholder = placeholder
	}
}

// WithWordDelimiter sets the delimiter to split words in a field.
func WithWordDelimiter(delimiter string) Option {
	return func(t *Table) {
		t.wordDelimiter = delimiter
	}
}

// WithMergeFields sets column indices to be merged.
func WithMergeFields(indices []int) Option {
	return func(t *Table) {
		t.mergedFields = indices
	}
}

// WithIgnoreFields sets column indices to be ignored.
func WithIgnoreFields(indices []int) Option {
	return func(t *Table) {
		t.ignoredFields = indices
	}
}

// WithEscape enables or disables HTML escaping.
func WithEscape(has bool) Option {
	return func(t *Table) {
		t.isEscape = has
	}
}

// WithBytesAsString controls how []uint8 is interpreted.
func WithBytesAsString(has bool) Option {
	return func(t *Table) {
		t.isBytesToString = has
	}
}
