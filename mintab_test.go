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

type escaped struct {
	Domain string
}

var (
	samples                      []sample
	samplesPtr                   []*sample
	slicePtr                     *[]sample
	nests                        []nested
	escapes                      []escaped
	irregulars                   []interface{}
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
	samplesPtr = make([]*sample, 0, len(samples))
	for i := range samples {
		samplesPtr = append(samplesPtr, &samples[i])
	}
	slicePtr = &samples
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
	escapes = []escaped{{Domain: "*.example.com"}}
	irregulars = []interface{}{sample{InstanceName: "i-1", SecurityGroupName: "sg-1", CidrBlock: []string{"10.0.0.0/16"}}, 1, "string", 2.5, struct{}{}}
}

func TestNewTable(t *testing.T) {
	type fields struct {
		format                int
		theme                 int
		hasHeader             bool
		emptyFieldPlaceholder string
		wordDelimiter         string
		mergedFields          []int
		ignoredFields         []int
		escapedTargets        []string
	}
	type want struct {
		got *Table
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "markdown+basic",
			fields: fields{
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
					mergedFields:          nil,
					ignoredFields:         nil,
					colorFlags:            nil,
					escapedTargets:        nil,
				},
			},
		},
		{
			name: "markdown+disableHeader",
			fields: fields{
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
					mergedFields:          nil,
					ignoredFields:         nil,
					colorFlags:            nil,
					escapedTargets:        nil,
				},
			},
		},
		{
			name: "markdown+merge",
			fields: fields{
				format:                MarkdownFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergedFields:          []int{0, 1},
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                MarkdownFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergedFields:          []int{0, 1},
					ignoredFields:         nil,
					colorFlags:            nil,
					escapedTargets:        nil,
				},
			},
		},
		{
			name: "markdown+ignore",
			fields: fields{
				format:                MarkdownFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				ignoredFields:         []int{2},
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                MarkdownFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergedFields:          nil,
					ignoredFields:         []int{2},
					colorFlags:            nil,
					escapedTargets:        nil,
				},
			},
		},
		{
			name: "markdown+emptyFieldPlaceholder",
			fields: fields{
				format:                MarkdownFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: "NULL",
				wordDelimiter:         defaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                MarkdownFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: "NULL",
					wordDelimiter:         defaultWordDelimiter,
					mergedFields:          nil,
					ignoredFields:         nil,
					colorFlags:            nil,
					escapedTargets:        nil,
				},
			},
		},
		{
			name: "markdown+wordDelimiter",
			fields: fields{
				format:                MarkdownFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         ",",
				mergedFields:          nil,
				ignoredFields:         nil,
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                MarkdownFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         ",",
					mergedFields:          nil,
					ignoredFields:         nil,
					colorFlags:            nil,
					escapedTargets:        nil,
				},
			},
		},
		{
			name: "markdown+escape",
			fields: fields{
				format:                MarkdownFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         ",",
				mergedFields:          nil,
				ignoredFields:         nil,
				escapedTargets:        []string{"*", "-"},
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                MarkdownFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         ",",
					mergedFields:          nil,
					ignoredFields:         nil,
					colorFlags:            nil,
					escapedTargets:        []string{"*", "-"},
				},
			},
		},
		{
			name: "markdown+edgecase",
			fields: fields{
				format:                MarkdownFormat,
				theme:                 NoneTheme,
				hasHeader:             false,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergedFields:          []int{0, 1},
				ignoredFields:         []int{2},
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                MarkdownFormat,
					theme:                 NoneTheme,
					hasHeader:             false,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergedFields:          []int{0, 1},
					ignoredFields:         []int{2},
					colorFlags:            nil,
					escapedTargets:        nil,
				},
			},
		},
		{
			name: "backlog+basic",
			fields: fields{
				format:                BacklogFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                BacklogFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergedFields:          nil,
					ignoredFields:         nil,
					colorFlags:            nil,
					escapedTargets:        nil,
				},
			},
		},
		{
			name: "backlog+disableHeader",
			fields: fields{
				format:                BacklogFormat,
				theme:                 NoneTheme,
				hasHeader:             false,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                BacklogFormat,
					theme:                 NoneTheme,
					hasHeader:             false,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergedFields:          nil,
					ignoredFields:         nil,
					colorFlags:            nil,
					escapedTargets:        nil,
				},
			},
		},
		{
			name: "backlog+merge",
			fields: fields{
				format:                BacklogFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergedFields:          []int{0, 1},
				ignoredFields:         nil,
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                BacklogFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergedFields:          []int{0, 1},
					ignoredFields:         nil,
					colorFlags:            nil,
					escapedTargets:        nil,
				},
			},
		},
		{
			name: "backlog+ignore",
			fields: fields{
				format:                BacklogFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         []int{2},
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                BacklogFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergedFields:          nil,
					ignoredFields:         []int{2},
					colorFlags:            nil,
					escapedTargets:        nil,
				},
			},
		},
		{
			name: "backlog+emptyFieldPlaceholder",
			fields: fields{
				format:                BacklogFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: "NULL",
				wordDelimiter:         defaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                BacklogFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: "NULL",
					wordDelimiter:         defaultWordDelimiter,
					mergedFields:          nil,
					ignoredFields:         nil,
					colorFlags:            nil,
					escapedTargets:        nil,
				},
			},
		},
		{
			name: "backlog+wordDelimiter",
			fields: fields{
				format:                BacklogFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         ",",
				mergedFields:          nil,
				ignoredFields:         nil,
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                BacklogFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         ",",
					mergedFields:          nil,
					ignoredFields:         nil,
					colorFlags:            nil,
					escapedTargets:        nil,
				},
			},
		},
		{
			name: "backlog+escape",
			fields: fields{
				format:                BacklogFormat,
				theme:                 NoneTheme,
				hasHeader:             true,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         ",",
				mergedFields:          nil,
				ignoredFields:         nil,
				escapedTargets:        []string{"*", "-"},
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                BacklogFormat,
					theme:                 NoneTheme,
					hasHeader:             true,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         ",",
					mergedFields:          nil,
					ignoredFields:         nil,
					colorFlags:            nil,
					escapedTargets:        []string{"*", "-"},
				},
			},
		},
		{
			name: "backlog+edgecase",
			fields: fields{
				format:                BacklogFormat,
				theme:                 NoneTheme,
				hasHeader:             false,
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergedFields:          []int{0, 1},
				ignoredFields:         []int{2},
			},
			want: want{
				got: &Table{
					data:                  nil,
					format:                BacklogFormat,
					theme:                 NoneTheme,
					hasHeader:             false,
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					mergedFields:          []int{0, 1},
					ignoredFields:         []int{2},
					colorFlags:            nil,
					escapedTargets:        nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTable(
				WithFormat(tt.fields.format),
				WithTheme(tt.fields.theme),
				WithHeader(tt.fields.hasHeader),
				WithEmptyFieldPlaceholder(tt.fields.emptyFieldPlaceholder),
				WithWordDelimiter(tt.fields.wordDelimiter),
				WithMergeFields(tt.fields.mergedFields),
				WithIgnoreFields(tt.fields.ignoredFields),
				WithEscapeTargets(tt.fields.escapedTargets),
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
		mergedFields          []int
		ignoredFields         []int
		escapedTargets        []string
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
				mergedFields:          nil,
				ignoredFields:         nil,
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
			name: "slice elements ptr",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
			},
			args: args{
				input: samplesPtr,
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
			name: "slice ptr",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
			},
			args: args{
				input: slicePtr,
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
				mergedFields:          nil,
				ignoredFields:         nil,
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
				mergedFields:          nil,
				ignoredFields:         nil,
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
				mergedFields:          []int{0, 1},
				ignoredFields:         nil,
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
					mergedFields:          []int{0, 1},
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
				mergedFields:          nil,
				ignoredFields:         []int{2},
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
					ignoredFields:         []int{2},

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
				mergedFields:          []int{0, 1},
				ignoredFields:         []int{2},
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
					mergedFields:          []int{0, 1},
					ignoredFields:         []int{2},
					colorFlags:            []bool{true, true, true, true, false, false, false, false, true, false},
				},
				err: nil,
			},
		},
		{
			name: "escape",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				escapedTargets:        []string{"*"},
			},
			args: args{
				input: escapes,
			},
			want: want{
				got: &Table{
					data: [][]string{
						{`\*.example.com`},
					},
					headers:               []string{"Domain"},
					emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
					wordDelimiter:         defaultWordDelimiter,
					colorFlags:            []bool{true},
					escapedTargets:        []string{"*"},
				},
				err: nil,
			},
		},
		{
			name: "edgecase",
			fields: fields{
				emptyFieldPlaceholder: "",
				wordDelimiter:         ",",
				mergedFields:          []int{0, 1},
				ignoredFields:         []int{2},
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
					mergedFields:          []int{0, 1},
					ignoredFields:         []int{2},
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
				err: fmt.Errorf("cannot parse input: elements of slice must be struct or pointer to struct"),
			},
		},
		{
			name: "int slice",
			args: args{
				input: []int{1, 2, 3},
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: elements of slice must be struct or pointer to struct"),
			},
		},
		{
			name: "bool slice",
			args: args{
				input: []bool{true, false, true},
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: elements of slice must be struct or pointer to struct"),
			},
		},
		{
			name: "rune slice",
			args: args{
				input: []rune{'a', 'b', 'c'},
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: elements of slice must be struct or pointer to struct"),
			},
		},
		{
			name: "string",
			args: args{
				input: "a",
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: must be a slice or a pointer to a slice"),
			},
		},
		{
			name: "int",
			args: args{
				input: 1,
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: must be a slice or a pointer to a slice"),
			},
		},
		{
			name: "bool",
			args: args{
				input: true,
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: must be a slice or a pointer to a slice"),
			},
		},
		{
			name: "rune",
			args: args{
				input: 'a',
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: must be a slice or a pointer to a slice"),
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
				err: fmt.Errorf("cannot parse input: elements of slice must be struct or pointer to struct"),
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
		{
			name: "iregular slice",
			args: args{
				input: irregulars,
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot parse input: elements of slice must not be empty interface"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{
				emptyFieldPlaceholder: tt.fields.emptyFieldPlaceholder,
				wordDelimiter:         tt.fields.wordDelimiter,
				mergedFields:          tt.fields.mergedFields,
				ignoredFields:         tt.fields.ignoredFields,
				escapedTargets:        tt.fields.escapedTargets,
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
		escapedTargets        []string
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
				mergedFields:          tt.fields.mergeFields,
				ignoredFields:         tt.fields.ignoreFields,
				colorFlags:            tt.fields.colorFlags,
				escapedTargets:        tt.fields.escapedTargets,
			}
			if got := table.Out(); !reflect.DeepEqual(got, tt.want.got) {
				t.Errorf("got: %v, want: %v", got, tt.want.got)
			}
		})
	}
}

func TestTable_formatValue(t *testing.T) {
	sp := func(s string) *string {
		return &s
	}
	type fields struct {
		emptyFieldPlaceholder string
		wordDelimiter         string
		escapedTargets        []string
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
			name: "string",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: "aaa",
			},
			want: want{
				got: "aaa",
				err: nil,
			},
		},
		{
			name: "escape",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
				escapedTargets:        []string{"*"},
			},
			args: args{
				v: "*.example.com",
			},
			want: want{
				got: `\*.example.com`,
				err: nil,
			},
		},
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
			name: "int",
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
			name: "uint",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: uint(123),
			},
			want: want{
				got: "123",
				err: nil,
			},
		},
		{
			name: "float",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: 123.456,
			},
			want: want{
				got: "123.456",
				err: nil,
			},
		},
		{
			name: "non-nil pointer string",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: new(string),
			},
			want: want{
				got: defaultEmptyFieldPlaceholder,
				err: nil,
			},
		},
		{
			name: "non-nil pointer int",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: new(int),
			},
			want: want{
				got: "0",
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
			name: "slice with empty strings",
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
			name: "int slice",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: []int{0, 1, 2},
			},
			want: want{
				got: "0" + defaultWordDelimiter + "1" + defaultWordDelimiter + "2",
				err: nil,
			},
		},
		{
			name: "uint slice",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: []uint{0, 1, 2},
			},
			want: want{
				got: "0" + defaultWordDelimiter + "1" + defaultWordDelimiter + "2",
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
		{
			name: "pointer to slice with normal strings",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: &[]string{"a", "b"},
			},
			want: want{
				got: "a" + defaultWordDelimiter + "b",
				err: nil,
			},
		},
		{
			name: "slice with pointer to strings",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: []*string{sp(""), sp("a"), sp("b")},
			},
			want: want{
				got: defaultEmptyFieldPlaceholder + defaultWordDelimiter + "a" + defaultWordDelimiter + "b",
				err: nil,
			},
		},
		{
			name: "slice with pointer to empty string",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: []*string{},
			},
			want: want{
				got: defaultEmptyFieldPlaceholder,
				err: nil,
			},
		},
		{
			name: "slice with nil pointer",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: []*int{nil},
			},
			want: want{
				got: defaultEmptyFieldPlaceholder,
				err: nil,
			},
		},
		{
			name: "slice with byte slice",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: []byte("string"),
			},
			want: want{
				got: "string",
				err: nil,
			},
		},
		{
			name: "nil ptr",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: (*int)(nil),
			},
			want: want{
				got: defaultEmptyFieldPlaceholder,
				err: nil,
			},
		},
		{
			name: "escape",
			fields: fields{
				emptyFieldPlaceholder: defaultEmptyFieldPlaceholder,
				wordDelimiter:         defaultWordDelimiter,
			},
			args: args{
				v: `aaa
  bbb
    ccc
`,
			},
			want: want{
				got: "aaa<br>&nbsp;&nbsp;bbb<br>&nbsp;&nbsp;&nbsp;&nbsp;ccc<br>",
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.args.v)
			table := &Table{
				emptyFieldPlaceholder: tt.fields.emptyFieldPlaceholder,
				wordDelimiter:         tt.fields.wordDelimiter,
				escapedTargets:        tt.fields.escapedTargets,
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

func TestTable_setHeader(t *testing.T) {
	type testHeader struct {
		ExportedString   string
		ExportedInt      int
		unexportedString string //nolint
		unexportedInt    int    //nolint
	}
	type fields struct {
		headers       []string
		ignoredFields []int
	}
	type args struct {
		typ reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "no initial headers",
			fields: fields{
				headers:       []string{},
				ignoredFields: []int{},
			},
			args: args{
				typ: reflect.TypeOf(testHeader{}),
			},
			want: []string{"ExportedString", "ExportedInt"},
		},
		{
			name: "ignoring fields",
			fields: fields{
				headers:       []string{},
				ignoredFields: []int{0},
			},
			args: args{
				typ: reflect.TypeOf(testHeader{}),
			},
			want: []string{"ExportedInt"},
		},
		{
			name: "with pre-existing headers",
			fields: fields{
				headers:       []string{"ExportedString", "ExportedInt"},
				ignoredFields: []int{},
			},
			args: args{
				typ: reflect.TypeOf(testHeader{}),
			},
			want: []string{"ExportedString", "ExportedInt"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{
				headers:       tt.fields.headers,
				ignoredFields: tt.fields.ignoredFields,
			}
			table.setHeader(tt.args.typ)
			if !reflect.DeepEqual(table.headers, tt.want) {
				t.Errorf("got: %v, want: %v", table.headers, tt.want)
			}
		})
	}
}

func TestTable_setData(t *testing.T) {
	type fields struct {
		headers []string
	}
	type args struct {
		v any
	}
	type want struct {
		got []bool
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "invalid field name",
			fields: fields{
				headers: []string{"aaa"},
			},
			args: args{
				v: samples,
			},
			want: want{
				got: nil,
				err: fmt.Errorf("field \"aaa\" does not exist"),
			},
		},
		{
			name: "invalid field",
			fields: fields{
				headers: []string{"BucketName", "Objects"},
			},
			args: args{
				v: nests,
			},
			want: want{
				got: nil,
				err: fmt.Errorf("cannot format field \"Objects\": elements of slice must not be nested"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.args.v)
			table := &Table{
				headers: tt.fields.headers,
			}
			if err := table.setData(v); err != nil {
				if err.Error() != tt.want.err.Error() {
					t.Fatalf("got: %v, want: %v", err.Error(), tt.want.err.Error())
				}
				return
			}
			if !reflect.DeepEqual(table.data, tt.want.got) {
				t.Errorf("got: %v, want: %v", table.data, tt.want.got)
			}
		})
	}
}

func TestTable_setColorFlags(t *testing.T) {
	type args struct {
		v any
	}
	type want struct {
		got []bool
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "ptr",
			args: args{
				v: slicePtr,
			},
			want: want{
				got: []bool{true, true, true, true, false, false, false, false, true, false},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.args.v)
			table := &Table{}
			table.setColorFlags(v)
			if !reflect.DeepEqual(table.colorFlags, tt.want.got) {
				t.Errorf("got: %v, want: %v", table.colorFlags, tt.want.got)
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
			name: "none",
			fields: fields{
				theme: NoneTheme,
			},
			want: want{
				got: &color.Color{},
			},
		},
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
			name: "default",
			fields: fields{
				theme: 9,
			},
			want: want{
				got: &color.Color{},
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
