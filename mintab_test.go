package mintab

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/fatih/color"
)

type sample struct {
	InstanceName      string
	SecurityGroupName string
	CidrBlock         []string
}

type object struct {
	ObjectID   int
	ObjectName string
}

type nested struct {
	BucketName string
	Objects    []object
}

var (
	samples                      []sample
	nests                        []nested
	defaultEmptyFieldPlaceholder string
	defaultWordDelimiter         string
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	defaultEmptyFieldPlaceholder = "N/A"
	defaultWordDelimiter = "<br>"
	samples = []sample{
		{InstanceName: "i-1", SecurityGroupName: "sg-1", CidrBlock: []string{"10.0.0.0/16"}},
		{InstanceName: "i-1", SecurityGroupName: "sg-1", CidrBlock: []string{"10.1.0.0/16"}},
		{InstanceName: "i-1", SecurityGroupName: "sg-2", CidrBlock: []string{"10.2.0.0/16"}},
		{InstanceName: "i-1", SecurityGroupName: "sg-2", CidrBlock: []string{"10.3.0.0/16"}},
		{InstanceName: "i-2", SecurityGroupName: "sg-1", CidrBlock: []string{"10.0.0.0/16", "0.0.0.0/0"}},
		{InstanceName: "i-2", SecurityGroupName: "sg-1", CidrBlock: []string{"10.1.0.0/16", "0.0.0.0/0"}},
		{InstanceName: "i-2", SecurityGroupName: "sg-2", CidrBlock: []string{"10.2.0.0/16", "0.0.0.0/0"}},
		{InstanceName: "i-2", SecurityGroupName: "sg-2", CidrBlock: []string{"10.3.0.0/16", "0.0.0.0/0"}},
		{InstanceName: "i-3", SecurityGroupName: "", CidrBlock: []string{"10.0.0.0/16", "0.0.0.0/0"}},
		{InstanceName: "i-4", SecurityGroupName: "sg-4", CidrBlock: []string{}},
	}
	nests = []nested{
		{
			BucketName: "bucket1",
			Objects: []object{
				{
					ObjectID:   11,
					ObjectName: "bucket1-obj1",
				},
				{
					ObjectID:   12,
					ObjectName: "bucket1-obj2",
				},
			},
		},
		{
			BucketName: "bucket2",
			Objects: []object{
				{
					ObjectID:   21,
					ObjectName: "bucket2-obj1",
				},
				{
					ObjectID:   22,
					ObjectName: "bucket2-obj2",
				},
				{
					ObjectID:   23,
					ObjectName: "bucket2-obj3",
				},
			},
		},
	}
}

func TestNewTable(t *testing.T) {
	type args struct {
		format                int
		theme                 int
		hasHeader             bool
		emptyFieldPlaceholder string
		wordDelimiter         string
		mergeFields           []int
		ignoreFields          []int
	}
	type want struct {
		got *Table
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "markdown+basic",
			args: args{
				format:                MarkdownFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                MarkdownFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergeFields:           nil,
					ignoreFields:          nil,
					colorFlags:            nil,
				},
			},
		},
		{
			name: "markdown+disableHeader",
			args: args{
				format:                MarkdownFormat,
				theme:                 NoneTheme,
				hasHeader:             false,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                MarkdownFormat,
					theme:                 NoneTheme,
					hasHeader:             false,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergeFields:           nil,
					ignoreFields:          nil,
					colorFlags:            nil,
				},
			},
		},
		{
			name: "markdown+merge",
			args: args{
				format:                MarkdownFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           []int{0, 1},
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                MarkdownFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergeFields:           []int{0, 1},
					ignoreFields:          nil,
					colorFlags:            nil,
				},
			},
		},
		{
			name: "markdown+ignore",
			args: args{
				format:                MarkdownFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				ignoreFields:          []int{2},
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                MarkdownFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergeFields:           nil,
					ignoreFields:          []int{2},
					colorFlags:            nil,
				},
			},
		},
		{
			name: "markdown+emptyFieldPlaceholder",
			args: args{
				format:                MarkdownFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: "NULL",
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                MarkdownFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: "NULL",
					wordDelimiter:         defaultWordDelimiter,
					mergeFields:           nil,
					ignoreFields:          nil,
					colorFlags:            nil,
				},
			},
		},
		{
			name: "markdown+wordDelimiter",
			args: args{
				format:                MarkdownFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         ",",
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                MarkdownFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         ",",
					mergeFields:           nil,
					ignoreFields:          nil,
					colorFlags:            nil,
				},
			},
		},
		{
			name: "markdown+edgecase",
			args: args{
				format:                MarkdownFormat,
				theme:                 NoneTheme,
				hasHeader:             false,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           []int{0, 1},
				ignoreFields:          []int{2},
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                MarkdownFormat,
					theme:                 NoneTheme,
					hasHeader:             false,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergeFields:           []int{0, 1},
					ignoreFields:          []int{2},
					colorFlags:            nil,
				},
			},
		},
		{
			name: "backlog+basic",
			args: args{
				format:                BacklogFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                BacklogFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergeFields:           nil,
					ignoreFields:          nil,
					colorFlags:            nil,
				},
			},
		},
		{
			name: "backlog+disableHeader",
			args: args{
				format:                BacklogFormat,
				theme:                 NoneTheme,
				hasHeader:             false,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                BacklogFormat,
					theme:                 NoneTheme,
					hasHeader:             false,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergeFields:           nil,
					ignoreFields:          nil,
					colorFlags:            nil,
				},
			},
		},
		{
			name: "backlog+merge",
			args: args{
				format:                BacklogFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           []int{0, 1},
				ignoreFields:          nil,
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                BacklogFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergeFields:           []int{0, 1},
					ignoreFields:          nil,
					colorFlags:            nil,
				},
			},
		},
		{
			name: "backlog+ignore",
			args: args{
				format:                BacklogFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           nil,
				ignoreFields:          []int{2},
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                BacklogFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergeFields:           nil,
					ignoreFields:          []int{2},
					colorFlags:            nil,
				},
			},
		},
		{
			name: "backlog+emptyFieldPlaceholder",
			args: args{
				format:                BacklogFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: "NULL",
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                BacklogFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: "NULL",
					wordDelimiter:         defaultWordDelimiter,
					mergeFields:           nil,
					ignoreFields:          nil,
					colorFlags:            nil,
				},
			},
		},
		{
			name: "backlog+wordDelimiter",
			args: args{
				format:                BacklogFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         ",",
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                BacklogFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         ",",
					mergeFields:           nil,
					ignoreFields:          nil,
					colorFlags:            nil,
				},
			},
		},
		{
			name: "backlog+edgecase",
			args: args{
				format:                BacklogFormat,
				theme:                 NoneTheme,
				hasHeader:             false,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           []int{0, 1},
				ignoreFields:          []int{2},
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                BacklogFormat,
					theme:                 NoneTheme,
					hasHeader:             false,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergeFields:           []int{0, 1},
					ignoreFields:          []int{2},
					colorFlags:            nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTable(
				WithFormat(tt.args.format),
				WithTheme(tt.args.theme),
				WithHeader(tt.args.hasHeader),
				WithEmptyFieldPlaceholder(tt.args.emptyFieldPlaceholder),
				WithWordDelimiter(tt.args.wordDelimiter),
				WithMergeFields(tt.args.mergeFields),
				WithIgnoreFields(tt.args.ignoreFields),
			)
			if !reflect.DeepEqual(got, tt.want.got) {
				t.Errorf("got: %v, want: %v", got, tt.want.got)
			}
		})
	}
}

func TestTable_Load(t *testing.T) {
	type fields struct {
		emptyFieldPlaceholder string
		wordDelimiter         string
		mergeFields           []int
		ignoreFields          []int
	}
	type args struct {
		input any
	}
	type want struct {
		got *Table
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "basic",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			args: args{
				input: samples,
			},
			want: want{
				got: &Table{
					data: [][]string{
						{"i-1", "sg-1", "10.0.0.0/16"},
						{"i-1", "sg-1", "10.1.0.0/16"},
						{"i-1", "sg-2", "10.2.0.0/16"},
						{"i-1", "sg-2", "10.3.0.0/16"},
						{"i-2", "sg-1", "10.0.0.0/16<br>0.0.0.0/0"},
						{"i-2", "sg-1", "10.1.0.0/16<br>0.0.0.0/0"},
						{"i-2", "sg-2", "10.2.0.0/16<br>0.0.0.0/0"},
						{"i-2", "sg-2", "10.3.0.0/16<br>0.0.0.0/0"},
						{"i-3", "N/A", "10.0.0.0/16<br>0.0.0.0/0"},
						{"i-4", "sg-4", "N/A"},
					},
					headers:               []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					colorFlags:            []bool{true, true, true, true, false, false, false, false, true, false},
				},
				err: nil,
			},
		},
		{
			name: "emptyFieldPlaceholder",
			fields: fields{
				emptyFieldPlaceholder: "NULL",
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			args: args{
				input: samples,
			},
			want: want{
				got: &Table{
					data: [][]string{
						{"i-1", "sg-1", "10.0.0.0/16"},
						{"i-1", "sg-1", "10.1.0.0/16"},
						{"i-1", "sg-2", "10.2.0.0/16"},
						{"i-1", "sg-2", "10.3.0.0/16"},
						{"i-2", "sg-1", "10.0.0.0/16<br>0.0.0.0/0"},
						{"i-2", "sg-1", "10.1.0.0/16<br>0.0.0.0/0"},
						{"i-2", "sg-2", "10.2.0.0/16<br>0.0.0.0/0"},
						{"i-2", "sg-2", "10.3.0.0/16<br>0.0.0.0/0"},
						{"i-3", "NULL", "10.0.0.0/16<br>0.0.0.0/0"},
						{"i-4", "sg-4", "NULL"},
					},
					headers:               []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
					emptyFieldPlaceholder: "NULL",
					wordDelimiter:         defaultWordDelimiter,
					colorFlags:            []bool{true, true, true, true, false, false, false, false, true, false},
				},
				err: nil,
			},
		},
		{
			name: "wordDelimiter",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         ",",
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			args: args{
				input: samples,
			},
			want: want{
				got: &Table{
					data: [][]string{
						{"i-1", "sg-1", "10.0.0.0/16"},
						{"i-1", "sg-1", "10.1.0.0/16"},
						{"i-1", "sg-2", "10.2.0.0/16"},
						{"i-1", "sg-2", "10.3.0.0/16"},
						{"i-2", "sg-1", "10.0.0.0/16,0.0.0.0/0"},
						{"i-2", "sg-1", "10.1.0.0/16,0.0.0.0/0"},
						{"i-2", "sg-2", "10.2.0.0/16,0.0.0.0/0"},
						{"i-2", "sg-2", "10.3.0.0/16,0.0.0.0/0"},
						{"i-3", "N/A", "10.0.0.0/16,0.0.0.0/0"},
						{"i-4", "sg-4", "N/A"},
					},
					headers:               []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         ",",
					colorFlags:            []bool{true, true, true, true, false, false, false, false, true, false},
				},
				err: nil,
			},
		},
		{
			name: "merge",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           []int{0, 1},
				ignoreFields:          nil,
			},
			args: args{
				input: samples,
			},
			want: want{
				got: &Table{
					data: [][]string{
						{"i-1", "sg-1", "10.0.0.0/16"},
						{"", "", "10.1.0.0/16"},
						{"", "sg-2", "10.2.0.0/16"},
						{"", "", "10.3.0.0/16"},
						{"i-2", "sg-1", "10.0.0.0/16<br>0.0.0.0/0"},
						{"", "", "10.1.0.0/16<br>0.0.0.0/0"},
						{"", "sg-2", "10.2.0.0/16<br>0.0.0.0/0"},
						{"", "", "10.3.0.0/16<br>0.0.0.0/0"},
						{"i-3", "N/A", "10.0.0.0/16<br>0.0.0.0/0"},
						{"i-4", "sg-4", "N/A"},
					},
					headers:               []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergeFields:           []int{0, 1},
					colorFlags:            []bool{true, true, true, true, false, false, false, false, true, false},
				},
				err: nil,
			},
		},
		{
			name: "ignore",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           nil,
				ignoreFields:          []int{2},
			},
			args: args{
				input: samples,
			},
			want: want{
				got: &Table{
					data: [][]string{
						{"i-1", "sg-1"},
						{"i-1", "sg-1"},
						{"i-1", "sg-2"},
						{"i-1", "sg-2"},
						{"i-2", "sg-1"},
						{"i-2", "sg-1"},
						{"i-2", "sg-2"},
						{"i-2", "sg-2"},
						{"i-3", "N/A"},
						{"i-4", "sg-4"},
					},
					headers:               []string{"InstanceName", "SecurityGroupName"},
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					ignoreFields:          []int{2},

					colorFlags: []bool{true, true, true, true, false, false, false, false, true, false},
				},
				err: nil,
			},
		},
		{
			name: "merge+ignore",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           []int{0, 1},
				ignoreFields:          []int{2},
			},
			args: args{
				input: samples,
			},
			want: want{
				got: &Table{
					data: [][]string{
						{"i-1", "sg-1"},
						{"", ""},
						{"", "sg-2"},
						{"", ""},
						{"i-2", "sg-1"},
						{"", ""},
						{"", "sg-2"},
						{"", ""},
						{"i-3", "N/A"},
						{"i-4", "sg-4"},
					},
					headers:               []string{"InstanceName", "SecurityGroupName"},
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergeFields:           []int{0, 1},
					ignoreFields:          []int{2},
					colorFlags:            []bool{true, true, true, true, false, false, false, false, true, false},
				},
				err: nil,
			},
		},
		{
			name: "edgecase",
			fields: fields{
				emptyFieldPlaceholder: "",
				wordDelimiter:         ",",
				mergeFields:           []int{0, 1},
				ignoreFields:          []int{2},
			},
			args: args{
				input: samples,
			},
			want: want{
				got: &Table{
					data: [][]string{
						{"i-1", "sg-1"},
						{"", ""},
						{"", "sg-2"},
						{"", ""},
						{"i-2", "sg-1"},
						{"", ""},
						{"", "sg-2"},
						{"", ""},
						{"i-3", ""},
						{"i-4", "sg-4"},
					},
					headers:               []string{"InstanceName", "SecurityGroupName"},
					emptyFieldPlaceholder: "",
					wordDelimiter:         ",",
					mergeFields:           []int{0, 1},
					ignoreFields:          []int{2},
					colorFlags:            []bool{true, true, true, true, false, false, false, false, true, false},
				},
				err: nil,
			},
		},
		{
			name: "string slice",
			args: args{
				input: []string{"a", "b", "c"},
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: must be struct"),
			},
		},
		{
			name: "int slice",
			args: args{
				input: []int{1, 2, 3},
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: must be struct"),
			},
		},
		{
			name: "bool slice",
			args: args{
				input: []bool{true, false, true},
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: must be struct"),
			},
		},
		{
			name: "rune slice",
			args: args{
				input: []rune{'a', 'b', 'c'},
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: must be struct"),
			},
		},
		{
			name: "string",
			args: args{
				input: "a",
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: must be slice"),
			},
		},
		{
			name: "int",
			args: args{
				input: 1,
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: must be slice"),
			},
		},
		{
			name: "bool",
			args: args{
				input: true,
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: must be slice"),
			},
		},
		{
			name: "rune",
			args: args{
				input: 'a',
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: must be slice"),
			},
		},
		{
			name: "slice in slice",
			args: args{
				input: [][]string{
					{"a", "b", "c"},
					{"x", "y", "z"},
				},
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: must be struct"),
			},
		},
		{
			name: "empty slice",
			args: args{
				input: []string{},
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: no elements in slice"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{
				emptyFieldPlaceholder: tt.fields.emptyFieldPlaceholder,
				wordDelimiter:         tt.fields.wordDelimiter,
				mergeFields:           tt.fields.mergeFields,
				ignoreFields:          tt.fields.ignoreFields,
			}
			if err := table.Load(tt.args.input); err != nil {
				if err.Error() != tt.want.err.Error() {
					t.Fatalf("got: %v, want: %v", err.Error(), tt.want.err.Error())
				}
				return
			}
			if !reflect.DeepEqual(table, tt.want.got) {
				t.Errorf("got: %v, want: %v", table, tt.want.got)
			}
		})
	}
}

func TestTable_Out(t *testing.T) {
	type fields struct {
		data                  [][]string
		headers               []string
		format                int
		theme                 int
		hasHeader             bool
		emptyFieldPlaceholder string
		wordDelimiter         string
		mergeFields           []int
		ignoreFields          []int
		colorFlags            []bool
	}
	type want struct {
		got string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "markdown+basic",
			fields: fields{
				data: [][]string{
					{"i-1", "sg-1", "10.0.0.0/16"},
					{"i-1", "sg-1", "10.1.0.0/16"},
					{"i-1", "sg-2", "10.2.0.0/16"},
					{"i-1", "sg-2", "10.3.0.0/16"},
					{"i-2", "sg-1", "10.0.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-1", "10.1.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-2", "10.2.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-2", "10.3.0.0/16<br>0.0.0.0/0"},
					{"i-3", "N/A", "10.0.0.0/16<br>0.0.0.0/0"},
					{"i-4", "sg-4", "N/A"},
				},
				headers:   []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
				format:    MarkdownFormat,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName | CidrBlock                |
|--------------|-------------------|--------------------------|
| i-1          | sg-1              | 10.0.0.0/16              |
| i-1          | sg-1              | 10.1.0.0/16              |
| i-1          | sg-2              | 10.2.0.0/16              |
| i-1          | sg-2              | 10.3.0.0/16              |
| i-2          | sg-1              | 10.0.0.0/16<br>0.0.0.0/0 |
| i-2          | sg-1              | 10.1.0.0/16<br>0.0.0.0/0 |
| i-2          | sg-2              | 10.2.0.0/16<br>0.0.0.0/0 |
| i-2          | sg-2              | 10.3.0.0/16<br>0.0.0.0/0 |
| i-3          | N/A               | 10.0.0.0/16<br>0.0.0.0/0 |
| i-4          | sg-4              | N/A                      |
`,
			},
		},
		{
			name: "markdown+disableHeader",
			fields: fields{
				data: [][]string{
					{"i-1", "sg-1", "10.0.0.0/16"},
					{"i-1", "sg-1", "10.1.0.0/16"},
					{"i-1", "sg-2", "10.2.0.0/16"},
					{"i-1", "sg-2", "10.3.0.0/16"},
					{"i-2", "sg-1", "10.0.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-1", "10.1.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-2", "10.2.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-2", "10.3.0.0/16<br>0.0.0.0/0"},
					{"i-3", "N/A", "10.0.0.0/16<br>0.0.0.0/0"},
					{"i-4", "sg-4", "N/A"},
				},
				headers:   []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
				format:    MarkdownFormat,
				theme:     NoneTheme,
				hasHeader: false,
			},
			want: want{
				got: `| i-1 | sg-1 | 10.0.0.0/16              |
| i-1 | sg-1 | 10.1.0.0/16              |
| i-1 | sg-2 | 10.2.0.0/16              |
| i-1 | sg-2 | 10.3.0.0/16              |
| i-2 | sg-1 | 10.0.0.0/16<br>0.0.0.0/0 |
| i-2 | sg-1 | 10.1.0.0/16<br>0.0.0.0/0 |
| i-2 | sg-2 | 10.2.0.0/16<br>0.0.0.0/0 |
| i-2 | sg-2 | 10.3.0.0/16<br>0.0.0.0/0 |
| i-3 | N/A  | 10.0.0.0/16<br>0.0.0.0/0 |
| i-4 | sg-4 | N/A                      |
`,
			},
		},
		{
			name: "markdown+merge",
			fields: fields{
				data: [][]string{
					{"i-1", "sg-1", "10.0.0.0/16"},
					{"", "", "10.1.0.0/16"},
					{"", "sg-2", "10.2.0.0/16"},
					{"", "", "10.3.0.0/16"},
					{"i-2", "sg-1", "10.0.0.0/16<br>0.0.0.0/0"},
					{"", "", "10.1.0.0/16<br>0.0.0.0/0"},
					{"", "sg-2", "10.2.0.0/16<br>0.0.0.0/0"},
					{"", "", "10.3.0.0/16<br>0.0.0.0/0"},
					{"i-3", "N/A", "10.0.0.0/16<br>0.0.0.0/0"},
					{"i-4", "sg-4", "N/A"},
				},
				headers:   []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
				format:    MarkdownFormat,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName | CidrBlock                |
|--------------|-------------------|--------------------------|
| i-1          | sg-1              | 10.0.0.0/16              |
|              |                   | 10.1.0.0/16              |
|              | sg-2              | 10.2.0.0/16              |
|              |                   | 10.3.0.0/16              |
| i-2          | sg-1              | 10.0.0.0/16<br>0.0.0.0/0 |
|              |                   | 10.1.0.0/16<br>0.0.0.0/0 |
|              | sg-2              | 10.2.0.0/16<br>0.0.0.0/0 |
|              |                   | 10.3.0.0/16<br>0.0.0.0/0 |
| i-3          | N/A               | 10.0.0.0/16<br>0.0.0.0/0 |
| i-4          | sg-4              | N/A                      |
`,
			},
		},
		{
			name: "markdown+ignore",
			fields: fields{
				data: [][]string{
					{"i-1", "sg-1"},
					{"i-1", "sg-1"},
					{"i-1", "sg-2"},
					{"i-1", "sg-2"},
					{"i-2", "sg-1"},
					{"i-2", "sg-1"},
					{"i-2", "sg-2"},
					{"i-2", "sg-2"},
					{"i-3", "N/A"},
					{"i-4", "sg-4"},
				},
				headers:   []string{"InstanceName", "SecurityGroupName"},
				format:    MarkdownFormat,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |
|--------------|-------------------|
| i-1          | sg-1              |
| i-1          | sg-1              |
| i-1          | sg-2              |
| i-1          | sg-2              |
| i-2          | sg-1              |
| i-2          | sg-1              |
| i-2          | sg-2              |
| i-2          | sg-2              |
| i-3          | N/A               |
| i-4          | sg-4              |
`,
			},
		},
		{
			name: "markdown+emptyFieldPlaceholder",
			fields: fields{
				data: [][]string{
					{"i-1", "sg-1", "10.0.0.0/16"},
					{"i-1", "sg-1", "10.1.0.0/16"},
					{"i-1", "sg-2", "10.2.0.0/16"},
					{"i-1", "sg-2", "10.3.0.0/16"},
					{"i-2", "sg-1", "10.0.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-1", "10.1.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-2", "10.2.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-2", "10.3.0.0/16<br>0.0.0.0/0"},
					{"i-3", "NULL", "10.0.0.0/16<br>0.0.0.0/0"},
					{"i-4", "sg-4", "NULL"},
				},
				headers:   []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
				format:    MarkdownFormat,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName | CidrBlock                |
|--------------|-------------------|--------------------------|
| i-1          | sg-1              | 10.0.0.0/16              |
| i-1          | sg-1              | 10.1.0.0/16              |
| i-1          | sg-2              | 10.2.0.0/16              |
| i-1          | sg-2              | 10.3.0.0/16              |
| i-2          | sg-1              | 10.0.0.0/16<br>0.0.0.0/0 |
| i-2          | sg-1              | 10.1.0.0/16<br>0.0.0.0/0 |
| i-2          | sg-2              | 10.2.0.0/16<br>0.0.0.0/0 |
| i-2          | sg-2              | 10.3.0.0/16<br>0.0.0.0/0 |
| i-3          | NULL              | 10.0.0.0/16<br>0.0.0.0/0 |
| i-4          | sg-4              | NULL                     |
`,
			},
		},
		{
			name: "markdown+wordDelimiter",
			fields: fields{
				data: [][]string{
					{"i-1", "sg-1", "10.0.0.0/16"},
					{"i-1", "sg-1", "10.1.0.0/16"},
					{"i-1", "sg-2", "10.2.0.0/16"},
					{"i-1", "sg-2", "10.3.0.0/16"},
					{"i-2", "sg-1", "10.0.0.0/16,0.0.0.0/0"},
					{"i-2", "sg-1", "10.1.0.0/16,0.0.0.0/0"},
					{"i-2", "sg-2", "10.2.0.0/16,0.0.0.0/0"},
					{"i-2", "sg-2", "10.3.0.0/16,0.0.0.0/0"},
					{"i-3", "N/A", "10.0.0.0/16,0.0.0.0/0"},
					{"i-4", "sg-4", "N/A"},
				},
				headers:   []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
				format:    MarkdownFormat,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName | CidrBlock             |
|--------------|-------------------|-----------------------|
| i-1          | sg-1              | 10.0.0.0/16           |
| i-1          | sg-1              | 10.1.0.0/16           |
| i-1          | sg-2              | 10.2.0.0/16           |
| i-1          | sg-2              | 10.3.0.0/16           |
| i-2          | sg-1              | 10.0.0.0/16,0.0.0.0/0 |
| i-2          | sg-1              | 10.1.0.0/16,0.0.0.0/0 |
| i-2          | sg-2              | 10.2.0.0/16,0.0.0.0/0 |
| i-2          | sg-2              | 10.3.0.0/16,0.0.0.0/0 |
| i-3          | N/A               | 10.0.0.0/16,0.0.0.0/0 |
| i-4          | sg-4              | N/A                   |
`,
			},
		},
		{
			name: "markdown+edgecase",
			fields: fields{
				data: [][]string{
					{"i-1", "sg-1"},
					{"", ""},
					{"", "sg-2"},
					{"", ""},
					{"i-2", "sg-1"},
					{"", ""},
					{"", "sg-2"},
					{"", ""},
					{"i-3", "N/A"},
					{"i-4", "sg-4"},
				},
				headers:   []string{"InstanceName", "SecurityGroupName"},
				format:    MarkdownFormat,
				theme:     NoneTheme,
				hasHeader: false,
			},
			want: want{
				got: `| i-1 | sg-1 |
|     |      |
|     | sg-2 |
|     |      |
| i-2 | sg-1 |
|     |      |
|     | sg-2 |
|     |      |
| i-3 | N/A  |
| i-4 | sg-4 |
`,
			},
		},
		{
			name: "backlog+basic",
			fields: fields{
				data: [][]string{
					{"i-1", "sg-1", "10.0.0.0/16"},
					{"i-1", "sg-1", "10.1.0.0/16"},
					{"i-1", "sg-2", "10.2.0.0/16"},
					{"i-1", "sg-2", "10.3.0.0/16"},
					{"i-2", "sg-1", "10.0.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-1", "10.1.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-2", "10.2.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-2", "10.3.0.0/16<br>0.0.0.0/0"},
					{"i-3", "N/A", "10.0.0.0/16<br>0.0.0.0/0"},
					{"i-4", "sg-4", "N/A"},
				},
				headers:   []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
				format:    BacklogFormat,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName | CidrBlock                |h
| i-1          | sg-1              | 10.0.0.0/16              |
| i-1          | sg-1              | 10.1.0.0/16              |
| i-1          | sg-2              | 10.2.0.0/16              |
| i-1          | sg-2              | 10.3.0.0/16              |
| i-2          | sg-1              | 10.0.0.0/16&br;0.0.0.0/0 |
| i-2          | sg-1              | 10.1.0.0/16&br;0.0.0.0/0 |
| i-2          | sg-2              | 10.2.0.0/16&br;0.0.0.0/0 |
| i-2          | sg-2              | 10.3.0.0/16&br;0.0.0.0/0 |
| i-3          | N/A               | 10.0.0.0/16&br;0.0.0.0/0 |
| i-4          | sg-4              | N/A                      |
`,
			},
		},
		{
			name: "backlog+disableHeader",
			fields: fields{
				data: [][]string{
					{"i-1", "sg-1", "10.0.0.0/16"},
					{"i-1", "sg-1", "10.1.0.0/16"},
					{"i-1", "sg-2", "10.2.0.0/16"},
					{"i-1", "sg-2", "10.3.0.0/16"},
					{"i-2", "sg-1", "10.0.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-1", "10.1.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-2", "10.2.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-2", "10.3.0.0/16<br>0.0.0.0/0"},
					{"i-3", "N/A", "10.0.0.0/16<br>0.0.0.0/0"},
					{"i-4", "sg-4", "N/A"},
				},
				headers:   []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
				format:    BacklogFormat,
				theme:     NoneTheme,
				hasHeader: false,
			},
			want: want{
				got: `| i-1 | sg-1 | 10.0.0.0/16              |
| i-1 | sg-1 | 10.1.0.0/16              |
| i-1 | sg-2 | 10.2.0.0/16              |
| i-1 | sg-2 | 10.3.0.0/16              |
| i-2 | sg-1 | 10.0.0.0/16&br;0.0.0.0/0 |
| i-2 | sg-1 | 10.1.0.0/16&br;0.0.0.0/0 |
| i-2 | sg-2 | 10.2.0.0/16&br;0.0.0.0/0 |
| i-2 | sg-2 | 10.3.0.0/16&br;0.0.0.0/0 |
| i-3 | N/A  | 10.0.0.0/16&br;0.0.0.0/0 |
| i-4 | sg-4 | N/A                      |
`,
			},
		},
		{
			name: "backlog+merge",
			fields: fields{
				data: [][]string{
					{"i-1", "sg-1", "10.0.0.0/16"},
					{"", "", "10.1.0.0/16"},
					{"", "sg-2", "10.2.0.0/16"},
					{"", "", "10.3.0.0/16"},
					{"i-2", "sg-1", "10.0.0.0/16<br>0.0.0.0/0"},
					{"", "", "10.1.0.0/16<br>0.0.0.0/0"},
					{"", "sg-2", "10.2.0.0/16<br>0.0.0.0/0"},
					{"", "", "10.3.0.0/16<br>0.0.0.0/0"},
					{"i-3", "N/A", "10.0.0.0/16<br>0.0.0.0/0"},
					{"i-4", "sg-4", "N/A"},
				},
				headers:   []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
				format:    BacklogFormat,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName | CidrBlock                |h
| i-1          | sg-1              | 10.0.0.0/16              |
|              |                   | 10.1.0.0/16              |
|              | sg-2              | 10.2.0.0/16              |
|              |                   | 10.3.0.0/16              |
| i-2          | sg-1              | 10.0.0.0/16&br;0.0.0.0/0 |
|              |                   | 10.1.0.0/16&br;0.0.0.0/0 |
|              | sg-2              | 10.2.0.0/16&br;0.0.0.0/0 |
|              |                   | 10.3.0.0/16&br;0.0.0.0/0 |
| i-3          | N/A               | 10.0.0.0/16&br;0.0.0.0/0 |
| i-4          | sg-4              | N/A                      |
`,
			},
		},
		{
			name: "backlog+ignore",
			fields: fields{
				data: [][]string{
					{"i-1", "sg-1"},
					{"i-1", "sg-1"},
					{"i-1", "sg-2"},
					{"i-1", "sg-2"},
					{"i-2", "sg-1"},
					{"i-2", "sg-1"},
					{"i-2", "sg-2"},
					{"i-2", "sg-2"},
					{"i-3", "N/A"},
					{"i-4", "sg-4"},
				},
				headers:   []string{"InstanceName", "SecurityGroupName"},
				format:    BacklogFormat,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |h
| i-1          | sg-1              |
| i-1          | sg-1              |
| i-1          | sg-2              |
| i-1          | sg-2              |
| i-2          | sg-1              |
| i-2          | sg-1              |
| i-2          | sg-2              |
| i-2          | sg-2              |
| i-3          | N/A               |
| i-4          | sg-4              |
`,
			},
		},
		{
			name: "backlog+emptyFieldPlaceholder",
			fields: fields{
				data: [][]string{
					{"i-1", "sg-1", "10.0.0.0/16"},
					{"i-1", "sg-1", "10.1.0.0/16"},
					{"i-1", "sg-2", "10.2.0.0/16"},
					{"i-1", "sg-2", "10.3.0.0/16"},
					{"i-2", "sg-1", "10.0.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-1", "10.1.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-2", "10.2.0.0/16<br>0.0.0.0/0"},
					{"i-2", "sg-2", "10.3.0.0/16<br>0.0.0.0/0"},
					{"i-3", "NULL", "10.0.0.0/16<br>0.0.0.0/0"},
					{"i-4", "sg-4", "NULL"},
				},
				headers:   []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
				format:    BacklogFormat,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName | CidrBlock                |h
| i-1          | sg-1              | 10.0.0.0/16              |
| i-1          | sg-1              | 10.1.0.0/16              |
| i-1          | sg-2              | 10.2.0.0/16              |
| i-1          | sg-2              | 10.3.0.0/16              |
| i-2          | sg-1              | 10.0.0.0/16&br;0.0.0.0/0 |
| i-2          | sg-1              | 10.1.0.0/16&br;0.0.0.0/0 |
| i-2          | sg-2              | 10.2.0.0/16&br;0.0.0.0/0 |
| i-2          | sg-2              | 10.3.0.0/16&br;0.0.0.0/0 |
| i-3          | NULL              | 10.0.0.0/16&br;0.0.0.0/0 |
| i-4          | sg-4              | NULL                     |
`,
			},
		},
		{
			name: "backlog+wordDelimiter",
			fields: fields{
				data: [][]string{
					{"i-1", "sg-1", "10.0.0.0/16"},
					{"i-1", "sg-1", "10.1.0.0/16"},
					{"i-1", "sg-2", "10.2.0.0/16"},
					{"i-1", "sg-2", "10.3.0.0/16"},
					{"i-2", "sg-1", "10.0.0.0/16,0.0.0.0/0"},
					{"i-2", "sg-1", "10.1.0.0/16,0.0.0.0/0"},
					{"i-2", "sg-2", "10.2.0.0/16,0.0.0.0/0"},
					{"i-2", "sg-2", "10.3.0.0/16,0.0.0.0/0"},
					{"i-3", "N/A", "10.0.0.0/16,0.0.0.0/0"},
					{"i-4", "sg-4", "N/A"},
				},
				headers:   []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
				format:    BacklogFormat,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName | CidrBlock             |h
| i-1          | sg-1              | 10.0.0.0/16           |
| i-1          | sg-1              | 10.1.0.0/16           |
| i-1          | sg-2              | 10.2.0.0/16           |
| i-1          | sg-2              | 10.3.0.0/16           |
| i-2          | sg-1              | 10.0.0.0/16,0.0.0.0/0 |
| i-2          | sg-1              | 10.1.0.0/16,0.0.0.0/0 |
| i-2          | sg-2              | 10.2.0.0/16,0.0.0.0/0 |
| i-2          | sg-2              | 10.3.0.0/16,0.0.0.0/0 |
| i-3          | N/A               | 10.0.0.0/16,0.0.0.0/0 |
| i-4          | sg-4              | N/A                   |
`,
			},
		},
		{
			name: "backlog+edgecase",
			fields: fields{
				data: [][]string{
					{"i-1", "sg-1"},
					{"", ""},
					{"", "sg-2"},
					{"", ""},
					{"i-2", "sg-1"},
					{"", ""},
					{"", "sg-2"},
					{"", ""},
					{"i-3", "N/A"},
					{"i-4", "sg-4"},
				},
				headers:   []string{"InstanceName", "SecurityGroupName"},
				format:    BacklogFormat,
				theme:     NoneTheme,
				hasHeader: false,
			},
			want: want{
				got: `| i-1 | sg-1 |
|     |      |
|     | sg-2 |
|     |      |
| i-2 | sg-1 |
|     |      |
|     | sg-2 |
|     |      |
| i-3 | N/A  |
| i-4 | sg-4 |
`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{
				data:                  tt.fields.data,
				headers:               tt.fields.headers,
				format:                tt.fields.format,
				theme:                 tt.fields.theme,
				hasHeader:             tt.fields.hasHeader,
				emptyFieldPlaceholder: tt.fields.emptyFieldPlaceholder,
				wordDelimiter:         tt.fields.wordDelimiter,
				mergeFields:           tt.fields.mergeFields,
				ignoreFields:          tt.fields.ignoreFields,
				colorFlags:            tt.fields.colorFlags,
			}
			if got := table.Out(); !reflect.DeepEqual(got, tt.want.got) {
				t.Errorf("got: %v, want: %v", got, tt.want.got)
			}
		})
	}
}

func TestTable_formatValue(t *testing.T) {
	type fields struct {
		emptyFieldPlaceholder string
		wordDelimiter         string
	}
	type args struct {
		v any
	}
	type want struct {
		got string
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "empty string",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: "",
			},
			want: want{
				got: defaultEmptyFieldPlaceholder,
				err: nil,
			},
		},
		{
			name: "nil slice",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: ([]string)(nil),
			},
			want: want{
				got: defaultEmptyFieldPlaceholder,
				err: nil,
			},
		},
		{
			name: "empty slice",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: []string{},
			},
			want: want{
				got: defaultEmptyFieldPlaceholder,
				err: nil,
			},
		},
		{
			name: "slice with empty string",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: []string{"", ""},
			},
			want: want{
				got: defaultEmptyFieldPlaceholder + defaultWordDelimiter + defaultEmptyFieldPlaceholder,
				err: nil,
			},
		},
		{
			name: "slice with normal strings",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: []string{"a", "b"},
			},
			want: want{
				got: "a" + defaultWordDelimiter + "b",
				err: nil,
			},
		},
		{
			name: "mixed slice",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: []string{"a", "", "b"},
			},
			want: want{
				got: "a" + defaultWordDelimiter + defaultEmptyFieldPlaceholder + defaultWordDelimiter + "b",
				err: nil,
			},
		},
		{
			name: "non-slice value",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: 123,
			},
			want: want{
				got: "123",
				err: nil,
			},
		},
		{
			name: "slice in slice",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: [][]string{
					{"a", "b", "c"},
					{"x", "y", "z"},
				},
			},
			want: want{
				got: "",
				err: fmt.Errorf("elements of slice must not be nested"),
			},
		},
		{
			name: "struct in slice",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: nests,
			},
			want: want{
				got: "",
				err: fmt.Errorf("elements of slice must not be nested"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.args.v)
			table := &Table{
				emptyFieldPlaceholder: tt.fields.emptyFieldPlaceholder,
				wordDelimiter:         tt.fields.wordDelimiter,
			}
			got, err := table.formatValue(v)
			if err != nil {
				if err.Error() != tt.want.err.Error() {
					t.Fatalf("got: %v, want: %v", err.Error(), tt.want.err.Error())
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want.got) {
				t.Errorf("got: %v, want: %v", got, tt.want.got)
			}
		})
	}
}

func TestTable_getOffset(t *testing.T) {
	type fields struct {
		format    int
		hasHeader bool
	}
	type want struct {
		got int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "markdown",
			fields: fields{
				format:    MarkdownFormat,
				hasHeader: true,
			},
			want: want{
				got: 2,
			},
		},
		{
			name: "backlog",
			fields: fields{
				format:    BacklogFormat,
				hasHeader: true,
			},
			want: want{
				got: 1,
			},
		},
		{
			name: "markdown+disableHeader",
			fields: fields{
				format:    MarkdownFormat,
				hasHeader: false,
			},
			want: want{
				got: 0,
			},
		},
		{
			name: "backlog+disableHeader",
			fields: fields{
				format:    BacklogFormat,
				hasHeader: false,
			},
			want: want{
				got: 0,
			},
		},
		{
			name: "default",
			fields: fields{
				format: 9,
			},
			want: want{
				got: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewTable(
				WithFormat(tt.fields.format),
				WithHeader(tt.fields.hasHeader),
			)
			got := table.getOffset()
			if !reflect.DeepEqual(got, tt.want.got) {
				t.Errorf("got: %v, want: %v", got, tt.want.got)
			}
		})
	}
}

func TestTable_getColor(t *testing.T) {
	type fields struct {
		theme int
	}
	type want struct {
		got *color.Color
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "dark",
			fields: fields{
				theme: DarkTheme,
			},
			want: want{
				got: color.New(color.BgHiBlack, color.FgHiWhite),
			},
		},
		{
			name: "light",
			fields: fields{
				theme: LightTheme,
			},
			want: want{
				got: color.New(color.BgHiWhite, color.FgHiBlack),
			},
		},
		{
			name: "none",
			fields: fields{
				theme: NoneTheme,
			},
			want: want{
				got: color.New(color.Reset),
			},
		},
		{
			name: "default",
			fields: fields{
				theme: 9,
			},
			want: want{
				got: color.New(color.Reset),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewTable(
				WithTheme(tt.fields.theme),
			)
			got := table.getColor()
			if !reflect.DeepEqual(got, tt.want.got) {
				t.Errorf("got: %v, want: %v", got, tt.want.got)
			}
		})
	}
}
