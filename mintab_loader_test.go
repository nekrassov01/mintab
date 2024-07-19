package mintab

import (
	"net"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestTable_Load(t *testing.T) {
	type args struct {
		input any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "interface{}",
			args: args{
				input: irregularsample,
			},
			wantErr: true,
		},
		{
			name: "ptr",
			args: args{
				input: basicsamplePtr,
			},
			wantErr: false,
		},
		{
			name: "struct",
			args: args{
				input: basicsampleNonSlice,
			},
			wantErr: false,
		},
		{
			name: "struct_empty",
			args: args{
				input: basicsampleNonSliceEmpty,
			},
			wantErr: false,
		},
		{
			name: "slice_ptr",
			args: args{
				input: basicsampleSlicePtr,
			},
			wantErr: false,
		},
		{
			name: "string",
			args: args{
				input: "aaa",
			},
			wantErr: true,
		},
		{
			name: "slice_empty",
			args: args{
				input: basicsampleEmpty,
			},
			wantErr: false,
		},
		{
			name: "slice_string",
			args: args{
				input: []string{"dummy"},
			},
			wantErr: true,
		},
		{
			name: "stringer",
			args: args{
				input: stringersample,
			},
			wantErr: false,
		},
		{
			name: "nestedStringer",
			args: args{
				input: nestedstringerSample,
			},
			wantErr: false,
		},
		{
			name: "input",
			args: args{
				input: basicinputsample,
			},
			wantErr: false,
		},
		{
			name: "input_irregular1",
			args: args{
				input: irregularinputsample1,
			},
			wantErr: true,
		},
		{
			name: "input_irregular2",
			args: args{
				input: irregularinputsample2,
			},
			wantErr: true,
		},
		{
			name: "input_irregular3",
			args: args{
				input: irregularinputsample3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Table{}
			if err := tr.Load(tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("Table.Load() error = %v, wantErr %v", err, tt.wantErr)
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
				rv: reflect.ValueOf(basicsample),
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
				rv: reflect.ValueOf(basicsample),
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
				rv: reflect.ValueOf(basicsample),
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
				rv: reflect.ValueOf(nonexportedsample),
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
				rv: basicinputsample,
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
				rv: basicinputsample,
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
				rv: basicinputsample,
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
				rv: irregularinputsample1,
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

func TestTable_setStructData(t *testing.T) {
	type fields struct {
		header                []string
		emptyFieldPlaceholder string
		wordDelimiter         string
		mergedFields          []int
		ignoredFields         []int
		colWidths             []int
		numColumns            int
		numRows               int
	}
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    [][]string
		wantErr bool
	}{
		{
			name: "text",
			fields: fields{
				header:                []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
				colWidths:             []int{0, 0, 0, 0},
				numColumns:            4,
				numRows:               6,
			},
			args: args{
				v: basicsample,
			},
			want: [][]string{
				{"i-1", "server-1", "lb-1", "tg-1"},
				{"i-2", "server-2", "lb-2\nlb-3", "tg-2"},
				{"i-3", "server-3", "lb-4", "tg-3\ntg-4"},
				{"i-4", "server-4", "-", "-"},
				{"i-5", "server-5", "lb-5", "-"},
				{"i-6", "server-6", "-", "tg-5\ntg-6\ntg-7\ntg-8"},
			},
			wantErr: false,
		},
		{
			name: "markdown",
			fields: fields{
				header:                []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				emptyFieldPlaceholder: MarkdownDefaultEmptyFieldPlaceholder,
				wordDelimiter:         MarkdownDefaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
				colWidths:             []int{0, 0, 0, 0},
				numColumns:            4,
				numRows:               6,
			},
			args: args{
				v: basicsample,
			},
			want: [][]string{
				{"i-1", "server-1", "lb-1", "tg-1"},
				{"i-2", "server-2", "lb-2<br>lb-3", "tg-2"},
				{"i-3", "server-3", "lb-4", "tg-3<br>tg-4"},
				{"i-4", "server-4", "\\-", "\\-"},
				{"i-5", "server-5", "lb-5", "\\-"},
				{"i-6", "server-6", "\\-", "tg-5<br>tg-6<br>tg-7<br>tg-8"},
			},
			wantErr: false,
		},
		{
			name: "backlog",
			fields: fields{
				header:                []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				emptyFieldPlaceholder: BacklogDefaultEmptyFieldPlaceholder,
				wordDelimiter:         BacklogDefaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
				colWidths:             []int{0, 0, 0, 0},
				numColumns:            4,
				numRows:               6,
			},
			args: args{
				v: basicsample,
			},
			want: [][]string{
				{"i-1", "server-1", "lb-1", "tg-1"},
				{"i-2", "server-2", "lb-2&br;lb-3", "tg-2"},
				{"i-3", "server-3", "lb-4", "tg-3&br;tg-4"},
				{"i-4", "server-4", "-", "-"},
				{"i-5", "server-5", "lb-5", "-"},
				{"i-6", "server-6", "-", "tg-5&br;tg-6&br;tg-7&br;tg-8"},
			},
			wantErr: false,
		},
		{
			name: "merge",
			fields: fields{
				header:                []string{"InstanceID", "InstanceName", "SecurityGroupID", "FlowDirection", "IPProtocol", "FromPort", "ToPort", "AddressType", "CidrBlock"},
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				mergedFields:          []int{0, 1, 2},
				ignoredFields:         nil,
				colWidths:             []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
				numColumns:            9,
				numRows:               8,
			},
			args: args{
				v: mergedsample,
			},
			want: [][]string{
				{"i-1", "server-1", "sg-1", "Ingress", "tcp", "22", "22", "SecurityGroup", "sg-10"},
				{"", "", "", "Egress", "-1", "0", "0", "Ipv4", "0.0.0.0/0"},
				{"", "", "sg-2", "Ingress", "tcp", "443", "443", "Ipv4", "0.0.0.0/0"},
				{"", "", "", "Egress", "-1", "0", "0", "Ipv4", "0.0.0.0/0"},
				{"i-2", "server-2", "sg-3", "Ingress", "icmp", "-1", "-1", "SecurityGroup", "sg-11"},
				{"", "", "", "Ingress", "tcp", "3389", "3389", "Ipv4", "10.1.0.0/16"},
				{"", "", "", "Ingress", "tcp", "0", "65535", "PrefixList", "pl-id/pl-name"},
				{"", "", "", "Egress", "-1", "0", "0", "Ipv4", "0.0.0.0/0"},
			},
			wantErr: false,
		},
		{
			name: "included_ptr",
			fields: fields{
				header:                []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
				colWidths:             []int{0, 0, 0, 0},
				numColumns:            4,
				numRows:               6,
			},
			args: args{
				v: basicsamplePtr,
			},
			want: [][]string{
				{"i-1", "server-1", "lb-1", "tg-1"},
				{"i-2", "server-2", "lb-2\nlb-3", "tg-2"},
				{"i-3", "server-3", "lb-4", "tg-3\ntg-4"},
				{"i-4", "server-4", "-", "-"},
				{"i-5", "server-5", "lb-5", "-"},
				{"i-6", "server-6", "-", "tg-5\ntg-6\ntg-7\ntg-8"},
			},
			wantErr: false,
		},
		{
			name: "invalid_field_name",
			fields: fields{
				header:     []string{"aaa"},
				numColumns: 4,
				numRows:    6,
			},
			args: args{
				v: basicsample,
			},
			want:    make([][]string, len(basicsample)),
			wantErr: true,
		},
		{
			name: "invalid_field",
			fields: fields{
				header:     []string{"BucketName", "Objects"},
				colWidths:  []int{0},
				numColumns: 2,
				numRows:    2,
			},
			args: args{
				v: nestedsample,
			},
			want:    make([][]string, len(nestedsample)),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Table{
				header:                tt.fields.header,
				emptyFieldPlaceholder: tt.fields.emptyFieldPlaceholder,
				wordDelimiter:         tt.fields.wordDelimiter,
				mergedFields:          tt.fields.mergedFields,
				ignoredFields:         tt.fields.ignoredFields,
				colWidths:             tt.fields.colWidths,
				numColumns:            tt.fields.numColumns,
				numRows:               tt.fields.numRows,
			}
			v := reflect.ValueOf(tt.args.v)
			if err := tr.setStructData(v); (err != nil) != tt.wantErr {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tr.data, tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", tr.data, tt.want)
			}
		})
	}
}

func TestTable_setInputData(t *testing.T) {
	type fields struct {
		header                []string
		emptyFieldPlaceholder string
		wordDelimiter         string
		mergedFields          []int
		ignoredFields         []int
		colWidths             []int
		numColumns            int
		numColumnsFirstRow    int
		numRows               int
	}
	type args struct {
		v Input
	}
	type want struct {
		header []string
		data   [][]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "text",
			fields: fields{
				header:                []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
				colWidths:             []int{0, 0, 0, 0},
				numColumns:            4,
				numColumnsFirstRow:    4,
				numRows:               6,
			},
			args: args{
				v: basicinputsample,
			},
			want: want{
				header: []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				data: [][]string{
					{"i-1", "server-1", "lb-1", "tg-1"},
					{"i-2", "server-2", "lb-2\nlb-3", "tg-2"},
					{"i-3", "server-3", "lb-4", "tg-3\ntg-4"},
					{"i-4", "server-4", "-", "-"},
					{"i-5", "server-5", "lb-5", "-"},
					{"i-6", "server-6", "-", "tg-5\ntg-6\ntg-7\ntg-8"},
				},
			},
			wantErr: false,
		},
		{
			name: "markdown",
			fields: fields{
				header:                []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				emptyFieldPlaceholder: MarkdownDefaultEmptyFieldPlaceholder,
				wordDelimiter:         MarkdownDefaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
				colWidths:             []int{0, 0, 0, 0},
				numColumns:            4,
				numColumnsFirstRow:    4,
				numRows:               6,
			},
			args: args{
				v: basicinputsample,
			},
			want: want{
				header: []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				data: [][]string{
					{"i-1", "server-1", "lb-1", "tg-1"},
					{"i-2", "server-2", "lb-2<br>lb-3", "tg-2"},
					{"i-3", "server-3", "lb-4", "tg-3<br>tg-4"},
					{"i-4", "server-4", "\\-", "\\-"},
					{"i-5", "server-5", "lb-5", "\\-"},
					{"i-6", "server-6", "\\-", "tg-5<br>tg-6<br>tg-7<br>tg-8"},
				},
			},
			wantErr: false,
		},
		{
			name: "backlog",
			fields: fields{
				header:                []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				emptyFieldPlaceholder: BacklogDefaultEmptyFieldPlaceholder,
				wordDelimiter:         BacklogDefaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
				colWidths:             []int{0, 0, 0, 0},
				numColumns:            4,
				numColumnsFirstRow:    4,
				numRows:               6,
			},
			args: args{
				v: basicinputsample,
			},
			want: want{
				header: []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				data: [][]string{
					{"i-1", "server-1", "lb-1", "tg-1"},
					{"i-2", "server-2", "lb-2&br;lb-3", "tg-2"},
					{"i-3", "server-3", "lb-4", "tg-3&br;tg-4"},
					{"i-4", "server-4", "-", "-"},
					{"i-5", "server-5", "lb-5", "-"},
					{"i-6", "server-6", "-", "tg-5&br;tg-6&br;tg-7&br;tg-8"},
				},
			},
			wantErr: false,
		},
		{
			name: "no_header",
			fields: fields{
				header:                []string{},
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
				colWidths:             []int{0, 0, 0, 0},
				numColumns:            4,
				numColumnsFirstRow:    4,
				numRows:               6,
			},
			args: args{
				v: noheadersample,
			},
			want: want{
				header: []string{},
				data: [][]string{
					{"i-1", "server-1", "lb-1", "tg-1"},
					{"i-2", "server-2", "lb-2\nlb-3", "tg-2"},
					{"i-3", "server-3", "lb-4", "tg-3\ntg-4"},
					{"i-4", "server-4", "-", "-"},
					{"i-5", "server-5", "lb-5", "-"},
					{"i-6", "server-6", "-", "tg-5\ntg-6\ntg-7\ntg-8"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid-number-of-columns",
			fields: fields{
				header:                []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
				colWidths:             []int{0, 0, 0, 0},
				numColumns:            4,
				numColumnsFirstRow:    4,
				numRows:               6,
			},
			args: args{
				v: irregularinputsample2,
			},
			want: want{
				header: []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				data: [][]string{
					{"i-1", "server-1", "lb-1", "tg-1"},
					{"i-2", "server-2", "lb-2\nlb-3", "tg-2"},
					{"i-3", "server-3", "lb-4", "tg-3\ntg-4"},
					{"i-4", "server-4", "-", "-"},
					nil,
					nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Table{
				header:                tt.fields.header,
				emptyFieldPlaceholder: tt.fields.emptyFieldPlaceholder,
				wordDelimiter:         tt.fields.wordDelimiter,
				mergedFields:          tt.fields.mergedFields,
				ignoredFields:         tt.fields.ignoredFields,
				colWidths:             tt.fields.colWidths,
				numColumns:            tt.fields.numColumns,
				numColumnsFirstRow:    tt.fields.numColumnsFirstRow,
				numRows:               tt.fields.numRows,
			}
			if err := tr.setInputData(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tr.header, tt.want.header) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", tr.header, tt.want.header)
			}
			if diff := cmp.Diff(tr.data, tt.want.data); diff != "" {
				t.Errorf(diff)
			}
			if !reflect.DeepEqual(tr.data, tt.want.data) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", tr.data, tt.want.data)
			}
		})
	}
}

func TestTable_setBorder(t *testing.T) {
	type fields struct {
		format      Format
		marginWidth int
		colWidths   []int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "text",
			fields: fields{
				format:      TextFormat,
				marginWidth: 1,
				colWidths:   []int{8, 12, 5},
			},
			want: "+----------+--------------+-------+\n",
		},
		{
			name: "markdown",
			fields: fields{
				format:      MarkdownFormat,
				marginWidth: 1,
				colWidths:   []int{8, 12, 5},
			},
			want: "|----------|--------------|-------|\n",
		},
		{
			name: "backlog",
			fields: fields{
				format:      BacklogFormat,
				marginWidth: 1,
				colWidths:   []int{8, 12, 5},
			},
			want: "|----------|--------------|-------|\n",
		},
		{
			name: "wide-margin",
			fields: fields{
				format:      TextFormat,
				marginWidth: 3,
				colWidths:   []int{8, 12, 5},
			},
			want: "+--------------+------------------+-----------+\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Table{
				format:      tt.fields.format,
				marginWidth: tt.fields.marginWidth,
				colWidths:   tt.fields.colWidths,
			}
			tr.setBorder()
			if !reflect.DeepEqual(tr.border, tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", tr.border, tt.want)
			}
		})
	}
}

func TestTable_formatStructField(t *testing.T) {
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
			got, err := tr.formatStructField(v)
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

func TestTable_formatInputField(t *testing.T) {
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
			name: "int8",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: int8(123),
			},
			want:    "123",
			wantErr: false,
		},
		{
			name: "int16",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: int16(123),
			},
			want:    "123",
			wantErr: false,
		},
		{
			name: "int32",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: int32(123),
			},
			want:    "123",
			wantErr: false,
		},
		{
			name: "int64",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: int64(123),
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
			name: "uint8",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: uint8(123),
			},
			want:    "123",
			wantErr: false,
		},
		{
			name: "uint16",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: uint16(123),
			},
			want:    "123",
			wantErr: false,
		},
		{
			name: "uint32",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: uint32(123),
			},
			want:    "123",
			wantErr: false,
		},
		{
			name: "uint64",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: uint64(123),
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
			name: "byte_slice_empty",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: []byte{},
			},
			want:    TextDefaultEmptyFieldPlaceholder,
			wantErr: false,
		},
		{
			name: "nil",
			fields: fields{
				format:                TextFormat,
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				isEscape:              false,
			},
			args: args{
				v: nil,
			},
			want:    TextDefaultEmptyFieldPlaceholder,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Table{
				format:                tt.fields.format,
				emptyFieldPlaceholder: tt.fields.emptyFieldPlaceholder,
				wordDelimiter:         tt.fields.wordDelimiter,
				isEscape:              tt.fields.isEscape,
			}
			got, err := tr.formatInputField(tt.args.v)
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
