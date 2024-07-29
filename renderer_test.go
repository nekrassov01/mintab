package mintab

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTable_Render(t *testing.T) {
	type fields struct {
		format             Format
		header             []string
		data               [][][]string
		colWidths          []int
		lineHeights        []int
		numRows            int
		numColumns         int
		numColumnsFirstRow int
		mergeFields        []int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "text",
			fields: fields{
				format: TextFormat,
				header: []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				data: [][][]string{
					{{"i-1"}, {"server-1"}, {"lb-1"}, {"tg-1"}},
					{{"i-2"}, {"server-2"}, {"lb-2", "lb-3"}, {"tg-2"}},
					{{"i-3"}, {"server-3"}, {"lb-4"}, {"tg-3", "tg-4"}},
					{{"i-4"}, {"server-4"}, {"-"}, {"-"}},
					{{"i-5"}, {"server-5"}, {"lb-5"}, {"-"}},
					{{"i-6"}, {"server-6"}, {"-"}, {"tg-5", "tg-6", "tg-7", "tg-8"}},
				},
				colWidths:          []int{10, 12, 10, 10},
				lineHeights:        []int{1, 2, 2, 1, 1, 4},
				numRows:            6,
				numColumns:         4,
				numColumnsFirstRow: 4,
				mergeFields:        []int{},
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
				format: CompressedTextFormat,
				header: []string{"InstanceID", "InstanceName", "VPCID", "SecurityGroupID", "FlowDirection", "IPProtocol", "FromPort", "ToPort", "AddressType", "CidrBlock"},
				data: [][][]string{
					{{"i-1"}, {"server-1"}, {"vpc-1"}, {"sg-1"}, {"Ingress"}, {"tcp"}, {"22"}, {"22"}, {"SecurityGroup"}, {"sg-10"}},
					{{""}, {""}, {""}, {""}, {"Egress"}, {"-1"}, {"0"}, {"0"}, {"Ipv4"}, {"0.0.0.0/0"}},
					{{""}, {""}, {""}, {"sg-2"}, {"Ingress"}, {"tcp"}, {"443"}, {"443"}, {"Ipv4"}, {"0.0.0.0/0"}},
					{{""}, {""}, {""}, {""}, {"Egress"}, {"-1"}, {"0"}, {"0"}, {"Ipv4"}, {"0.0.0.0/0"}},
					{{"i-2"}, {"server-2"}, {"vpc-1"}, {"sg-3"}, {"Ingress"}, {"icmp"}, {"-1"}, {"-1"}, {"SecurityGroup"}, {"sg-11"}},
					{{""}, {""}, {""}, {""}, {"Ingress"}, {"tcp"}, {"3389"}, {"3389"}, {"Ipv4"}, {"10.1.0.0/16"}},
					{{""}, {""}, {""}, {""}, {"Ingress"}, {"tcp"}, {"0"}, {"65535"}, {"PrefixList"}, {"pl-id/pl-name"}},
					{{""}, {""}, {""}, {""}, {"Egress"}, {"-1"}, {"0"}, {"0"}, {"Ipv4"}, {"0.0.0.0/0"}},
				},
				colWidths:          []int{10, 12, 5, 15, 13, 10, 8, 6, 13, 13},
				lineHeights:        []int{1, 1, 1, 1, 1, 1, 1, 1},
				numRows:            8,
				numColumns:         10,
				numColumnsFirstRow: 10,
				mergeFields:        []int{0, 1, 2, 3},
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
				format: MarkdownFormat,
				header: []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				data: [][][]string{
					{{"i-1"}, {"server-1"}, {"lb-1"}, {"tg-1"}},
					{{"i-2"}, {"server-2"}, {"lb-2<br>lb-3"}, {"tg-2"}},
					{{"i-3"}, {"server-3"}, {"lb-4"}, {"tg-3<br>tg-4"}},
					{{"i-4"}, {"server-4"}, {"\\-"}, {"\\-"}},
					{{"i-5"}, {"server-5"}, {"lb-5"}, {"\\-"}},
					{{"i-6"}, {"server-6"}, {"\\-"}, {"tg-5<br>tg-6<br>tg-7<br>tg-8"}},
				},
				colWidths:          []int{10, 12, 12, 28},
				lineHeights:        []int{1, 1, 1, 1, 1, 1},
				numRows:            6,
				numColumns:         4,
				numColumnsFirstRow: 4,
				mergeFields:        []int{},
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
				format: BacklogFormat,
				header: []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				data: [][][]string{
					{{"i-1"}, {"server-1"}, {"lb-1"}, {"tg-1"}},
					{{"i-2"}, {"server-2"}, {"lb-2&br;lb-3"}, {"tg-2"}},
					{{"i-3"}, {"server-3"}, {"lb-4"}, {"tg-3&br;tg-4"}},
					{{"i-4"}, {"server-4"}, {"-"}, {"-"}},
					{{"i-5"}, {"server-5"}, {"lb-5"}, {"-"}},
					{{"i-6"}, {"server-6"}, {"-"}, {"tg-5&br;tg-6&br;tg-7&br;tg-8"}},
				},
				colWidths:          []int{10, 12, 12, 28},
				lineHeights:        []int{1, 1, 1, 1, 1, 1},
				numRows:            6,
				numColumns:         4,
				numColumnsFirstRow: 4,
				mergeFields:        []int{},
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
		{
			name: "nil",
			fields: fields{
				format:             TextFormat,
				header:             []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
				data:               [][][]string{},
				colWidths:          []int{10, 12, 10, 10},
				lineHeights:        []int{},
				numRows:            0,
				numColumns:         4,
				numColumnsFirstRow: 4,
				mergeFields:        []int{},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			tr := New(buf)
			tr.format = tt.fields.format
			tr.header = tt.fields.header
			tr.data = tt.fields.data
			tr.colWidths = tt.fields.colWidths
			tr.lineHeights = tt.fields.lineHeights
			tr.numRows = tt.fields.numRows
			tr.numColumns = tt.fields.numColumns
			tr.numColumnsFirstRow = tt.fields.numColumnsFirstRow
			tr.mergedFields = tt.fields.mergeFields
			tr.setBorder()
			tr.Render()
			if !reflect.DeepEqual(buf.String(), tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", buf.String(), tt.want)
			}
			if diff := cmp.Diff(buf.String(), tt.want); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestTable_printHeader(t *testing.T) {
	type fields struct {
		header      []string
		format      Format
		colWidths   []int
		numColumns  int
		marginWidth int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "text",
			fields: fields{
				header:      []string{"a", "bb", "ccc"},
				format:      TextFormat,
				colWidths:   []int{1, 2, 3},
				numColumns:  3,
				marginWidth: 1,
			},
			want: "| a | bb | ccc |\n",
		},
		{
			name: "markdown",
			fields: fields{
				header:      []string{"a", "bb", "ccc"},
				format:      MarkdownFormat,
				colWidths:   []int{1, 2, 3},
				numColumns:  3,
				marginWidth: 1,
			},
			want: "| a | bb | ccc |\n",
		},
		{
			name: "backlog",
			fields: fields{
				header:      []string{"a", "bb", "ccc"},
				format:      BacklogFormat,
				colWidths:   []int{1, 2, 3},
				numColumns:  3,
				marginWidth: 1,
			},
			want: "| a | bb | ccc |h\n",
		},
		{
			name: "margin",
			fields: fields{
				header:      []string{"a", "bb", "ccc"},
				format:      TextFormat,
				colWidths:   []int{1, 2, 3},
				numColumns:  3,
				marginWidth: 3,
			},
			want: "|   a   |   bb   |   ccc   |\n",
		},
		{
			name: "long",
			fields: fields{
				header:      []string{"a", "bb", "ccc"},
				format:      TextFormat,
				colWidths:   []int{10, 2, 3},
				numColumns:  3,
				marginWidth: 1,
			},
			want: "| a          | bb | ccc |\n",
		},
		{
			name: "short",
			fields: fields{
				header:      []string{"a", "bb", "ccc"},
				format:      TextFormat,
				colWidths:   []int{1, 2, 1},
				numColumns:  3,
				marginWidth: 1,
			},
			want: "| a | bb | ccc |\n",
		},
		{
			name: "nil",
			fields: fields{
				header:      []string{},
				format:      TextFormat,
				colWidths:   []int{},
				numColumns:  0,
				marginWidth: 1,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			tr := New(buf, WithMargin(tt.fields.marginWidth))
			tr.format = tt.fields.format
			tr.header = tt.fields.header
			tr.colWidths = tt.fields.colWidths
			tr.numColumns = tt.fields.numColumns
			tr.marginWidth = tt.fields.marginWidth
			tr.printHeader()
			if !reflect.DeepEqual(buf.String(), tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", buf.String(), tt.want)
			}
		})
	}
}

func TestTable_printData(t *testing.T) {
	type fields struct {
		data         [][][]string
		format       Format
		colWidths    []int
		lineHeights  []int
		numColumns   int
		hasHeader    bool
		mergedFields []int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "text",
			fields: fields{
				data: [][][]string{
					{{"i-1"}, {"server-1"}, {"lb-1"}, {"tg-1"}},
					{{"i-2"}, {"server-2"}, {"lb-2", "lb-3"}, {"tg-2"}},
					{{"i-3"}, {"server-3"}, {"lb-4"}, {"tg-3", "tg-4"}},
					{{"i-4"}, {"server-4"}, {"-"}, {"-"}},
					{{"i-5"}, {"server-5"}, {"lb-5"}, {"-"}},
					{{"i-6"}, {"server-6"}, {"-"}, {"tg-5", "tg-6", "tg-7", "tg-8"}},
				},
				format:       TextFormat,
				colWidths:    []int{10, 12, 10, 10},
				lineHeights:  []int{1, 2, 2, 1, 1, 4},
				numColumns:   4,
				hasHeader:    false,
				mergedFields: []int{},
			},
			want: `+------------+--------------+------------+------------+
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
			name: "text_with_compress",
			fields: fields{
				data: [][][]string{
					{{"i-1"}, {"server-1"}, {"vpc-1"}, {"sg-1"}, {"Ingress"}, {"tcp"}, {"22"}, {"22"}, {"SecurityGroup"}, {"sg-10"}},
					{{""}, {""}, {""}, {""}, {"Egress"}, {"-1"}, {"0"}, {"0"}, {"Ipv4"}, {"0.0.0.0/0"}},
					{{""}, {""}, {""}, {"sg-2"}, {"Ingress"}, {"tcp"}, {"443"}, {"443"}, {"Ipv4"}, {"0.0.0.0/0"}},
					{{""}, {""}, {""}, {""}, {"Egress"}, {"-1"}, {"0"}, {"0"}, {"Ipv4"}, {"0.0.0.0/0"}},
					{{"i-2"}, {"server-2"}, {"vpc-1"}, {"sg-3"}, {"Ingress"}, {"icmp"}, {"-1"}, {"-1"}, {"SecurityGroup"}, {"sg-11"}},
					{{""}, {""}, {""}, {""}, {"Ingress"}, {"tcp"}, {"3389"}, {"3389"}, {"Ipv4"}, {"10.1.0.0/16"}},
					{{""}, {""}, {""}, {""}, {"Ingress"}, {"tcp"}, {"0"}, {"65535"}, {"PrefixList"}, {"pl-id/pl-name"}},
					{{""}, {""}, {""}, {""}, {"Egress"}, {"-1"}, {"0"}, {"0"}, {"Ipv4"}, {"0.0.0.0/0"}},
				},
				format:       CompressedTextFormat,
				colWidths:    []int{10, 12, 5, 15, 13, 10, 8, 6, 13, 13},
				lineHeights:  []int{1, 1, 1, 1, 1, 1, 1, 1},
				numColumns:   10,
				hasHeader:    false,
				mergedFields: []int{0, 1, 2, 3},
			},
			want: `+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
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
			name: "markdown_with_border",
			fields: fields{
				data: [][][]string{
					{{"i-1"}, {"server-1"}, {"lb-1"}, {"tg-1"}},
					{{"i-2"}, {"server-2"}, {"lb-2<br>lb-3"}, {"tg-2"}},
					{{"i-3"}, {"server-3"}, {"lb-4"}, {"tg-3<br>tg-4"}},
					{{"i-4"}, {"server-4"}, {"\\-"}, {"\\-"}},
					{{"i-5"}, {"server-5"}, {"lb-5"}, {"\\-"}},
					{{"i-6"}, {"server-6"}, {"\\-"}, {"tg-5<br>tg-6<br>tg-7<br>tg-8"}},
				},
				format:       MarkdownFormat,
				colWidths:    []int{10, 12, 12, 28},
				lineHeights:  []int{1, 1, 1, 1, 1, 1},
				numColumns:   4,
				hasHeader:    false,
				mergedFields: []int{},
			},
			want: `|------------|--------------|--------------|------------------------------|
| i-1        | server-1     | lb-1         | tg-1                         |
| i-2        | server-2     | lb-2<br>lb-3 | tg-2                         |
| i-3        | server-3     | lb-4         | tg-3<br>tg-4                 |
| i-4        | server-4     | \-           | \-                           |
| i-5        | server-5     | lb-5         | \-                           |
| i-6        | server-6     | \-           | tg-5<br>tg-6<br>tg-7<br>tg-8 |
`,
		},
		{
			name: "markdown_with_noborder",
			fields: fields{
				data: [][][]string{
					{{"i-1"}, {"server-1"}, {"lb-1"}, {"tg-1"}},
					{{"i-2"}, {"server-2"}, {"lb-2<br>lb-3"}, {"tg-2"}},
					{{"i-3"}, {"server-3"}, {"lb-4"}, {"tg-3<br>tg-4"}},
					{{"i-4"}, {"server-4"}, {"\\-"}, {"\\-"}},
					{{"i-5"}, {"server-5"}, {"lb-5"}, {"\\-"}},
					{{"i-6"}, {"server-6"}, {"\\-"}, {"tg-5<br>tg-6<br>tg-7<br>tg-8"}},
				},
				format:       MarkdownFormat,
				colWidths:    []int{10, 12, 12, 28},
				lineHeights:  []int{1, 1, 1, 1, 1, 1},
				numColumns:   0,
				hasHeader:    false,
				mergedFields: []int{},
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
				data: [][][]string{
					{{"i-1"}, {"server-1"}, {"lb-1"}, {"tg-1"}},
					{{"i-2"}, {"server-2"}, {"lb-2&br;lb-3"}, {"tg-2"}},
					{{"i-3"}, {"server-3"}, {"lb-4"}, {"tg-3&br;tg-4"}},
					{{"i-4"}, {"server-4"}, {"-"}, {"-"}},
					{{"i-5"}, {"server-5"}, {"lb-5"}, {"-"}},
					{{"i-6"}, {"server-6"}, {"-"}, {"tg-5&br;tg-6&br;tg-7&br;tg-8"}},
				},
				format:       BacklogFormat,
				colWidths:    []int{10, 12, 12, 28},
				lineHeights:  []int{1, 1, 1, 1, 1, 1},
				numColumns:   4,
				hasHeader:    false,
				mergedFields: []int{},
			},
			want: `| i-1        | server-1     | lb-1         | tg-1                         |
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
			tr.data = tt.fields.data
			tr.format = tt.fields.format
			tr.colWidths = tt.fields.colWidths
			tr.lineHeights = tt.fields.lineHeights
			tr.numColumns = tt.fields.numColumns
			tr.hasHeader = tt.fields.hasHeader
			tr.mergedFields = tt.fields.mergedFields
			tr.setBorder()
			tr.printData()
			if !reflect.DeepEqual(buf.String(), tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", buf.String(), tt.want)
			}
		})
	}
}

func TestTable_printBorder(t *testing.T) {
	type fields struct {
		format               Format
		marginWidth          int
		marginWidthBothSides int
		colWidths            []int
	}
	type want struct {
		border string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "text",
			fields: fields{
				format:               TextFormat,
				marginWidth:          1,
				marginWidthBothSides: 2,
				colWidths:            []int{8, 12, 5},
			},
			want: want{
				border: "+----------+--------------+-------+\n",
			},
		},
		{
			name: "markdown",
			fields: fields{
				format:               MarkdownFormat,
				marginWidth:          1,
				marginWidthBothSides: 2,
				colWidths:            []int{8, 12, 5},
			},
			want: want{
				border: "|----------|--------------|-------|\n",
			},
		},
		{
			name: "backlog",
			fields: fields{
				format:               BacklogFormat,
				marginWidth:          1,
				marginWidthBothSides: 2,
				colWidths:            []int{8, 12, 5},
			},
			want: want{
				border: "|----------|--------------|-------|\n",
			},
		},
		{
			name: "wide-margin",
			fields: fields{
				format:               TextFormat,
				marginWidth:          3,
				marginWidthBothSides: 6,
				colWidths:            []int{8, 12, 5},
			},
			want: want{
				border: "+--------------+------------------+-----------+\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			tr := New(buf)
			tr.format = tt.fields.format
			tr.marginWidth = tt.fields.marginWidth
			tr.marginWidthBothSides = tt.fields.marginWidthBothSides
			tr.colWidths = tt.fields.colWidths
			tr.setBorder()
			tr.print(tr.border)
			if !reflect.DeepEqual(buf.String(), tt.want.border) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", buf.String(), tt.want.border)
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
			name: "float",
			args: args{
				s: "0.1",
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
		{
			name: "string",
			args: args{
				s: "0.0.1",
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
