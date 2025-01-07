package mintab

import (
	"encoding/json"
	"fmt"
)

// A Format represents the output format.
type Format int

const (
	TextFormat           Format = iota // Text table format.
	CompressedTextFormat               // Compressed text table format.
	MarkdownFormat                     // Markdown table format.
	BacklogFormat                      // Backlog-specific table format.
)

// MarshalJSON marshals a Format into JSON.
func (t Format) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// String returns the string representation of a Format.
func (t Format) String() string {
	switch t {
	case TextFormat:
		return "text"
	case CompressedTextFormat:
		return "compressed"
	case MarkdownFormat:
		return "markdown"
	case BacklogFormat:
		return "backlog"
	default:
		return ""
	}
}

// Values returns the list of supported formats.
func (t Format) Values() []Format {
	return []Format{
		TextFormat,
		CompressedTextFormat,
		MarkdownFormat,
		BacklogFormat,
	}
}

// ParseFormat parses a string into a Format.
func ParseFormat(s string) (Format, error) {
	switch s {
	case TextFormat.String():
		return TextFormat, nil
	case CompressedTextFormat.String():
		return CompressedTextFormat, nil
	case MarkdownFormat.String():
		return MarkdownFormat, nil
	case BacklogFormat.String():
		return BacklogFormat, nil
	default:
		return 0, fmt.Errorf("unsupported format: %q", s)
	}
}
