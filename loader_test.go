package mintab

import (
	"net"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestTable_Load(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "input",
			args: args{
				v: basicTestInput,
			},
			wantErr: false,
		},
		{
			name: "input_ptr",
			args: args{
				v: basicTestInputPtr,
			},
			wantErr: false,
		},
		{
			name: "input_ptr_error",
			args: args{
				v: &nestedTestInput1,
			},
			wantErr: true,
		},
		{
			name: "input_merged",
			args: args{
				v: mergedTestInput,
			},
			wantErr: false,
		},
		{
			name: "input_escaped",
			args: args{
				v: escapedTestInput,
			},
			wantErr: false,
		},
		{
			name: "input_noheader",
			args: args{
				v: noHeaderTestInput,
			},
			wantErr: false,
		},
		{
			name: "input_nested_1",
			args: args{
				v: nestedTestInput1,
			},
			wantErr: true,
		},
		{
			name: "input_nested_2",
			args: args{
				v: nestedTestInput2,
			},
			wantErr: true,
		},
		{
			name: "input_invalid_header_indices",
			args: args{
				v: invalidHeaderIndicesTestInput,
			},
			wantErr: true,
		},
		{
			name: "input_invalid_data_indices",
			args: args{
				v: invalidDataIndicesTestInput,
			},
			wantErr: true,
		},
		{
			name: "input_empty",
			args: args{
				v: emptyTestInput,
			},
			wantErr: false,
		},
		{
			name: "struct_slice",
			args: args{
				v: basicTestStructSlice,
			},
			wantErr: false,
		},
		{
			name: "struct_slice_empty",
			args: args{
				v: basicTestStructSliceEmpty,
			},
			wantErr: false,
		},
		{
			name: "struct_non_slice",
			args: args{
				v: basicTestStructNonSlice,
			},
			wantErr: false,
		},
		{
			name: "struct_non_slice_empty",
			args: args{
				v: basicTestStructNonSliceEmpty,
			},
			wantErr: false,
		},
		{
			name: "struct_ptr_slice",
			args: args{
				v: basicTestStructPtrSlice,
			},
			wantErr: false,
		},
		{
			name: "struct_slice_ptr",
			args: args{
				v: basicTestStructSlicePtr,
			},
			wantErr: false,
		},
		{
			name: "struct_slice_merged",
			args: args{
				v: basicTestStructSliceEmpty,
			},
			wantErr: false,
		},
		{
			name: "struct_slice_escaped",
			args: args{
				v: escapedTestStructSlice,
			},
			wantErr: false,
		},
		{
			name: "struct_slice_nested",
			args: args{
				v: nestedTestStructSlice,
			},
			wantErr: true,
		},
		{
			name: "struct_slice_stringer",
			args: args{
				v: stringerTestStructSlice,
			},
			wantErr: false,
		},
		{
			name: "struct_slice_non_exported",
			args: args{
				v: nonExportedTestStructSlice,
			},
			wantErr: true,
		},
		{
			name: "struct_slice_non_type",
			args: args{
				v: nonTypeTestStructSlice,
			},
			wantErr: true,
		},
		{
			name: "nil",
			args: args{
				v: nil,
			},
			wantErr: false,
		},
		{
			name: "string",
			args: args{
				v: "aaa",
			},
			wantErr: true,
		},
		{
			name: "int",
			args: args{
				v: 1,
			},
			wantErr: true,
		},
		{
			name: "string_slice",
			args: args{
				v: []string{"dummy"},
			},
			wantErr: true,
		},
		{
			name: "int_slice",
			args: args{
				v: []string{"dummy"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Table{}
			if err := tr.Load(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Table.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTable_setFormat(t *testing.T) {
	type fields struct {
		format Format
	}
	type want struct {
		emptyFieldPlaceholder string
		wordDelimiter         string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "text",
			fields: fields{
				format: TextFormat,
			},
			want: want{
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
			},
		},
		{
			name: "markdown",
			fields: fields{
				format: MarkdownFormat,
			},
			want: want{
				emptyFieldPlaceholder: MarkdownDefaultEmptyFieldPlaceholder,
				wordDelimiter:         MarkdownDefaultWordDelimiter,
			},
		},
		{
			name: "backlog",
			fields: fields{
				format: BacklogFormat,
			},
			want: want{
				emptyFieldPlaceholder: BacklogDefaultEmptyFieldPlaceholder,
				wordDelimiter:         BacklogDefaultWordDelimiter,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := New(os.Stdout, WithFormat(tt.fields.format))
			tr.setFormat()
			if !reflect.DeepEqual(tr.emptyFieldPlaceholder, tt.want.emptyFieldPlaceholder) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", tr.emptyFieldPlaceholder, tt.want.emptyFieldPlaceholder)
			}
			if !reflect.DeepEqual(tr.wordDelimiter, tt.want.wordDelimiter) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", tr.wordDelimiter, tt.want.wordDelimiter)
			}
		})
	}
}

func TestTable_setStructHeader(t *testing.T) {
	type fields struct {
		ignoredFields []int
	}
	type args struct {
		rv reflect.Value
	}
	type want struct {
		header    []string
		colWidths []int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				ignoredFields: nil,
			},
			args: args{
				rv: reflect.ValueOf(basicTestStructSlice),
			},
			want: want{
				header:    []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				colWidths: []int{10, 12, 10, 10},
			},
			wantErr: false,
		},
		{
			name: "ignore",
			fields: fields{
				ignoredFields: []int{1},
			},
			args: args{
				rv: reflect.ValueOf(basicTestStructSlice),
			},
			want: want{
				header:    []string{"InstanceID", "AttachedLB", "AttachedTG"},
				colWidths: []int{10, 10, 10},
			},
			wantErr: false,
		},
		{
			name: "ignore-signed-int",
			fields: fields{
				ignoredFields: []int{-10},
			},
			args: args{
				rv: reflect.ValueOf(basicTestStructSlice),
			},
			want: want{
				header:    []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				colWidths: []int{10, 12, 10, 10},
			},
			wantErr: false,
		},
		{
			name: "non-exported",
			fields: fields{
				ignoredFields: nil,
			},
			args: args{
				rv: reflect.ValueOf(nonExportedTestStructSlice),
			},
			want: want{
				header:    []string{},
				colWidths: []int{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Table{
				ignoredFields: tt.fields.ignoredFields,
			}
			if err := tr.setStructHeader(tt.args.rv); (err != nil) != tt.wantErr {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tr.header, tt.want.header) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", tr.header, tt.want.header)
			}
			if !reflect.DeepEqual(tr.colWidths, tt.want.colWidths) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", tr.colWidths, tt.want.colWidths)
			}
		})
	}
}

func TestTable_setInputHeader(t *testing.T) {
	type fields struct {
		ignoredFields []int
	}
	type args struct {
		rv Input
	}
	type want struct {
		header    []string
		colWidths []int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				ignoredFields: nil,
			},
			args: args{
				rv: basicTestInput,
			},
			want: want{
				header:    []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				colWidths: []int{10, 12, 10, 10},
			},
			wantErr: false,
		},
		{
			name: "ignore",
			fields: fields{
				ignoredFields: []int{1},
			},
			args: args{
				rv: basicTestInput,
			},
			want: want{
				header:    []string{"InstanceID", "AttachedLB", "AttachedTG"},
				colWidths: []int{10, 10, 10},
			},
			wantErr: false,
		},
		{
			name: "ignore-signed-int",
			fields: fields{
				ignoredFields: []int{-10},
			},
			args: args{
				rv: basicTestInput,
			},
			want: want{
				header:    []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				colWidths: []int{10, 12, 10, 10},
			},
			wantErr: false,
		},
		{
			name: "irregular-columns-number",
			fields: fields{
				ignoredFields: nil,
			},
			args: args{
				rv: invalidHeaderIndicesTestInput,
			},
			want: want{
				header:    nil,
				colWidths: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Table{
				ignoredFields: tt.fields.ignoredFields,
			}
			if err := tr.setInputHeader(tt.args.rv); (err != nil) != tt.wantErr {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tr.header, tt.want.header) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", tr.header, tt.want.header)
			}
			if !reflect.DeepEqual(tr.colWidths, tt.want.colWidths) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", tr.colWidths, tt.want.colWidths)
			}
		})
	}
}

func TestTable_setBorder(t *testing.T) {
	type fields struct {
		format               Format
		marginWidth          int
		marginWidthBothSides int
		colWidths            []int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "text",
			fields: fields{
				format:               TextFormat,
				marginWidth:          1,
				marginWidthBothSides: 2,
				colWidths:            []int{8, 12, 5},
			},
			want: "+----------+--------------+-------+\n",
		},
		{
			name: "markdown",
			fields: fields{
				format:               MarkdownFormat,
				marginWidth:          1,
				marginWidthBothSides: 2,
				colWidths:            []int{8, 12, 5},
			},
			want: "|----------|--------------|-------|\n",
		},
		{
			name: "backlog",
			fields: fields{
				format:               BacklogFormat,
				marginWidth:          1,
				marginWidthBothSides: 2,
				colWidths:            []int{8, 12, 5},
			},
			want: "|----------|--------------|-------|\n",
		},
		{
			name: "wide-margin",
			fields: fields{
				format:               TextFormat,
				marginWidth:          3,
				marginWidthBothSides: 6,
				colWidths:            []int{8, 12, 5},
			},
			want: "+--------------+------------------+-----------+\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Table{
				format:               tt.fields.format,
				marginWidth:          tt.fields.marginWidth,
				marginWidthBothSides: tt.fields.marginWidthBothSides,
				colWidths:            tt.fields.colWidths,
			}
			tr.setBorder()
			if !reflect.DeepEqual(tr.border, tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", tr.border, tt.want)
			}
		})
	}
}

func TestTable_formatField(t *testing.T) {
	sp := func(s string) *string {
		return &s
	}
	type fields struct {
		format                Format
		emptyFieldPlaceholder string
		wordDelimiter         string
		isEscape              bool
	}
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "string",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: "aaa",
			},
			want:    "aaa",
			wantErr: false,
		},
		{
			name: "string_empty",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: "",
			},
			want:    TextDefaultEmptyFieldPlaceholder,
			wantErr: false,
		},
		{
			name: "byte_slice",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: []byte("aaa"),
			},
			want:    "aaa",
			wantErr: false,
		},
		{
			name: "escape",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              true,
			},
			args: args{
				v: `<>"'& *\_|`,
			},
			want:    "&lt;&gt;&quot;&lsquo;&amp;&nbsp;&#42;&#92;&#95;&#124;",
			wantErr: false,
		},
		{
			name: "asterisk_prefix_at_markdown",
			fields: fields{
				format:                MarkdownFormat,
				emptyFieldPlaceholder: MarkdownDefaultEmptyFieldPlaceholder,
				wordDelimiter:         MarkdownDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: "*.example.com",
			},
			want:    "\\*.example.com",
			wantErr: false,
		},
		{
			name: "int",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: 123,
			},
			want:    "123",
			wantErr: false,
		},
		{
			name: "int_signed",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: -123,
			},
			want:    "-123",
			wantErr: false,
		},
		{
			name: "uint",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: uint(123),
			},
			want:    "123",
			wantErr: false,
		},
		{
			name: "float",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: 123.456,
			},
			want:    "123.456",
			wantErr: false,
		},
		{
			name: "float32",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: float32(1.5),
			},
			want:    "1.5",
			wantErr: false,
		},
		{
			name: "ptr",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: sp("aaa"),
			},
			want:    "aaa",
			wantErr: false,
		},
		{
			name: "nil_ptr",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: (*string)(nil),
			},
			want:    TextDefaultEmptyFieldPlaceholder,
			wantErr: false,
		},
		{
			name: "non_nil_ptr_string",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: new(string),
			},
			want:    TextDefaultEmptyFieldPlaceholder,
			wantErr: false,
		},
		{
			name: "non_nil_ptr_int",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: new(int),
			},
			want:    "0",
			wantErr: false,
		},
		{
			name: "slice_string",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: []string{"a", "b"},
			},
			want:    "a" + TextDefaultWordDelimiter + "b",
			wantErr: false,
		},
		{
			name: "slice_string_included_empty",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: []string{"a", "", "b"},
			},
			want:    "a" + TextDefaultWordDelimiter + TextDefaultEmptyFieldPlaceholder + TextDefaultWordDelimiter + "b",
			wantErr: false,
		},
		{
			name: "slice_int",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: []int{0, 1, 2},
			},
			want:    "0" + TextDefaultWordDelimiter + "1" + TextDefaultWordDelimiter + "2",
			wantErr: false,
		},
		{
			name: "slice_uint",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: []uint{0, 1, 2},
			},
			want:    "0" + TextDefaultWordDelimiter + "1" + TextDefaultWordDelimiter + "2",
			wantErr: false,
		},
		{
			name: "slice_float32",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: []float32{0.1, 1.25, 2.001},
			},
			want:    "0.1" + TextDefaultWordDelimiter + "1.25" + TextDefaultWordDelimiter + "2.001",
			wantErr: false,
		},
		{
			name: "slice_float64",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: []float64{0.1, 1.25, 2.001},
			},
			want:    "0.1" + TextDefaultWordDelimiter + "1.25" + TextDefaultWordDelimiter + "2.001",
			wantErr: false,
		},
		{
			name: "slice_nil",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: ([]string)(nil),
			},
			want:    TextDefaultEmptyFieldPlaceholder,
			wantErr: false,
		},
		{
			name: "slice_empty",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: []string{},
			},
			want:    TextDefaultEmptyFieldPlaceholder,
			wantErr: false,
		},
		{
			name: "slice_with_byte_slice",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: []byte("aaa"),
			},
			want:    "aaa",
			wantErr: false,
		},
		{
			name: "slice_slice",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: [][]string{
					{"a", "b", "c"},
					{"x", "y", "z"},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "slice_struct",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: []struct {
					key    string
					values []int
				}{
					{key: "key1", values: []int{1, 2, 3}},
					{key: "key2", values: []int{4, 5, 6}},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "slice_ptr",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: &[]string{"a", "b"},
			},
			want:    "a" + TextDefaultWordDelimiter + "b",
			wantErr: false,
		},
		{
			name: "slice_with_ptr_to_strings",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: []*string{sp(""), sp("a"), sp("b")},
			},
			want:    TextDefaultEmptyFieldPlaceholder + TextDefaultWordDelimiter + "a" + TextDefaultWordDelimiter + "b",
			wantErr: false,
		},
		{
			name: "slice_with_ptr_to_string_empty",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: []*string{},
			},
			want:    TextDefaultEmptyFieldPlaceholder,
			wantErr: false,
		},
		{
			name: "slice_with_nil_ptr",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: []*int{nil},
			},
			want:    TextDefaultEmptyFieldPlaceholder,
			wantErr: false,
		},
		{
			name: "slice_with_ptr_mixed",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: []*string{nil, sp(""), sp("aaa")},
			},
			want:    TextDefaultEmptyFieldPlaceholder + TextDefaultWordDelimiter + TextDefaultEmptyFieldPlaceholder + TextDefaultWordDelimiter + "aaa",
			wantErr: false,
		},
		{
			name: "stringer_duration",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: 123 * time.Hour,
			},
			want:    "123h0m0s",
			wantErr: false,
		},
		{
			name: "stringer_ipaddress",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: net.IPv4bcast,
			},
			want:    "255.255.255.255",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.args.v)
			tr := &Table{
				format:                tt.fields.format,
				emptyFieldPlaceholder: tt.fields.emptyFieldPlaceholder,
				wordDelimiter:         tt.fields.wordDelimiter,
				isEscape:              tt.fields.isEscape,
			}
			got, err := tr.formatField(v)
			if (err != nil) != tt.wantErr {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, tt.want)
			}
		})
	}
}

func TestTable_sanitize(t *testing.T) {
	type fields struct {
		format        Format
		wordDelimiter string
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "text",
			fields: fields{
				format:        TextFormat,
				wordDelimiter: "\n",
			},
			args: args{
				s: "aaa\nbbb\nccc",
			},
			want: "aaa\nbbb\nccc",
		},
		{
			name: "markdown",
			fields: fields{
				format:        MarkdownFormat,
				wordDelimiter: "<br>",
			},
			args: args{
				s: "aaa\nbbb\nccc",
			},
			want: "aaa<br>bbb<br>ccc",
		},
		{
			name: "backlog",
			fields: fields{
				format:        BacklogFormat,
				wordDelimiter: "&br;",
			},
			args: args{
				s: "aaa\nbbb\nccc",
			},
			want: "aaa&br;bbb&br;ccc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Table{
				format:        tt.fields.format,
				wordDelimiter: tt.fields.wordDelimiter,
			}
			tr.setFormat()
			got := tr.sanitize(tt.args.s)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, tt.want)
			}
		})
	}
}

func TestTable_escape(t *testing.T) {
	type fields struct {
		builder strings.Builder
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "wildcard",
			fields: fields{
				builder: strings.Builder{},
			},
			args: args{
				s: "*.example.com",
			},
			want: "&#42;.example.com",
		},
		{
			name: "html",
			fields: fields{
				builder: strings.Builder{},
			},
			args: args{
				s: "<span style=\"color:#d70910;\">red</span>",
			},
			want: "&lt;span&nbsp;style=&quot;color:#d70910;&quot;&gt;red&lt;/span&gt;",
		},
		{
			name: "json",
			fields: fields{
				builder: strings.Builder{},
			},
			args: args{
				s: `{
  "key": [
    "value1",
    "value2",
    "value3",
  ]
}`,
			},
			want: `{
&nbsp;&nbsp;&quot;key&quot;:&nbsp;[
&nbsp;&nbsp;&nbsp;&nbsp;&quot;value1&quot;,
&nbsp;&nbsp;&nbsp;&nbsp;&quot;value2&quot;,
&nbsp;&nbsp;&nbsp;&nbsp;&quot;value3&quot;,
&nbsp;&nbsp;]
}`,
		},
		{
			name: "other",
			fields: fields{
				builder: strings.Builder{},
			},
			args: args{
				s: `'&\_|`,
			},
			want: "&lsquo;&amp;&#92;&#95;&#124;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Table{
				//	builder: tt.fields.builder,
			}
			got := tr.escape(tt.args.s)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, tt.want)
			}
		})
	}
}
