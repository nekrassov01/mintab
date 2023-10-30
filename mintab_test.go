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

type object struct {
	ObjectID   int
	ObjectName string
}

type nested struct {
	BucketName string
	Objects    []object
}

var nests []nested

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

func TestNew(t *testing.T) {
	type args struct {
		input        any
		mergeFields  []int
		ignoreFields []int
	}
	type want struct {
		got *Table
		err error
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
				got: &Table{
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
				err: nil,
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
				got: &Table{
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
				err: nil,
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
				got: &Table{
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
				err: nil,
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
				got: &Table{
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
			got, err := New(tt.args.input, WithMergeFields(tt.args.mergeFields), WithIgnoreFields(tt.args.ignoreFields))
			if err != nil {
				if err.Error() != tt.want.err.Error() {
					t.Fatalf("got: %v, want: %v", err.Error(), tt.want.err.Error())
				}
				return
			}
			if !reflect.DeepEqual(got.Data, tt.want.got.Data) {
				t.Errorf("got: %v, want: %v", got.Data, tt.want.got.Data)
			}
			if !reflect.DeepEqual(got.headers, tt.want.got.headers) {
				t.Errorf("got: %v, want: %v", got.headers, tt.want.got.headers)
			}
			if !reflect.DeepEqual(got.colorFlags, tt.want.got.colorFlags) {
				t.Errorf("got: %v, want: %v", got.colorFlags, tt.want.got.colorFlags)
			}
		})
	}
}

func Test_formatValue(t *testing.T) {
	type args struct {
		v                     any
		emptyFieldPlaceholder string
		wordDelimiter         string
	}
	type want struct {
		got string
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty string",
			args: args{
				v:                     "",
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			want: want{
				got: defaultEmptyFieldPlaceholder,
				err: nil,
			},
		},
		{
			name: "nil slice",
			args: args{
				v:                     ([]string)(nil),
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			want: want{
				got: defaultEmptyFieldPlaceholder,
				err: nil,
			},
		},
		{
			name: "empty slice",
			args: args{
				v:                     []string{},
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			want: want{
				got: defaultEmptyFieldPlaceholder,
				err: nil,
			},
		},
		{
			name: "slice with empty string",
			args: args{
				v:                     []string{"", ""},
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			want: want{
				got: defaultEmptyFieldPlaceholder + defaultWordDelimiter + defaultEmptyFieldPlaceholder,
				err: nil,
			},
		},
		{
			name: "slice with normal strings",
			args: args{
				v:                     []string{"a", "b"},
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			want: want{
				got: "a" + defaultWordDelimiter + "b",
				err: nil,
			},
		},
		{
			name: "mixed slice",
			args: args{
				v:                     []string{"a", "", "b"},
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			want: want{
				got: "a" + defaultWordDelimiter + defaultEmptyFieldPlaceholder + defaultWordDelimiter + "b",
				err: nil,
			},
		},
		{
			name: "non-slice value",
			args: args{
				v:                     123,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			want: want{
				got: "123",
				err: nil,
			},
		},
		{
			name: "slice in slice",
			args: args{
				v: [][]string{
					{"a", "b", "c"},
					{"x", "y", "z"},
				},
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			want: want{
				got: "",
				err: fmt.Errorf("elements of slice must not be nested"),
			},
		},
		{
			name: "struct in slice",
			args: args{
				v:                     nests,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
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
			got, err := formatValue(v, tt.args.emptyFieldPlaceholder, tt.args.wordDelimiter)
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

func TestTable_Out(t *testing.T) {
	type args struct {
		Data                  [][]string
		headers               []string
		format                TableFormat
		theme                 TableTheme
		hasHeader             bool
		emptyFieldPlaceholder string
		wordDelimiter         string
		mergeFields           []int
		ignoreFields          []int
		colorFlags            []bool
	}
	type want struct {
		got string
		err error
	}
	tests := []struct {
		name string
		args *args
		want want
	}{
		{
			name: "markdown+basic",
			args: &args{
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
				headers:   []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
				format:    Markdown,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |        CidrBlock         |
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
				err: nil,
			},
		},
		{
			name: "markdown+disableHeader",
			args: &args{
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
				headers:   []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
				format:    Markdown,
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
				err: nil,
			},
		},
		{
			name: "markdown+merge",
			args: &args{
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
				headers:   []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
				format:    Markdown,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |        CidrBlock         |
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
				err: nil,
			},
		},
		{
			name: "markdown+ignore",
			args: &args{
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
				headers:   []string{"InstanceName", "SecurityGroupName"},
				format:    Markdown,
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
				err: nil,
			},
		},
		{
			name: "markdown+emptyFieldPlaceholder",
			args: &args{
				Data: [][]string{
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
				format:    Markdown,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |        CidrBlock         |
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
				err: nil,
			},
		},
		{
			name: "markdown+wordDelimiter",
			args: &args{
				Data: [][]string{
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
				format:    Markdown,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |       CidrBlock       |
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
				err: nil,
			},
		},
		{
			name: "markdown+edgecase",
			args: &args{
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
				headers:   []string{"InstanceName", "SecurityGroupName"},
				format:    Markdown,
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
				err: nil,
			},
		},
		{
			name: "backlog+basic",
			args: &args{
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
				headers:   []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
				format:    Backlog,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |        CidrBlock         |h
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
				err: nil,
			},
		},
		{
			name: "backlog+disableHeader",
			args: &args{
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
				headers:   []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
				format:    Backlog,
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
				err: nil,
			},
		},
		{
			name: "backlog+merge",
			args: &args{
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
				headers:   []string{"InstanceName", "SecurityGroupName", "CidrBlock"},
				format:    Backlog,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |        CidrBlock         |h
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
				err: nil,
			},
		},
		{
			name: "backlog+ignore",
			args: &args{
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
				headers:   []string{"InstanceName", "SecurityGroupName"},
				format:    Backlog,
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
				err: nil,
			},
		},
		{
			name: "backlog+emptyFieldPlaceholder",
			args: &args{
				Data: [][]string{
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
				format:    Backlog,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |        CidrBlock         |h
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
				err: nil,
			},
		},
		{
			name: "backlog+wordDelimiter",
			args: &args{
				Data: [][]string{
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
				format:    Backlog,
				theme:     NoneTheme,
				hasHeader: true,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |       CidrBlock       |h
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
				err: nil,
			},
		},
		{
			name: "backlog+edgecase",
			args: &args{
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
				headers:   []string{"InstanceName", "SecurityGroupName"},
				format:    Backlog,
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
				err: nil,
			},
		},
		{
			name: "table data is nil",
			args: &args{
				Data: nil,
			},
			want: want{
				got: "",
				err: fmt.Errorf("cannot parse table: empty data"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{
				Data:                  tt.args.Data,
				headers:               tt.args.headers,
				format:                tt.args.format,
				theme:                 tt.args.theme,
				hasHeader:             tt.args.hasHeader,
				emptyFieldPlaceholder: tt.args.emptyFieldPlaceholder,
				wordDelimiter:         tt.args.wordDelimiter,
				mergeFields:           tt.args.mergeFields,
				ignoreFields:          tt.args.ignoreFields,
				colorFlags:            tt.args.colorFlags,
			}
			got, err := table.Out()
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

func TestNewAndOut(t *testing.T) {
	type args struct {
		input                 any
		format                TableFormat
		theme                 TableTheme
		hasHeader             bool
		emptyFieldPlaceholder string
		wordDelimiter         string
		mergeFields           []int
		ignoreFields          []int
	}
	type want struct {
		got string
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "markdown+basic",
			args: args{
				input:                 samples,
				format:                Markdown,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |        CidrBlock         |
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
				err: nil,
			},
		},
		{
			name: "markdown+disableHeader",
			args: args{
				input:                 samples,
				format:                Markdown,
				theme:                 NoneTheme,
				hasHeader:             false,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           nil,
				ignoreFields:          nil,
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
				err: nil,
			},
		},
		{
			name: "markdown+merge",
			args: args{
				input:                 samples,
				format:                Markdown,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           []int{0, 1},
				ignoreFields:          nil,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |        CidrBlock         |
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
				err: nil,
			},
		},
		{
			name: "markdown+ignore",
			args: args{
				input:                 samples,
				format:                Markdown,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           nil,
				ignoreFields:          []int{2},
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
				err: nil,
			},
		},
		{
			name: "markdown+emptyFieldPlaceholder",
			args: args{
				input:                 samples,
				format:                Markdown,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: "NULL",
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |        CidrBlock         |
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
				err: nil,
			},
		},
		{
			name: "markdown+wordDelimiter",
			args: args{
				input:                 samples,
				format:                Markdown,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         ",",
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |       CidrBlock       |
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
				err: nil,
			},
		},
		{
			name: "markdown+edgecase",
			args: args{
				input:                 samples,
				format:                Markdown,
				theme:                 NoneTheme,
				hasHeader:             false,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           []int{0, 1},
				ignoreFields:          []int{2},
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
				err: nil,
			},
		},
		{
			name: "backlog+basic",
			args: args{
				input:                 samples,
				format:                Backlog,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |        CidrBlock         |h
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
				err: nil,
			},
		},
		{
			name: "backlog+disableHeader",
			args: args{
				input:                 samples,
				format:                Backlog,
				theme:                 NoneTheme,
				hasHeader:             false,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           nil,
				ignoreFields:          nil,
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
				err: nil,
			},
		},
		{
			name: "backlog+merge",
			args: args{
				input:                 samples,
				format:                Backlog,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           []int{0, 1},
				ignoreFields:          nil,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |        CidrBlock         |h
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
				err: nil,
			},
		},
		{
			name: "backlog+ignore",
			args: args{
				input:                 samples,
				format:                Backlog,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           nil,
				ignoreFields:          []int{2},
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
				err: nil,
			},
		},
		{
			name: "backlog+emptyFieldPlaceholder",
			args: args{
				input:                 samples,
				format:                Backlog,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: "NULL",
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |        CidrBlock         |h
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
				err: nil,
			},
		},
		{
			name: "backlog+wordDelimiter",
			args: args{
				input:                 samples,
				format:                Backlog,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         ",",
				mergeFields:           nil,
				ignoreFields:          nil,
			},
			want: want{
				got: `| InstanceName | SecurityGroupName |       CidrBlock       |h
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
				err: nil,
			},
		},
		{
			name: "backlog+edgecase",
			args: args{
				input:                 samples,
				format:                Backlog,
				theme:                 NoneTheme,
				hasHeader:             false,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergeFields:           []int{0, 1},
				ignoreFields:          []int{2},
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
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table, err := New(
				tt.args.input,
				WithTableFormat(tt.args.format),
				WithTableTheme(tt.args.theme),
				WithTableHeader(tt.args.hasHeader),
				WithEmptyFieldPlaceholder(tt.args.emptyFieldPlaceholder),
				WithWordDelimiter(tt.args.wordDelimiter),
				WithMergeFields(tt.args.mergeFields),
				WithIgnoreFields(tt.args.ignoreFields),
			)
			if err != nil {
				t.Fatal(err)
			}
			got, err := table.Out()
			if err != nil && err.Error() != tt.want.err.Error() {
				t.Fatalf("got: %v, want: %v", err.Error(), tt.want.err.Error())
			}
			if !reflect.DeepEqual(got, tt.want.got) {
				t.Errorf("got: %v, want: %v", got, tt.want.got)
			}
		})
	}
}

func Test_getOffset(t *testing.T) {
	type args struct {
		format    TableFormat
		hasHeader bool
	}
	type want struct {
		got int
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "markdown",
			args: args{
				format:    Markdown,
				hasHeader: true,
			},
			want: want{
				got: 2,
				err: nil,
			},
		},
		{
			name: "backlog",
			args: args{
				format:    Backlog,
				hasHeader: true,
			},
			want: want{
				got: 1,
				err: nil,
			},
		},
		{
			name: "markdown+disableHeader",
			args: args{
				format:    Markdown,
				hasHeader: false,
			},
			want: want{
				got: 0,
				err: nil,
			},
		},
		{
			name: "backlog+disableHeader",
			args: args{
				format:    Backlog,
				hasHeader: false,
			},
			want: want{
				got: 0,
				err: nil,
			},
		},
		{
			name: "default",
			args: args{
				format: "",
			},
			want: want{
				got: 0,
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getOffset(tt.args.format, tt.args.hasHeader)
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

func Test_getColor(t *testing.T) {
	type args struct {
		tableTheme TableTheme
	}
	type want struct {
		got *color.Color
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "dark",
			args: args{
				tableTheme: DarkTheme,
			},
			want: want{
				got: color.New(color.BgHiBlack, color.FgHiWhite),
				err: nil,
			},
		},
		{
			name: "light",
			args: args{
				tableTheme: LightTheme,
			},
			want: want{
				got: color.New(color.BgHiWhite, color.FgHiBlack),
				err: nil,
			},
		},
		{
			name: "none",
			args: args{
				tableTheme: NoneTheme,
			},
			want: want{
				got: color.New(color.Reset),
				err: nil,
			},
		},
		{
			name: "default",
			args: args{
				tableTheme: "",
			},
			want: want{
				got: nil,
				err: fmt.Errorf("invalid table theme detected"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getColor(tt.args.tableTheme)
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
