package mintab

import (
	"reflect"
	"testing"
)

func TestFormat_String(t *testing.T) {
	tests := []struct {
		name string
		o    Format
		want string
	}{
		{
			name: "text",
			o:    TextFormat,
			want: "text",
		},
		{
			name: "compressed",
			o:    CompressedTextFormat,
			want: "compressed",
		},
		{
			name: "markdown",
			o:    MarkdownFormat,
			want: "markdown",
		},
		{
			name: "backlog",
			o:    BacklogFormat,
			want: "backlog",
		},
		{
			name: "other",
			o:    9,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.String(); got != tt.want {
				t.Errorf("Format.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormat_Values(t *testing.T) {
	tests := []struct {
		name string
		tr   Format
		want []Format
	}{
		{
			name: "all formats",
			tr:   TextFormat,
			want: []Format{TextFormat, CompressedTextFormat, MarkdownFormat, BacklogFormat},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Values(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Format.Values() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFormat(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Format
		wantErr bool
	}{
		{
			name:    "parse text",
			args:    args{s: "text"},
			want:    TextFormat,
			wantErr: false,
		},
		{
			name:    "parse compressed",
			args:    args{s: "compressed"},
			want:    CompressedTextFormat,
			wantErr: false,
		},
		{
			name:    "parse markdown",
			args:    args{s: "markdown"},
			want:    MarkdownFormat,
			wantErr: false,
		},
		{
			name:    "parse backlog",
			args:    args{s: "backlog"},
			want:    BacklogFormat,
			wantErr: false,
		},
		{
			name:    "invalid format",
			args:    args{s: "invalid"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty string",
			args:    args{s: ""},
			want:    0,
			wantErr: true,
		},
		{
			name:    "case sensitivity",
			args:    args{s: "Text"},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFormat(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
