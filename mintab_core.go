package mintab

import (
	"io"
	"strings"
)

const (
	// Default placeholder when field is empty in text table format.
	TextDefaultEmptyFieldPlaceholder = "-"

	// Default word delimiter in text table format.
	TextDefaultWordDelimiter = textNewLine

	// Default placeholder when field is empty in markdown table format.
	MarkdownDefaultEmptyFieldPlaceholder = "\\" + TextDefaultEmptyFieldPlaceholder

	// Default word delimiter in markdown table format.
	MarkdownDefaultWordDelimiter = markdownNewLine

	// Default placeholder when field is empty in backlog table format.
	BacklogDefaultEmptyFieldPlaceholder = TextDefaultEmptyFieldPlaceholder

	// Default word delimiter in backlog table format.
	BacklogDefaultWordDelimiter = backlogNewLine

	textNewLine     = "\n"
	markdownNewLine = "<br>"
	backlogNewLine  = "&br;"
)

// A Format represents the output format.
type Format int

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

// Input is a struct for loading values into Table
type Input struct {
	Header []string // table header
	Data   [][]any  // matrix with any types
}

// Table represents a table structure for rendering data.
type Table struct {
	w                     io.Writer       // Destination for table output
	b                     strings.Builder // Internal string builder
	data                  [][]string      // Matrix after parsing
	multilineData         [][][]string    // Matrix after parsing with each field divided by new lines
	header                []string        // Table header after parsing
	format                Format          // Output table format: text|compressed-text|markdown|backlog
	newLine               string          // New line string representation: "\n"|"<br>"|"&br;"
	emptyFieldPlaceholder string          // Placeholder for empty fields
	wordDelimiter         string          // Delimiter for words within a field
	lineHeights           []int           // Heights of lines with fields containing line breaks
	colWidths             []int           // Max widths of each columns
	numColumns            int             // Number of columns
	numColumnsFirstRow    int             // Number of columns of the first data row
	numRows               int             // Number of rows
	border                string          // Border line based on column widths
	tableWidth            int             // Table full width
	marginWidth           int             // Margin size around the field
	margin                string          // Whitespaces around the field
	hasHeader             bool            // Whether header rendering
	isEscape              bool            // Whether HTML escaping
	mergedFields          []int           // Indices of columns to merge
	ignoredFields         []int           // Indices of columns to ignore
	isMerge               bool            // Track whether to merge fields
	prevRow               []string        // Retain previous row
}

// New instantiates a new Table with the writer and options.
func New(w io.Writer, opts ...Option) *Table {
	var b strings.Builder
	t := &Table{
		w:                     w,
		b:                     b,
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
	if width < 0 {
		width = 0
	}
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
		t.isEscape = has
	}
}
