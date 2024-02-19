package mintab

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"
)

type basicSample struct {
	InstanceID   string
	InstanceName string
	AttachedLB   []string
	AttachedTG   []string
}

type object struct {
	ObjectID   int
	ObjectName string
}

type nestedSample struct {
	BucketName string
	Objects    []object
}

type mergedSample struct {
	InstanceID      string
	InstanceName    string
	VPCID           string
	SecurityGroupID string
	FlowDirection   string
	IPProtocol      string
	FromPort        int
	ToPort          int
	AddressType     string
	CidrBlock       string
}

var (
	basicsample         []basicSample
	basicsampleEmpty    []basicSample
	nestedsample        []nestedSample
	mergedsample        []mergedSample
	basicsamplePtr      []*basicSample
	basicsampleSlicePtr *[]basicSample
	irregularsample     []interface{}
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	basicsample = []basicSample{
		{InstanceID: "i-1", InstanceName: "server-1", AttachedLB: []string{"lb-1"}, AttachedTG: []string{"tg-1"}},
		{InstanceID: "i-2", InstanceName: "server-2", AttachedLB: []string{"lb-2", "lb-3"}, AttachedTG: []string{"tg-2"}},
		{InstanceID: "i-3", InstanceName: "server-3", AttachedLB: []string{"lb-4"}, AttachedTG: []string{"tg-3", "tg-4"}},
		{InstanceID: "i-4", InstanceName: "server-4", AttachedLB: []string{}, AttachedTG: []string{}},
		{InstanceID: "i-5", InstanceName: "server-5", AttachedLB: []string{"lb-5"}, AttachedTG: []string{}},
		{InstanceID: "i-6", InstanceName: "server-6", AttachedLB: []string{}, AttachedTG: []string{"tg-5", "tg-6", "tg-7", "tg-8"}},
	}
	basicsampleEmpty = []basicSample{
		{},
	}
	nestedsample = []nestedSample{
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
	mergedsample = []mergedSample{
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-1", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 22, ToPort: 22, AddressType: "SecurityGroup", CidrBlock: "sg-10"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-1", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-2", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 443, ToPort: 443, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-2", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "icmp", FromPort: -1, ToPort: -1, AddressType: "SecurityGroup", CidrBlock: "sg-11"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 3389, ToPort: 3389, AddressType: "Ipv4", CidrBlock: "10.1.0.0/16"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 0, ToPort: 65535, AddressType: "PrefixList", CidrBlock: "pl-id/pl-name"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
	}
	basicsamplePtr = make([]*basicSample, 0, len(basicsample))
	for i := range basicsample {
		basicsamplePtr = append(basicsamplePtr, &basicsample[i])
	}
	basicsampleSlicePtr = &basicsample
	irregularsample = []interface{}{
		basicSample{
			InstanceID:   "i-1",
			InstanceName: "server-1",
			AttachedLB:   []string{"lb-1"},
			AttachedTG:   []string{"tg-1"},
		},
		1,
		"string",
		2.5,
		struct{}{},
	}
}

func TestFormat_String(t *testing.T) {
	tests := []struct {
		name string
		o    Format
		want string
	}{
		{
			name: "text",
			o:    FormatText,
			want: "text",
		},
		{
			name: "markdown",
			o:    FormatMarkdown,
			want: "markdown",
		},
		{
			name: "backlog",
			o:    FormatBacklog,
			want: "backlog",
		},
		{
			name: "other",
			o:    3,
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

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name  string
		args  args
		want  *Table
		wantW string
	}{
		{
			name: "default",
			args: args{
				opts: []Option{},
			},
			want: &Table{
				writer:                &bytes.Buffer{},
				data:                  nil,
				header:                nil,
				format:                FormatText,
				border:                "",
				margin:                1,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
				columnWidths:          nil,
				hasHeader:             true,
				hasEscape:             false,
				compress:              false,
			},
		},
		{
			name: "not-default",
			args: args{
				opts: []Option{
					WithFormat(FormatMarkdown),
					WithHeader(false),
					WithMargin(2),
					WithEmptyFieldPlaceholder(MarkdownDefaultEmptyFieldPlaceholder),
					WithWordDelimiter(MarkdownDefaultWordDelimiter),
					WithMergeFields([]int{0}),
					WithIgnoreFields([]int{0}),
					WithEscape(true),
					WithCompress(true),
				},
			},
			want: &Table{
				writer:                &bytes.Buffer{},
				data:                  nil,
				header:                nil,
				format:                FormatMarkdown,
				border:                "",
				margin:                2,
				emptyFieldPlaceholder: MarkdownDefaultEmptyFieldPlaceholder,
				wordDelimiter:         MarkdownDefaultWordDelimiter,
				mergedFields:          []int{0},
				ignoredFields:         []int{0},
				columnWidths:          nil,
				hasHeader:             false,
				hasEscape:             true,
				compress:              true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			var got *Table
			if got = New(w, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot\n%v\nwant\n%v\n", got, tt.want)
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("\ngot\n%v\nwant\n%v\n", gotW, tt.wantW)
			}
		})
	}
}

func TestTable_Load(t *testing.T) {
	type fields struct {
		header []string
		margin int
	}
	type args struct {
		input any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "err_margin",
			fields: fields{
				margin: -1,
			},
			args: args{
				input: basicsample,
			},
			wantErr: true,
		},
		{
			name:   "err_interface{}",
			fields: fields{},
			args: args{
				input: irregularsample,
			},
			wantErr: true,
		},
		{
			name:   "err_ptr",
			fields: fields{},
			args: args{
				input: basicsamplePtr,
			},
			wantErr: false,
		},
		{
			name:   "err_slice_ptr",
			fields: fields{},
			args: args{
				input: basicsampleSlicePtr,
			},
			wantErr: false,
		},
		{
			name:   "err_slice_ptr",
			fields: fields{},
			args: args{
				input: "aaa",
			},
			wantErr: true,
		},
		{
			name:   "err_slice_ptr",
			fields: fields{},
			args: args{
				input: []string{},
			},
			wantErr: true,
		},
		{
			name:   "err_struct_in_slice",
			fields: fields{},
			args: args{
				input: []string{"dummy"},
			},
			wantErr: true,
		},
		{
			name:   "err_irregular",
			fields: fields{},
			args: args{
				input: irregularsample,
			},
			wantErr: true,
		},
		{
			name: "err_irregular",
			fields: fields{
				header: []string{"dummmy"},
			},
			args: args{
				input: basicsample,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Table{
				header: tt.fields.header,
				margin: tt.fields.margin,
			}
			if err := tr.Load(tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("Table.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTable_Out(t *testing.T) {
	type fields struct {
		format       Format
		header       []string
		data         [][]string
		columnWidths []int
		compress     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "text",
			fields: fields{
				format: FormatText,
				header: []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				data: [][]string{
					{"i-1", "server-1", "lb-1", "tg-1"},
					{"i-2", "server-2", "lb-2\nlb-3", "tg-2"},
					{"i-3", "server-3", "lb-4", "tg-3\ntg-4"},
					{"i-4", "server-4", "-", "-"},
					{"i-5", "server-5", "lb-5", "-"},
					{"i-6", "server-6", "-", "tg-5\ntg-6\ntg-7\ntg-8"},
				},
				columnWidths: []int{10, 12, 10, 10},
				compress:     false,
			},
			want: `+------------+--------------+------------+------------+
| InstanceID | InstanceName | AttachedLB | AttachedTG |
+------------+--------------+------------+------------+
| i-1        | server-1     | lb-1       | tg-1       |
+------------+--------------+------------+------------+
| i-2        | server-2     | lb-2       | tg-2       |
|            |              | lb-3       |            |
+------------+--------------+------------+------------+
| i-3        | server-3     | lb-4       | tg-3       |
|            |              |            | tg-4       |
+------------+--------------+------------+------------+
| i-4        | server-4     | -          | -          |
+------------+--------------+------------+------------+
| i-5        | server-5     | lb-5       | -          |
+------------+--------------+------------+------------+
| i-6        | server-6     | -          | tg-5       |
|            |              |            | tg-6       |
|            |              |            | tg-7       |
|            |              |            | tg-8       |
+------------+--------------+------------+------------+

`,
		},
		{
			name: "text_with_compressed",
			fields: fields{
				format: FormatText,
				header: []string{"InstanceID", "InstanceName", "VPCID", "SecurityGroupID", "FlowDirection", "IPProtocol", "FromPort", "ToPort", "AddressType", "CidrBlock"},
				data: [][]string{
					{"i-1", "server-1", "vpc-1", "sg-1", "Ingress", "tcp", "22", "22", "SecurityGroup", "sg-10"},
					{"", "", "", "", "Egress", "-1", "0", "0", "Ipv4", "0.0.0.0/0"},
					{"", "", "", "sg-2", "Ingress", "tcp", "443", "443", "Ipv4", "0.0.0.0/0"},
					{"", "", "", "", "Egress", "-1", "0", "0", "Ipv4", "0.0.0.0/0"},
					{"i-2", "server-2", "vpc-1", "sg-3", "Ingress", "icmp", "-1", "-1", "SecurityGroup", "sg-11"},
					{"", "", "", "", "Ingress", "tcp", "3389", "3389", "Ipv4", "10.1.0.0/16"},
					{"", "", "", "", "Ingress", "tcp", "0", "65535", "PrefixList", "pl-id/pl-name"},
					{"", "", "", "", "Egress", "-1", "0", "0", "Ipv4", "0.0.0.0/0"},
				},
				columnWidths: []int{10, 12, 5, 15, 13, 10, 8, 6, 13, 13},
				compress:     true,
			},
			want: `+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| InstanceID | InstanceName | VPCID | SecurityGroupID | FlowDirection | IPProtocol | FromPort | ToPort | AddressType   | CidrBlock     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-1        | server-1     | vpc-1 | sg-1            | Ingress       | tcp        |       22 |     22 | SecurityGroup | sg-10         |
|            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
|            |              |       | sg-2            | Ingress       | tcp        |      443 |    443 | Ipv4          | 0.0.0.0/0     |
|            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-2        | server-2     | vpc-1 | sg-3            | Ingress       | icmp       |       -1 |     -1 | SecurityGroup | sg-11         |
|            |              |       |                 | Ingress       | tcp        |     3389 |   3389 | Ipv4          | 10.1.0.0/16   |
|            |              |       |                 | Ingress       | tcp        |        0 |  65535 | PrefixList    | pl-id/pl-name |
|            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+

`,
		},
		{
			name: "markdown",
			fields: fields{
				format: FormatMarkdown,
				header: []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				data: [][]string{
					{"i-1", "server-1", "lb-1", "tg-1"},
					{"i-2", "server-2", "lb-2<br>lb-3", "tg-2"},
					{"i-3", "server-3", "lb-4", "tg-3<br>tg-4"},
					{"i-4", "server-4", "\\-", "\\-"},
					{"i-5", "server-5", "lb-5", "\\-"},
					{"i-6", "server-6", "\\-", "tg-5<br>tg-6<br>tg-7<br>tg-8"},
				},
				columnWidths: []int{10, 12, 12, 28},
				compress:     false,
			},
			want: `| InstanceID | InstanceName | AttachedLB   | AttachedTG                   |
|------------|--------------|--------------|------------------------------|
| i-1        | server-1     | lb-1         | tg-1                         |
| i-2        | server-2     | lb-2<br>lb-3 | tg-2                         |
| i-3        | server-3     | lb-4         | tg-3<br>tg-4                 |
| i-4        | server-4     | \-           | \-                           |
| i-5        | server-5     | lb-5         | \-                           |
| i-6        | server-6     | \-           | tg-5<br>tg-6<br>tg-7<br>tg-8 |

`,
		},
		{
			name: "backlog",
			fields: fields{
				format: FormatBacklog,
				header: []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				data: [][]string{
					{"i-1", "server-1", "lb-1", "tg-1"},
					{"i-2", "server-2", "lb-2&br;lb-3", "tg-2"},
					{"i-3", "server-3", "lb-4", "tg-3&br;tg-4"},
					{"i-4", "server-4", "-", "-"},
					{"i-5", "server-5", "lb-5", "-"},
					{"i-6", "server-6", "-", "tg-5&br;tg-6&br;tg-7&br;tg-8"},
				},
				columnWidths: []int{10, 12, 12, 28},
				compress:     false,
			},
			want: `| InstanceID | InstanceName | AttachedLB   | AttachedTG                   |h
| i-1        | server-1     | lb-1         | tg-1                         |
| i-2        | server-2     | lb-2&br;lb-3 | tg-2                         |
| i-3        | server-3     | lb-4         | tg-3&br;tg-4                 |
| i-4        | server-4     | -            | -                            |
| i-5        | server-5     | lb-5         | -                            |
| i-6        | server-6     | -            | tg-5&br;tg-6&br;tg-7&br;tg-8 |

`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			tr := New(buf)
			tr.format = tt.fields.format
			tr.header = tt.fields.header
			tr.data = tt.fields.data
			tr.columnWidths = tt.fields.columnWidths
			tr.compress = tt.fields.compress
			tr.setBorder()
			tr.Out()
			if !reflect.DeepEqual(buf.String(), tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", buf.String(), tt.want)
			}
		})
	}
}

func TestTable_printHeader(t *testing.T) {
	type fields struct {
		header       []string
		format       Format
		margin       int
		columnWidths []int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "text",
			fields: fields{
				header:       []string{"a", "bb", "ccc"},
				format:       FormatText,
				margin:       1,
				columnWidths: []int{1, 2, 3},
			},
			want: "| a | bb | ccc |\n",
		},
		{
			name: "markdown",
			fields: fields{
				header:       []string{"a", "bb", "ccc"},
				format:       FormatMarkdown,
				margin:       1,
				columnWidths: []int{1, 2, 3},
			},
			want: "| a | bb | ccc |\n",
		},
		{
			name: "backlog",
			fields: fields{
				header:       []string{"a", "bb", "ccc"},
				format:       FormatBacklog,
				margin:       1,
				columnWidths: []int{1, 2, 3},
			},
			want: "| a | bb | ccc |h\n",
		},
		{
			name: "margin",
			fields: fields{
				header:       []string{"a", "bb", "ccc"},
				format:       FormatText,
				margin:       3,
				columnWidths: []int{1, 2, 3},
			},
			want: "|   a   |   bb   |   ccc   |\n",
		},
		{
			name: "long",
			fields: fields{
				header:       []string{"a", "bb", "ccc"},
				format:       FormatText,
				margin:       1,
				columnWidths: []int{10, 2, 3},
			},
			want: "| a          | bb | ccc |\n",
		},
		{
			name: "short",
			fields: fields{
				header:       []string{"a", "bb", "ccc"},
				format:       FormatText,
				margin:       1,
				columnWidths: []int{1, 2, 1},
			},
			want: "| a | bb | ccc |\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			tr := New(buf, WithMargin(tt.fields.margin))
			tr.format = tt.fields.format
			tr.header = tt.fields.header
			tr.margin = tt.fields.margin
			tr.columnWidths = tt.fields.columnWidths
			tr.printHeader()
			if !reflect.DeepEqual(buf.String(), tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", buf.String(), tt.want)
			}
		})
	}
}

func TestTable_printData(t *testing.T) {
	type fields struct {
		data         [][]string
		format       Format
		columnWidths []int
		compress     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "text",
			fields: fields{
				data: [][]string{
					{"i-1", "server-1", "lb-1", "tg-1"},
					{"i-2", "server-2", "lb-2\nlb-3", "tg-2"},
					{"i-3", "server-3", "lb-4", "tg-3\ntg-4"},
					{"i-4", "server-4", "-", "-"},
					{"i-5", "server-5", "lb-5", "-"},
					{"i-6", "server-6", "-", "tg-5\ntg-6\ntg-7\ntg-8"},
				},
				format:       FormatText,
				columnWidths: []int{10, 12, 10, 10},
				compress:     false,
			},
			want: `| i-1        | server-1     | lb-1       | tg-1       |
+------------+--------------+------------+------------+
| i-2        | server-2     | lb-2       | tg-2       |
|            |              | lb-3       |            |
+------------+--------------+------------+------------+
| i-3        | server-3     | lb-4       | tg-3       |
|            |              |            | tg-4       |
+------------+--------------+------------+------------+
| i-4        | server-4     | -          | -          |
+------------+--------------+------------+------------+
| i-5        | server-5     | lb-5       | -          |
+------------+--------------+------------+------------+
| i-6        | server-6     | -          | tg-5       |
|            |              |            | tg-6       |
|            |              |            | tg-7       |
|            |              |            | tg-8       |
`,
		},
		{
			name: "markdown",
			fields: fields{
				data: [][]string{
					{"i-1", "server-1", "lb-1", "tg-1"},
					{"i-2", "server-2", "lb-2<br>lb-3", "tg-2"},
					{"i-3", "server-3", "lb-4", "tg-3<br>tg-4"},
					{"i-4", "server-4", "\\-", "\\-"},
					{"i-5", "server-5", "lb-5", "\\-"},
					{"i-6", "server-6", "\\-", "tg-5<br>tg-6<br>tg-7<br>tg-8"},
				},
				format:       FormatMarkdown,
				columnWidths: []int{10, 12, 12, 28},
				compress:     false,
			},
			want: `| i-1        | server-1     | lb-1         | tg-1                         |
| i-2        | server-2     | lb-2<br>lb-3 | tg-2                         |
| i-3        | server-3     | lb-4         | tg-3<br>tg-4                 |
| i-4        | server-4     | \-           | \-                           |
| i-5        | server-5     | lb-5         | \-                           |
| i-6        | server-6     | \-           | tg-5<br>tg-6<br>tg-7<br>tg-8 |
`,
		},
		{
			name: "backlog",
			fields: fields{
				data: [][]string{
					{"i-1", "server-1", "lb-1", "tg-1"},
					{"i-2", "server-2", "lb-2&br;lb-3", "tg-2"},
					{"i-3", "server-3", "lb-4", "tg-3&br;tg-4"},
					{"i-4", "server-4", "-", "-"},
					{"i-5", "server-5", "lb-5", "-"},
					{"i-6", "server-6", "-", "tg-5&br;tg-6&br;tg-7&br;tg-8"},
				},
				format:       FormatBacklog,
				columnWidths: []int{10, 12, 12, 28},
				compress:     false,
			},
			want: `| i-1        | server-1     | lb-1         | tg-1                         |
| i-2        | server-2     | lb-2&br;lb-3 | tg-2                         |
| i-3        | server-3     | lb-4         | tg-3&br;tg-4                 |
| i-4        | server-4     | -            | -                            |
| i-5        | server-5     | lb-5         | -                            |
| i-6        | server-6     | -            | tg-5&br;tg-6&br;tg-7&br;tg-8 |
`,
		},
		{
			name: "text_with_compress",
			fields: fields{
				data: [][]string{
					{"i-1", "server-1", "vpc-1", "sg-1", "Ingress", "tcp", "22", "22", "SecurityGroup", "sg-10"},
					{"", "", "", "", "Egress", "-1", "0", "0", "Ipv4", "0.0.0.0/0"},
					{"", "", "", "sg-2", "Ingress", "tcp", "443", "443", "Ipv4", "0.0.0.0/0"},
					{"", "", "", "", "Egress", "-1", "0", "0", "Ipv4", "0.0.0.0/0"},
					{"i-2", "server-2", "vpc-1", "sg-3", "Ingress", "icmp", "-1", "-1", "SecurityGroup", "sg-11"},
					{"", "", "", "", "Ingress", "tcp", "3389", "3389", "Ipv4", "10.1.0.0/16"},
					{"", "", "", "", "Ingress", "tcp", "0", "65535", "PrefixList", "pl-id/pl-name"},
					{"", "", "", "", "Egress", "-1", "0", "0", "Ipv4", "0.0.0.0/0"},
				},
				format:       FormatText,
				columnWidths: []int{10, 12, 5, 15, 13, 10, 8, 6, 13, 13},
				compress:     true,
			},
			want: `| i-1        | server-1     | vpc-1 | sg-1            | Ingress       | tcp        |       22 |     22 | SecurityGroup | sg-10         |
|            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
|            |              |       | sg-2            | Ingress       | tcp        |      443 |    443 | Ipv4          | 0.0.0.0/0     |
|            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-2        | server-2     | vpc-1 | sg-3            | Ingress       | icmp       |       -1 |     -1 | SecurityGroup | sg-11         |
|            |              |       |                 | Ingress       | tcp        |     3389 |   3389 | Ipv4          | 10.1.0.0/16   |
|            |              |       |                 | Ingress       | tcp        |        0 |  65535 | PrefixList    | pl-id/pl-name |
|            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			tr := New(buf)
			tr.data = tt.fields.data
			tr.format = tt.fields.format
			tr.columnWidths = tt.fields.columnWidths
			tr.compress = tt.fields.compress
			tr.setBorder()
			tr.printData()
			if !reflect.DeepEqual(buf.String(), tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", buf.String(), tt.want)
			}
		})
	}
}

func TestTable_printDataBorder(t *testing.T) {
	type fields struct {
		margin       int
		columnWidths []int
	}
	type args struct {
		row []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "basic",
			fields: fields{
				margin:       1,
				columnWidths: []int{4, 5, 6},
			},
			args: args{
				row: []string{"a", "bb", "ccc"},
			},
			want: "+------+-------+--------+\n",
		},
		{
			name: "margin",
			fields: fields{
				margin:       3,
				columnWidths: []int{4, 5, 6},
			},
			args: args{
				row: []string{"a", "bb", "ccc"},
			},
			want: "+----------+-----------+------------+\n",
		},
		{
			name: "long",
			fields: fields{
				margin:       1,
				columnWidths: []int{10, 5, 6},
			},
			args: args{
				row: []string{"a", "bb", "ccc"},
			},
			want: "+------------+-------+--------+\n",
		},
		{
			name: "short",
			fields: fields{
				margin:       1,
				columnWidths: []int{4, 5, 6},
			},
			args: args{
				row: []string{"aaaaaa", "bb", "ccc"},
			},
			want: "+------+-------+--------+\n",
		},
		{
			name: "empty_field_included_1",
			fields: fields{
				margin:       1,
				columnWidths: []int{4, 5, 6},
			},
			args: args{
				row: []string{"", "bb", "ccc"},
			},
			want: "+      +-------+--------+\n",
		},
		{
			name: "empty_field_included_2",
			fields: fields{
				margin:       1,
				columnWidths: []int{4, 5, 6},
			},
			args: args{
				row: []string{"", "", "ccc"},
			},
			want: "+      +       +--------+\n",
		},
		{
			name: "empty_field_included_3",
			fields: fields{
				margin:       1,
				columnWidths: []int{4, 5, 6},
			},
			args: args{
				row: []string{"", "", ""},
			},
			want: "+      +       +        +\n",
		},
		{
			name: "empty_field_included_4",
			fields: fields{
				margin:       1,
				columnWidths: []int{4, 5, 6},
			},
			args: args{
				row: []string{"", "bb", ""},
			},
			want: "+      +-------+        +\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			tr := New(buf)
			tr.margin = tt.fields.margin
			tr.columnWidths = tt.fields.columnWidths
			tr.printDataBorder(tt.args.row)
			if !reflect.DeepEqual(buf.String(), tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", buf.String(), tt.want)
			}
		})
	}
}

func TestTable_printBorder(t *testing.T) {
	type fields struct {
		format       Format
		margin       int
		columnWidths []int
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "text",
			fields: fields{
				format:       FormatText,
				margin:       1,
				columnWidths: []int{8, 12, 5},
			},
			want: "+----------+--------------+-------+\n",
		},
		{
			name: "markdown",
			fields: fields{
				format:       FormatMarkdown,
				margin:       1,
				columnWidths: []int{8, 12, 5},
			},
			want: "|----------|--------------|-------|\n",
		},
		{
			name: "backlog",
			fields: fields{
				format:       FormatBacklog,
				margin:       1,
				columnWidths: []int{8, 12, 5},
			},
			want: "|----------|--------------|-------|\n",
		},
		{
			name: "wide-margin",
			fields: fields{
				format:       FormatText,
				margin:       3,
				columnWidths: []int{8, 12, 5},
			},
			want: "+--------------+------------------+-----------+\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			tr := New(buf)
			tr.format = tt.fields.format
			tr.margin = tt.fields.margin
			tr.columnWidths = tt.fields.columnWidths
			tr.setBorder()
			tr.printBorder()
			if !reflect.DeepEqual(buf.String(), tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", buf.String(), tt.want)
			}
		})
	}
}

func TestTable_setAttr(t *testing.T) {
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
				format: FormatText,
			},
			want: want{
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
			},
		},
		{
			name: "markdown",
			fields: fields{
				format: FormatMarkdown,
			},
			want: want{
				emptyFieldPlaceholder: MarkdownDefaultEmptyFieldPlaceholder,
				wordDelimiter:         MarkdownDefaultWordDelimiter,
			},
		},
		{
			name: "backlog",
			fields: fields{
				format: FormatBacklog,
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
			tr.setAttr()
			if !reflect.DeepEqual(tr.emptyFieldPlaceholder, tt.want.emptyFieldPlaceholder) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", tr.emptyFieldPlaceholder, tt.want.emptyFieldPlaceholder)
			}
			if !reflect.DeepEqual(tr.wordDelimiter, tt.want.wordDelimiter) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", tr.wordDelimiter, tt.want.wordDelimiter)
			}
		})
	}
}

func TestTable_setHeader(t *testing.T) {
	type fields struct {
		ignoredFields []int
	}
	type args struct {
		typ reflect.Type
	}
	type want struct {
		header       []string
		columnWidths []int
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
				ignoredFields: nil,
			},
			want: want{
				header:       []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				columnWidths: []int{10, 12, 10, 10},
			},
		},
		{
			name: "ignore",
			fields: fields{
				ignoredFields: []int{1},
			},
			want: want{
				header:       []string{"InstanceID", "AttachedLB", "AttachedTG"},
				columnWidths: []int{10, 10, 10},
			},
		},
		{
			name: "ignore-signed-int",
			fields: fields{
				ignoredFields: []int{-10},
			},
			want: want{
				header:       []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				columnWidths: []int{10, 12, 10, 10},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Table{
				ignoredFields: tt.fields.ignoredFields,
			}
			if err := tr.Load(basicsampleEmpty); err != nil {
				t.Fatal(err)
			}
			tr.setHeader(tt.args.typ)
			if !reflect.DeepEqual(tr.header, tt.want.header) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", tr.header, tt.want.header)
			}
			if !reflect.DeepEqual(tr.columnWidths, tt.want.columnWidths) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", tr.columnWidths, tt.want.columnWidths)
			}
		})
	}
}

func TestTable_setData(t *testing.T) {
	type fields struct {
		header                []string
		emptyFieldPlaceholder string
		wordDelimiter         string
		mergedFields          []int
		ignoredFields         []int
		columnWidths          []int
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
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
				columnWidths:          []int{0, 0, 0, 0},
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
				columnWidths:          []int{0, 0, 0, 0},
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
				columnWidths:          []int{0, 0, 0, 0},
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
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				mergedFields:          []int{0, 1, 2},
				ignoredFields:         nil,
				columnWidths:          []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
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
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
				columnWidths:          []int{0, 0, 0, 0},
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
			name: "slice_ptr",
			fields: fields{
				header:                []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
				columnWidths:          []int{0, 0, 0, 0},
			},
			args: args{
				v: basicsampleSlicePtr,
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
				header: []string{"aaa"},
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
				header:       []string{"BucketName", "Objects"},
				columnWidths: []int{0},
			},
			args: args{
				v: nestedsample,
			},
			want:    make([][]string, len(nestedsample)),
			wantErr: true,
		},
		{
			name: "not_slice",
			fields: fields{
				header:       []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				columnWidths: []int{0, 0, 0, 0},
			},
			args: args{
				v: "aaa",
			},
			want:    nil,
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
				columnWidths:          tt.fields.columnWidths,
			}
			v := reflect.ValueOf(tt.args.v)
			if err := tr.setData(v); (err != nil) != tt.wantErr {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tr.data, tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", tr.data, tt.want)
			}
		})
	}
}

func TestTable_setBorder(t *testing.T) {
	type fields struct {
		format       Format
		margin       int
		columnWidths []int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "text",
			fields: fields{
				format:       FormatText,
				margin:       1,
				columnWidths: []int{8, 12, 5},
			},
			want: "+----------+--------------+-------+",
		},
		{
			name: "markdown",
			fields: fields{
				format:       FormatMarkdown,
				margin:       1,
				columnWidths: []int{8, 12, 5},
			},
			want: "|----------|--------------|-------|",
		},
		{
			name: "backlog",
			fields: fields{
				format:       FormatBacklog,
				margin:       1,
				columnWidths: []int{8, 12, 5},
			},
			want: "|----------|--------------|-------|",
		},
		{
			name: "wide-margin",
			fields: fields{
				format:       FormatText,
				margin:       3,
				columnWidths: []int{8, 12, 5},
			},
			want: "+--------------+------------------+-----------+",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Table{
				format:       tt.fields.format,
				margin:       tt.fields.margin,
				columnWidths: tt.fields.columnWidths,
			}
			tr.setBorder()
			if !reflect.DeepEqual(tr.border, tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", tr.border, tt.want)
			}
		})
	}
}

func TestTable_getMargin(t *testing.T) {
	type fields struct {
		margin int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "basic",
			fields: fields{
				margin: 1,
			},
			want: " ",
		},
		{
			name: "basic",
			fields: fields{
				margin: 2,
			},
			want: "  ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Table{
				margin: tt.fields.margin,
			}
			got := tr.getMargin()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, tt.want)
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
		hasEscape             bool
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
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
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
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
			},
			args: args{
				v: "",
			},
			want:    DefaultEmptyFieldPlaceholder,
			wantErr: false,
		},
		{
			name: "byte_slice",
			fields: fields{
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
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
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             true,
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
				format:                FormatMarkdown,
				emptyFieldPlaceholder: MarkdownDefaultEmptyFieldPlaceholder,
				wordDelimiter:         MarkdownDefaultWordDelimiter,
				hasEscape:             false,
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
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
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
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
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
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
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
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
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
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
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
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
			},
			args: args{
				v: (*string)(nil),
			},
			want:    DefaultEmptyFieldPlaceholder,
			wantErr: false,
		},
		{
			name: "non_nil_ptr_string",
			fields: fields{
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
			},
			args: args{
				v: new(string),
			},
			want:    DefaultEmptyFieldPlaceholder,
			wantErr: false,
		},
		{
			name: "non_nil_ptr_int",
			fields: fields{
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
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
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
			},
			args: args{
				v: []string{"a", "b"},
			},
			want:    "a" + DefaultWordDelimiter + "b",
			wantErr: false,
		},
		{
			name: "slice_string_included_empty",
			fields: fields{
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
			},
			args: args{
				v: []string{"a", "", "b"},
			},
			want:    "a" + DefaultWordDelimiter + DefaultEmptyFieldPlaceholder + DefaultWordDelimiter + "b",
			wantErr: false,
		},
		{
			name: "slice_int",
			fields: fields{
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
			},
			args: args{
				v: []int{0, 1, 2},
			},
			want:    "0" + DefaultWordDelimiter + "1" + DefaultWordDelimiter + "2",
			wantErr: false,
		},
		{
			name: "slice_uint",
			fields: fields{
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
			},
			args: args{
				v: []uint{0, 1, 2},
			},
			want:    "0" + DefaultWordDelimiter + "1" + DefaultWordDelimiter + "2",
			wantErr: false,
		},
		{
			name: "slice_nil",
			fields: fields{
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
			},
			args: args{
				v: ([]string)(nil),
			},
			want:    DefaultEmptyFieldPlaceholder,
			wantErr: false,
		},
		{
			name: "slice_empty",
			fields: fields{
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
			},
			args: args{
				v: []string{},
			},
			want:    DefaultEmptyFieldPlaceholder,
			wantErr: false,
		},
		{
			name: "slice_with_byte_slice",
			fields: fields{
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
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
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
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
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
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
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
			},
			args: args{
				v: &[]string{"a", "b"},
			},
			want:    "a" + DefaultWordDelimiter + "b",
			wantErr: false,
		},
		{
			name: "slice_with_ptr_to_strings",
			fields: fields{
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
			},
			args: args{
				v: []*string{sp(""), sp("a"), sp("b")},
			},
			want:    DefaultEmptyFieldPlaceholder + DefaultWordDelimiter + "a" + DefaultWordDelimiter + "b",
			wantErr: false,
		},
		{
			name: "slice_with_ptr_to_string_empty",
			fields: fields{
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
			},
			args: args{
				v: []*string{},
			},
			want:    DefaultEmptyFieldPlaceholder,
			wantErr: false,
		},
		{
			name: "slice_with_nil_ptr",
			fields: fields{
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
			},
			args: args{
				v: []*int{nil},
			},
			want:    DefaultEmptyFieldPlaceholder,
			wantErr: false,
		},
		{
			name: "slice_with_ptr_mixed",
			fields: fields{
				format:                FormatText,
				emptyFieldPlaceholder: DefaultEmptyFieldPlaceholder,
				wordDelimiter:         DefaultWordDelimiter,
				hasEscape:             false,
			},
			args: args{
				v: []*string{nil, sp(""), sp("aaa")},
			},
			want:    DefaultEmptyFieldPlaceholder + DefaultWordDelimiter + DefaultEmptyFieldPlaceholder + DefaultWordDelimiter + "aaa",
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
				hasEscape:             tt.fields.hasEscape,
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

func TestTable_replaceNL(t *testing.T) {
	type fields struct {
		format Format
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
				format: FormatText,
			},
			args: args{
				s: "aaa\nbbb\nccc",
			},
			want: "aaa\nbbb\nccc",
		},
		{
			name: "markdown",
			fields: fields{
				format: FormatMarkdown,
			},
			args: args{
				s: "aaa\nbbb\nccc",
			},
			want: "aaa<br>bbb<br>ccc",
		},
		{
			name: "backlog",
			fields: fields{
				format: FormatBacklog,
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
				format: tt.fields.format,
			}
			got := tr.replaceNL(tt.args.s)
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
				builder: tt.fields.builder,
			}
			got := tr.escape(tt.args.s)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, tt.want)
			}
		})
	}
}

func Test_pad(t *testing.T) {
	type args struct {
		s string
		w int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic",
			args: args{
				s: "dummy",
				w: 10,
			},
			want: "dummy     ",
		},
		{
			name: "japanese",
			args: args{
				s: "あいうえお",
				w: 20,
			},
			want: "あいうえお          ",
		},
		{
			name: "short",
			args: args{
				s: "dummy",
				w: 2,
			},
			want: "dummy",
		},
		{
			name: "int1",
			args: args{
				s: "0",
				w: 10,
			},
			want: "         0",
		},
		{
			name: "int2",
			args: args{
				s: "-1",
				w: 10,
			},
			want: "        -1",
		},
		{
			name: "int3",
			args: args{
				s: "01",
				w: 10,
			},
			want: "        01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pad(tt.args.s, tt.args.w)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, tt.want)
			}
		})
	}
}

func Test_padR(t *testing.T) {
	type args struct {
		s string
		w int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic",
			args: args{
				s: "dummy",
				w: 10,
			},
			want: "dummy     ",
		},
		{
			name: "japanese",
			args: args{
				s: "あいうえお",
				w: 20,
			},
			want: "あいうえお          ",
		},
		{
			name: "short",
			args: args{
				s: "dummy",
				w: 2,
			},
			want: "dummy",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := padR(tt.args.s, tt.args.w)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, tt.want)
			}
		})
	}
}

func Test_padL(t *testing.T) {
	type args struct {
		s string
		w int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic",
			args: args{
				s: "dummy",
				w: 10,
			},
			want: "     dummy",
		},
		{
			name: "japanese",
			args: args{
				s: "あいうえお",
				w: 20,
			},
			want: "          あいうえお",
		},
		{
			name: "short",
			args: args{
				s: "dummy",
				w: 2,
			},
			want: "dummy",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := padL(tt.args.s, tt.args.w)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, tt.want)
			}
		})
	}
}

func Test_isNum(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "int1",
			args: args{
				s: "1",
			},
			want: true,
		},
		{
			name: "int2",
			args: args{
				s: "-1",
			},
			want: true,
		},
		{
			name: "int3",
			args: args{
				s: "01",
			},
			want: true,
		},
		{
			name: "string",
			args: args{
				s: "dummy",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isNum(tt.args.s)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, tt.want)
			}
		})
	}
}
