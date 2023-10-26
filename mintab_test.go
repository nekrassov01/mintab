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

var samples []sample

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
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
}

func TestNew(t *testing.T) {
	type args struct {
		input        any
		mergeFields  []int
		ignoreFields []int
	}
	type want struct {
		table *Table
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "basic",
			args: args{
				input:        samples,
				mergeFields:  nil,
				ignoreFields: nil,
			},
			want: want{
				table: &Table{
					Data: [][]string{
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
					headers:    []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
					colorFlags: []bool{true, true, true, true, false, false, false, false, true, false},
				},
			},
		},
		{
			name: "merge",
			args: args{
				input:        samples,
				mergeFields:  []int{0, 1},
				ignoreFields: nil,
			},
			want: want{
				table: &Table{
					Data: [][]string{
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
					headers:    []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
					colorFlags: []bool{true, true, true, true, false, false, false, false, true, false},
				},
			},
		},
		{
			name: "ignore",
			args: args{
				input:        samples,
				mergeFields:  nil,
				ignoreFields: []int{2},
			},
			want: want{
				table: &Table{
					Data: [][]string{
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
					headers:    []string{"InstanceName", "SecurityGroupName"},
					colorFlags: []bool{true, true, true, true, false, false, false, false, true, false},
				},
			},
		},
		{
			name: "merge+ignore",
			args: args{
				input:        samples,
				mergeFields:  []int{0, 1},
				ignoreFields: []int{2},
			},
			want: want{
				table: &Table{
					Data: [][]string{
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
					headers:    []string{"InstanceName", "SecurityGroupName"},
					colorFlags: []bool{true, true, true, true, false, false, false, false, true, false},
				},
			},
		},
		{
			name: "string slice",
			args: args{
				input: []string{"a", "b", "c"},
			},
			want: want{
				table: nil,
			},
		},
		{
			name: "int slice",
			args: args{
				input: []int{1, 2, 3},
			},
			want: want{
				table: nil,
			},
		},
		{
			name: "bool slice",
			args: args{
				input: []bool{true, false, true},
			},
			want: want{
				table: nil,
			},
		},
		{
			name: "rune slice",
			args: args{
				input: []rune{'a', 'b', 'c'},
			},
			want: want{
				table: nil,
			},
		},
		{
			name: "string",
			args: args{
				input: "a",
			},
			want: want{
				table: nil,
			},
		},
		{
			name: "int",
			args: args{
				input: 1,
			},
			want: want{
				table: nil,
			},
		},
		{
			name: "bool",
			args: args{
				input: true,
			},
			want: want{
				table: nil,
			},
		},
		{
			name: "rune",
			args: args{
				input: 'a',
			},
			want: want{
				table: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := New(tt.args.input, WithMergeFields(tt.args.mergeFields), WithIgnoreFields(tt.args.ignoreFields))
			if table == nil && tt.want.table == nil {
				return
			}
			if !reflect.DeepEqual(table.Data, tt.want.table.Data) {
				t.Errorf("got: %v, want: %v", table.Data, tt.want.table.Data)
			}
			if !reflect.DeepEqual(table.headers, tt.want.table.headers) {
				t.Errorf("got: %v, want: %v", table.headers, tt.want.table.headers)
			}
			if !reflect.DeepEqual(table.colorFlags, tt.want.table.colorFlags) {
				t.Errorf("got: %v, want: %v", table.colorFlags, tt.want.table.colorFlags)
			}
		})
	}
}

func Test_formatValue(t *testing.T) {
	type args struct {
		v                     any
		emptyFieldPlaceholder string
		wordDelimitter        string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty string",
			args: args{
				v:                     "",
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        defaultWordDelimitter,
			},
			want: defaultEmptyFieldPlaceholder,
		},
		{
			name: "nil slice",
			args: args{
				v:                     ([]string)(nil),
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        defaultWordDelimitter,
			},
			want: defaultEmptyFieldPlaceholder,
		},
		{
			name: "empty slice",
			args: args{
				v:                     []string{},
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        defaultWordDelimitter,
			},
			want: defaultEmptyFieldPlaceholder,
		},
		{
			name: "slice with empty string",
			args: args{
				v:                     []string{"", ""},
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        defaultWordDelimitter,
			},
			want: "N/A<br>N/A",
		},
		{
			name: "slice with normal strings",
			args: args{
				v:                     []string{"a", "b"},
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        defaultWordDelimitter,
			},
			want: "a<br>b",
		},
		{
			name: "mixed slice",
			args: args{
				v:                     []string{"a", "", "b"},
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        defaultWordDelimitter,
			},
			want: "a<br>N/A<br>b",
		},
		{
			name: "non-slice value",
			args: args{
				v:                     123,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        defaultWordDelimitter,
			},
			want: "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.args.v)
			got := formatValue(v, tt.args.emptyFieldPlaceholder, tt.args.wordDelimitter)
			if got != tt.want {
				t.Errorf("got: %v, got: %v", tt.want, got)
			}
		})
	}
}

func TestOut(t *testing.T) {
	type args struct {
		input                 any
		format                TableFormat
		theme                 TableTheme
		hasHeader             bool
		emptyFieldPlaceholder string
		wordDelimitter        string
		mergeFields           []int
		ignoreFields          []int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "markdown+basic",
			args: args{
				input:                 samples,
				format:                Markdown,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        defaultWordDelimitter,
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: `| InstanceName | SecurityGroupName |        CidrBlock         |
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
		{
			name: "markdown+disableHeader",
			args: args{
				input:                 samples,
				format:                Markdown,
				theme:                 NoneTheme,
				hasHeader:             false,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        defaultWordDelimitter,
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: `| i-1 | sg-1 | 10.0.0.0/16              |
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
		{
			name: "markdown+merge",
			args: args{
				input:                 samples,
				format:                Markdown,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        defaultWordDelimitter,
				mergeFields:           []int{0, 1},
				ignoreFields:          nil,
			},
			want: `| InstanceName | SecurityGroupName |        CidrBlock         |
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
		{
			name: "markdown+ignore",
			args: args{
				input:                 samples,
				format:                Markdown,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        defaultWordDelimitter,
				mergeFields:           nil,
				ignoreFields:          []int{2},
			},
			want: `| InstanceName | SecurityGroupName |
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
		{
			name: "markdown+emptyFieldPlaceholder",
			args: args{
				input:                 samples,
				format:                Markdown,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: "NULL",
				wordDelimitter:        defaultWordDelimitter,
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: `| InstanceName | SecurityGroupName |        CidrBlock         |
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
		{
			name: "markdown+wordDelimitter",
			args: args{
				input:                 samples,
				format:                Markdown,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        ",",
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: `| InstanceName | SecurityGroupName |       CidrBlock       |
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
		{
			name: "markdown+edgecase",
			args: args{
				input:                 samples,
				format:                Markdown,
				theme:                 NoneTheme,
				hasHeader:             false,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        defaultWordDelimitter,
				mergeFields:           []int{0, 1},
				ignoreFields:          []int{2},
			},
			want: `| i-1 | sg-1 |
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
		{
			name: "backlog+basic",
			args: args{
				input:                 samples,
				format:                Backlog,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        defaultWordDelimitter,
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: `| InstanceName | SecurityGroupName |        CidrBlock         |h
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
		{
			name: "backlog+disableHeader",
			args: args{
				input:                 samples,
				format:                Backlog,
				theme:                 NoneTheme,
				hasHeader:             false,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        defaultWordDelimitter,
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: `| i-1 | sg-1 | 10.0.0.0/16              |
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
		{
			name: "backlog+merge",
			args: args{
				input:                 samples,
				format:                Backlog,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        defaultWordDelimitter,
				mergeFields:           []int{0, 1},
				ignoreFields:          nil,
			},
			want: `| InstanceName | SecurityGroupName |        CidrBlock         |h
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
		{
			name: "backlog+ignore",
			args: args{
				input:                 samples,
				format:                Backlog,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        defaultWordDelimitter,
				mergeFields:           nil,
				ignoreFields:          []int{2},
			},
			want: `| InstanceName | SecurityGroupName |h
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
		{
			name: "backlog+emptyFieldPlaceholder",
			args: args{
				input:                 samples,
				format:                Backlog,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: "NULL",
				wordDelimitter:        defaultWordDelimitter,
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: `| InstanceName | SecurityGroupName |        CidrBlock         |h
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
		{
			name: "backlog+wordDelimitter",
			args: args{
				input:                 samples,
				format:                Backlog,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        ",",
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: `| InstanceName | SecurityGroupName |       CidrBlock       |h
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
		{
			name: "backlog+edgecase",
			args: args{
				input:                 samples,
				format:                Backlog,
				theme:                 NoneTheme,
				hasHeader:             false,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimitter:        defaultWordDelimitter,
				mergeFields:           []int{0, 1},
				ignoreFields:          []int{2},
			},
			want: `| i-1 | sg-1 |
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := New(
				tt.args.input,
				WithTableFormat(tt.args.format),
				WithTableTheme(tt.args.theme),
				WithTableHeader(tt.args.hasHeader),
				WithEmptyFieldPlaceholder(tt.args.emptyFieldPlaceholder),
				WithWordDelimitter(tt.args.wordDelimitter),
				WithMergeFields(tt.args.mergeFields),
				WithIgnoreFields(tt.args.ignoreFields),
			)
			got := table.Out()
			fmt.Println(got)
			if got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func Test_getOffset(t *testing.T) {
	type args struct {
		format    TableFormat
		hasHeader bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "markdown",
			args: args{
				format:    Markdown,
				hasHeader: true,
			},
			want: 2,
		},
		{
			name: "backlog",
			args: args{
				format:    Backlog,
				hasHeader: true,
			},
			want: 1,
		},
		{
			name: "markdown+disableHeader",
			args: args{
				format:    Markdown,
				hasHeader: false,
			},
			want: 0,
		},
		{
			name: "backlog+disableHeader",
			args: args{
				format:    Backlog,
				hasHeader: false,
			},
			want: 0,
		},
		{
			name: "default",
			args: args{
				format: "",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getOffset(tt.args.format, tt.args.hasHeader)
			if got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func Test_getColor(t *testing.T) {
	type args struct {
		tableTheme TableTheme
	}
	tests := []struct {
		name string
		args args
		want *color.Color
	}{
		{
			name: "dark",
			args: args{
				tableTheme: DarkTheme,
			},
			want: color.New(color.BgHiBlack, color.FgHiWhite),
		},
		{
			name: "light",
			args: args{
				tableTheme: LightTheme,
			},
			want: color.New(color.BgHiWhite, color.FgHiBlack),
		},
		{
			name: "none",
			args: args{
				tableTheme: NoneTheme,
			},
			want: color.New(color.Reset),
		},
		{
			name: "default",
			args: args{
				tableTheme: "",
			},
			want: color.New(color.Reset),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getColor(tt.args.tableTheme)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}
