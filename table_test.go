package mintab

import (
	"bytes"
	"net"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type basicTestStruct struct {
	InstanceID   string
	InstanceName string
	AttachedLB   []string
	AttachedTG   []string
}

type mergedTestStruct struct {
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

type escapedTestStruct struct {
	Name  string
	Value string
}

type nestedTestStructChild struct {
	ObjectID   int
	ObjectName string
}

type nestedTestStruct struct {
	BucketName string
	Objects    []nestedTestStructChild
}

type stringerTestStruct struct {
	ElapsedTime []time.Duration
	IPAddress   []net.IP
	NestedBytes [][]byte
}

type nonExportedTestStruct struct {
	f1 string
	f2 string
}

var jsonSample = `{
  "key": [
    "value1",
    "value2",
    "value3",
  ]
}`

var (
	basicTestInput                Input
	basicTestInputPtr             *Input
	mergedTestInput               Input
	escapedTestInput              Input
	noHeaderTestInput             Input
	nestedTestInput1              Input
	nestedTestInput2              Input
	invalidHeaderIndicesTestInput Input
	invalidDataIndicesTestInput   Input
	emptyTestInput                Input
	basicTestStructSlice          []basicTestStruct
	basicTestStructSliceEmpty     []basicTestStruct
	basicTestStructNonSlice       basicTestStruct
	basicTestStructNonSliceEmpty  basicTestStruct
	basicTestStructPtrSlice       []*basicTestStruct
	basicTestStructSlicePtr       *[]basicTestStruct
	mergedTestStructSlice         []mergedTestStruct
	escapedTestStructSlice        []escapedTestStruct
	nestedTestStructSlice         []nestedTestStruct
	stringerTestStructSlice       []stringerTestStruct
	nonExportedTestStructSlice    []nonExportedTestStruct
	nonTypeTestStructSlice        []interface{}
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	basicTestInput = Input{
		Header: []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
		Data: [][]any{
			{"i-1", "server-1", []string{"lb-1"}, []string{"tg-1"}},
			{"i-2", "server-2", []string{"lb-2", "lb-3"}, []string{"tg-2"}},
			{"i-3", "server-3", []string{"lb-4"}, []string{"tg-3", "tg-4"}},
			{"i-4", "server-4", []string{}, []string{}},
			{"i-5", "server-5", []string{"lb-5"}, []string{}},
			{"i-6", "server-6", []string{}, []string{"tg-5", "tg-6", "tg-7", "tg-8"}},
		},
	}

	basicTestInputPtr = &Input{
		Header: []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
		Data: [][]any{
			{"i-1", "server-1", []string{"lb-1"}, []string{"tg-1"}},
			{"i-2", "server-2", []string{"lb-2", "lb-3"}, []string{"tg-2"}},
			{"i-3", "server-3", []string{"lb-4"}, []string{"tg-3", "tg-4"}},
			{"i-4", "server-4", []string{}, []string{}},
			{"i-5", "server-5", []string{"lb-5"}, []string{}},
			{"i-6", "server-6", []string{}, []string{"tg-5", "tg-6", "tg-7", "tg-8"}},
		},
	}

	mergedTestInput = Input{
		Header: []string{"InstanceID", "InstanceName", "VPCID", "SecurityGroupID", "FlowDirection", "IPProtocol", "FromPort", "ToPort", "AddressType", "CidrBlock"},
		Data: [][]any{
			{"i-1", "server-1", "vpc-1", "sg-1", "Ingress", "tcp", 22, 22, "SecurityGroup", "sg-10"},
			{"i-1", "server-1", "vpc-1", "sg-1", "Egress", "-1", 0, 0, "Ipv4", "0.0.0.0/0"},
			{"i-1", "server-1", "vpc-1", "sg-2", "Ingress", "tcp", 443, 443, "Ipv4", "0.0.0.0/0"},
			{"i-1", "server-1", "vpc-1", "sg-2", "Egress", "-1", 0, 0, "Ipv4", "0.0.0.0/0"},
			{"i-2", "server-2", "vpc-1", "sg-3", "Ingress", "icmp", -1, -1, "SecurityGroup", "sg-11"},
			{"i-2", "server-2", "vpc-1", "sg-3", "Ingress", "tcp", 3389, 3389, "Ipv4", "10.1.0.0/16"},
			{"i-2", "server-2", "vpc-1", "sg-3", "Ingress", "tcp", 0, 65535, "PrefixList", "pl-id/pl-name"},
			{"i-2", "server-2", "vpc-1", "sg-3", "Egress", "-1", 0, 0, "Ipv4", "0.0.0.0/0"},
		},
	}

	escapedTestInput = Input{
		Header: []string{"Name", "Value"},
		Data: [][]any{
			{"wildcard domain", "*.example.com"},
			{"empty field placeholder", ""},
			{"html tag", "<span style=\"color:#d70910;\">red</span>"},
			{"JSON", jsonSample},
		},
	}

	noHeaderTestInput = Input{
		Header: []string{},
		Data: [][]any{
			{"i-1", "server-1", []string{"lb-1"}, []string{"tg-1"}},
			{"i-2", "server-2", []string{"lb-2", "lb-3"}, []string{"tg-2"}},
			{"i-3", "server-3", []string{"lb-4"}, []string{"tg-3", "tg-4"}},
			{"i-4", "server-4", []string{}, []string{}},
			{"i-5", "server-5", []string{"lb-5"}, []string{}},
			{"i-6", "server-6", []string{}, []string{"tg-5", "tg-6", "tg-7", "tg-8"}},
		},
	}

	nestedTestInput1 = Input{
		Header: []string{},
		Data: [][]any{
			{
				struct {
					Field1 string
					Field2 string
				}{
					Field1: "aaa",
					Field2: "bbb",
				},
			},
			{
				struct {
					Field1 string
					Field2 string
				}{
					Field1: "ccc",
					Field2: "ddd",
				},
			},
		},
	}

	nestedTestInput2 = Input{
		Header: []string{},
		Data: [][]any{
			{
				basicTestInput,
			},
		},
	}

	invalidHeaderIndicesTestInput = Input{
		Header: []string{"InstanceID", "InstanceName", "AttachedLB"}, // number of columns error
		Data: [][]any{
			{"i-1", "server-1", []string{"lb-1"}, []string{"tg-1"}},
			{"i-2", "server-2", []string{"lb-2", "lb-3"}, []string{"tg-2"}},
			{"i-3", "server-3", []string{"lb-4"}, []string{"tg-3", "tg-4"}},
			{"i-4", "server-4", []string{}, []string{}},
			{"i-5", "server-5", []string{"lb-5"}, []string{}},
			{"i-6", "server-6", []string{}, []string{"tg-5", "tg-6", "tg-7", "tg-8"}},
		},
	}

	invalidDataIndicesTestInput = Input{
		Header: []string{"InstanceID", "InstanceName", "AttachedLB", "AttachedTG"},
		Data: [][]any{
			{"i-1", "server-1", []string{"lb-1"}, []string{"tg-1"}},
			{"i-2", "server-2", []string{"lb-2", "lb-3"}, []string{"tg-2"}},
			{"i-3", "server-3", []string{"lb-4"}, []string{"tg-3", "tg-4"}},
			{"i-4", "server-4", []string{}, []string{}},
			{"i-5", "server-5", []string{"lb-5"}}, // number of columns error
			{"i-6", "server-6", []string{}, []string{"tg-5", "tg-6", "tg-7", "tg-8"}},
		},
	}

	emptyTestInput = Input{}

	basicTestStructSlice = []basicTestStruct{
		{InstanceID: "i-1", InstanceName: "server-1", AttachedLB: []string{"lb-1"}, AttachedTG: []string{"tg-1"}},
		{InstanceID: "i-2", InstanceName: "server-2", AttachedLB: []string{"lb-2", "lb-3"}, AttachedTG: []string{"tg-2"}},
		{InstanceID: "i-3", InstanceName: "server-3", AttachedLB: []string{"lb-4"}, AttachedTG: []string{"tg-3", "tg-4"}},
		{InstanceID: "i-4", InstanceName: "server-4", AttachedLB: []string{}, AttachedTG: []string{}},
		{InstanceID: "i-5", InstanceName: "server-5", AttachedLB: []string{"lb-5"}, AttachedTG: []string{}},
		{InstanceID: "i-6", InstanceName: "server-6", AttachedLB: []string{}, AttachedTG: []string{"tg-5", "tg-6", "tg-7", "tg-8"}},
	}

	basicTestStructSliceEmpty = []basicTestStruct{}

	basicTestStructNonSlice = basicTestStruct{
		InstanceID:   "i-1",
		InstanceName: "server-1",
		AttachedLB:   []string{"lb-1"},
		AttachedTG:   []string{"tg-1"},
	}

	basicTestStructNonSliceEmpty = basicTestStruct{}

	basicTestStructPtrSlice = make([]*basicTestStruct, 0, len(basicTestStructSlice))
	for i := range basicTestStructSlice {
		basicTestStructPtrSlice = append(basicTestStructPtrSlice, &basicTestStructSlice[i])
	}

	basicTestStructSlicePtr = &basicTestStructSlice

	mergedTestStructSlice = []mergedTestStruct{
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-1", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 22, ToPort: 22, AddressType: "SecurityGroup", CidrBlock: "sg-10"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-1", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-2", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 443, ToPort: 443, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-2", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "icmp", FromPort: -1, ToPort: -1, AddressType: "SecurityGroup", CidrBlock: "sg-11"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 3389, ToPort: 3389, AddressType: "Ipv4", CidrBlock: "10.1.0.0/16"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 0, ToPort: 65535, AddressType: "PrefixList", CidrBlock: "pl-id/pl-name"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
	}

	escapedTestStructSlice = []escapedTestStruct{
		{Name: "wildcard domain", Value: "*.example.com"},
		{Name: "empty field placeholder", Value: ""},
		{Name: "html tag", Value: "<span style=\"color:#d70910;\">red</span>"},
		{Name: "JSON", Value: jsonSample},
	}

	nestedTestStructSlice = []nestedTestStruct{
		{
			BucketName: "bucket1",
			Objects: []nestedTestStructChild{
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
			Objects: []nestedTestStructChild{
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

	stringerTestStructSlice = []stringerTestStruct{
		{
			ElapsedTime: []time.Duration{
				123 * time.Hour,
				234 * time.Minute,
				345 * time.Second,
			},
			IPAddress: []net.IP{
				net.IPv4(192, 168, 1, 1),
				net.IPv4(10, 0, 0, 1),
				net.ParseIP("2001:db8::68"),
			},
			NestedBytes: [][]byte{
				[]byte("aaa"),
				[]byte("bbb"),
				[]byte("ccc"),
			},
		},
	}

	nonExportedTestStructSlice = []nonExportedTestStruct{
		{
			f1: "f1",
			f2: "f2",
		},
	}

	nonTypeTestStructSlice = []interface{}{
		basicTestStruct{
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

func TestTable(t *testing.T) {
	type args struct {
		opts []Option
		v    any
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "input_text",
			args: args{
				opts: []Option{},
				v:    basicTestInput,
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
			wantErr: false,
		},
		{
			name: "input_text_ptr",
			args: args{
				opts: []Option{},
				v:    basicTestInputPtr,
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
			wantErr: false,
		},
		{
			name: "input_compressed",
			args: args{
				opts: []Option{WithFormat(CompressedTextFormat)},
				v:    basicTestInput,
			},
			want: `+------------+--------------+------------+------------+
| InstanceID | InstanceName | AttachedLB | AttachedTG |
+------------+--------------+------------+------------+
| i-1        | server-1     | lb-1       | tg-1       |
| i-2        | server-2     | lb-2       | tg-2       |
|            |              | lb-3       |            |
| i-3        | server-3     | lb-4       | tg-3       |
|            |              |            | tg-4       |
| i-4        | server-4     | -          | -          |
| i-5        | server-5     | lb-5       | -          |
| i-6        | server-6     | -          | tg-5       |
|            |              |            | tg-6       |
|            |              |            | tg-7       |
|            |              |            | tg-8       |
+------------+--------------+------------+------------+
`,
			wantErr: false,
		},
		{
			name: "input_markdown",
			args: args{
				opts: []Option{WithFormat(MarkdownFormat)},
				v:    basicTestInput,
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
			wantErr: false,
		},
		{
			name: "input_backlog",
			args: args{
				opts: []Option{WithFormat(BacklogFormat)},
				v:    basicTestInput,
			},
			want: `| InstanceID | InstanceName | AttachedLB   | AttachedTG                   |h
| i-1        | server-1     | lb-1         | tg-1                         |
| i-2        | server-2     | lb-2&br;lb-3 | tg-2                         |
| i-3        | server-3     | lb-4         | tg-3&br;tg-4                 |
| i-4        | server-4     | -            | -                            |
| i-5        | server-5     | lb-5         | -                            |
| i-6        | server-6     | -            | tg-5&br;tg-6&br;tg-7&br;tg-8 |
`,
			wantErr: false,
		},
		{
			name: "input_disable_header",
			args: args{
				opts: []Option{WithHeader(false)},
				v:    basicTestInput,
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
			wantErr: false,
		},
		{
			name: "input_emptyFieldPlaceholder",
			args: args{
				opts: []Option{WithEmptyFieldPlaceholder("")},
				v:    basicTestInput,
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
| i-4        | server-4     |            |            |
+------------+--------------+------------+------------+
| i-5        | server-5     | lb-5       |            |
+------------+--------------+------------+------------+
| i-6        | server-6     |            | tg-5       |
|            |              |            | tg-6       |
|            |              |            | tg-7       |
|            |              |            | tg-8       |
+------------+--------------+------------+------------+
`,
			wantErr: false,
		},
		{
			name: "input_wordDelimiter",
			args: args{
				opts: []Option{WithWordDelimiter(",")},
				v:    basicTestInput,
			},
			want: `+------------+--------------+------------+---------------------+
| InstanceID | InstanceName | AttachedLB | AttachedTG          |
+------------+--------------+------------+---------------------+
| i-1        | server-1     | lb-1       | tg-1                |
+------------+--------------+------------+---------------------+
| i-2        | server-2     | lb-2,lb-3  | tg-2                |
+------------+--------------+------------+---------------------+
| i-3        | server-3     | lb-4       | tg-3,tg-4           |
+------------+--------------+------------+---------------------+
| i-4        | server-4     | -          | -                   |
+------------+--------------+------------+---------------------+
| i-5        | server-5     | lb-5       | -                   |
+------------+--------------+------------+---------------------+
| i-6        | server-6     | -          | tg-5,tg-6,tg-7,tg-8 |
+------------+--------------+------------+---------------------+
`,
			wantErr: false,
		},
		{
			name: "input_mergeFields_off",
			args: args{
				opts: []Option{},
				v:    mergedTestInput,
			},
			want: `+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| InstanceID | InstanceName | VPCID | SecurityGroupID | FlowDirection | IPProtocol | FromPort | ToPort | AddressType   | CidrBlock     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-1        | server-1     | vpc-1 | sg-1            | Ingress       | tcp        |       22 |     22 | SecurityGroup | sg-10         |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-1        | server-1     | vpc-1 | sg-1            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-1        | server-1     | vpc-1 | sg-2            | Ingress       | tcp        |      443 |    443 | Ipv4          | 0.0.0.0/0     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-1        | server-1     | vpc-1 | sg-2            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-2        | server-2     | vpc-1 | sg-3            | Ingress       | icmp       |       -1 |     -1 | SecurityGroup | sg-11         |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-2        | server-2     | vpc-1 | sg-3            | Ingress       | tcp        |     3389 |   3389 | Ipv4          | 10.1.0.0/16   |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-2        | server-2     | vpc-1 | sg-3            | Ingress       | tcp        |        0 |  65535 | PrefixList    | pl-id/pl-name |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-2        | server-2     | vpc-1 | sg-3            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
`,
			wantErr: false,
		},
		{
			name: "input_mergeFields_on",
			args: args{
				opts: []Option{WithMergeFields([]int{0, 1, 2, 3})},
				v:    mergedTestInput,
			},
			want: `+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| InstanceID | InstanceName | VPCID | SecurityGroupID | FlowDirection | IPProtocol | FromPort | ToPort | AddressType   | CidrBlock     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-1        | server-1     | vpc-1 | sg-1            | Ingress       | tcp        |       22 |     22 | SecurityGroup | sg-10         |
+            +              +       +                 +---------------+------------+----------+--------+---------------+---------------+
|            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
+            +              +       +-----------------+---------------+------------+----------+--------+---------------+---------------+
|            |              |       | sg-2            | Ingress       | tcp        |      443 |    443 | Ipv4          | 0.0.0.0/0     |
+            +              +       +                 +---------------+------------+----------+--------+---------------+---------------+
|            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-2        | server-2     | vpc-1 | sg-3            | Ingress       | icmp       |       -1 |     -1 | SecurityGroup | sg-11         |
+            +              +       +                 +---------------+------------+----------+--------+---------------+---------------+
|            |              |       |                 | Ingress       | tcp        |     3389 |   3389 | Ipv4          | 10.1.0.0/16   |
+            +              +       +                 +---------------+------------+----------+--------+---------------+---------------+
|            |              |       |                 | Ingress       | tcp        |        0 |  65535 | PrefixList    | pl-id/pl-name |
+            +              +       +                 +---------------+------------+----------+--------+---------------+---------------+
|            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
`,
			wantErr: false,
		},
		{
			name: "input_mergeFields_compressed",
			args: args{
				opts: []Option{WithFormat(CompressedTextFormat), WithMergeFields([]int{0, 1, 2, 3})},
				v:    mergedTestInput,
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
			wantErr: false,
		},
		{
			name: "input_ignoreFields",
			args: args{
				opts: []Option{WithIgnoreFields([]int{1, 2, 3})},
				v:    basicTestInput,
			},
			want: `+------------+
| InstanceID |
+------------+
| i-1        |
+------------+
| i-2        |
+------------+
| i-3        |
+------------+
| i-4        |
+------------+
| i-5        |
+------------+
| i-6        |
+------------+
`,
			wantErr: false,
		},
		{
			name: "input_multiline",
			args: args{
				opts: []Option{},
				v:    escapedTestInput,
			},
			want: `+-------------------------+-----------------------------------------+
| Name                    | Value                                   |
+-------------------------+-----------------------------------------+
| wildcard domain         | *.example.com                           |
+-------------------------+-----------------------------------------+
| empty field placeholder | -                                       |
+-------------------------+-----------------------------------------+
| html tag                | <span style="color:#d70910;">red</span> |
+-------------------------+-----------------------------------------+
| JSON                    | {                                       |
|                         |   "key": [                              |
|                         |     "value1",                           |
|                         |     "value2",                           |
|                         |     "value3",                           |
|                         |   ]                                     |
|                         | }                                       |
+-------------------------+-----------------------------------------+
`,
			wantErr: false,
		},
		{
			name: "input_escaped",
			args: args{
				opts: []Option{WithFormat(MarkdownFormat), WithEscape(true)},
				v:    escapedTestInput,
			},
			want: `| Name                              | Value                                                                                                                                                                                                       |
|-----------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| wildcard&nbsp;domain              | &#42;.example.com                                                                                                                                                                                           |
| empty&nbsp;field&nbsp;placeholder | \-                                                                                                                                                                                                          |
| html&nbsp;tag                     | &lt;span&nbsp;style=&quot;color:#d70910;&quot;&gt;red&lt;/span&gt;                                                                                                                                          |
| JSON                              | {<br>&nbsp;&nbsp;&quot;key&quot;:&nbsp;[<br>&nbsp;&nbsp;&nbsp;&nbsp;&quot;value1&quot;,<br>&nbsp;&nbsp;&nbsp;&nbsp;&quot;value2&quot;,<br>&nbsp;&nbsp;&nbsp;&nbsp;&quot;value3&quot;,<br>&nbsp;&nbsp;]<br>} |
`,
			wantErr: false,
		},
		{
			name: "input_no_header",
			args: args{
				opts: []Option{},
				v:    noHeaderTestInput,
			},
			want: `+-----+----------+------+------+
| i-1 | server-1 | lb-1 | tg-1 |
+-----+----------+------+------+
| i-2 | server-2 | lb-2 | tg-2 |
|     |          | lb-3 |      |
+-----+----------+------+------+
| i-3 | server-3 | lb-4 | tg-3 |
|     |          |      | tg-4 |
+-----+----------+------+------+
| i-4 | server-4 | -    | -    |
+-----+----------+------+------+
| i-5 | server-5 | lb-5 | -    |
+-----+----------+------+------+
| i-6 | server-6 | -    | tg-5 |
|     |          |      | tg-6 |
|     |          |      | tg-7 |
|     |          |      | tg-8 |
+-----+----------+------+------+
`,
			wantErr: false,
		},
		{
			name: "struct_text",
			args: args{
				opts: []Option{},
				v:    basicTestStructSlice,
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
			wantErr: false,
		},
		{
			name: "struct_text_ptr1",
			args: args{
				opts: []Option{},
				v:    basicTestStructSlicePtr,
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
			wantErr: false,
		},
		{
			name: "struct_text_ptr2",
			args: args{
				opts: []Option{},
				v:    basicTestStructPtrSlice,
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
			wantErr: false,
		},
		{
			name: "struct_compressed",
			args: args{
				opts: []Option{WithFormat(CompressedTextFormat)},
				v:    basicTestStructSlice,
			},
			want: `+------------+--------------+------------+------------+
| InstanceID | InstanceName | AttachedLB | AttachedTG |
+------------+--------------+------------+------------+
| i-1        | server-1     | lb-1       | tg-1       |
| i-2        | server-2     | lb-2       | tg-2       |
|            |              | lb-3       |            |
| i-3        | server-3     | lb-4       | tg-3       |
|            |              |            | tg-4       |
| i-4        | server-4     | -          | -          |
| i-5        | server-5     | lb-5       | -          |
| i-6        | server-6     | -          | tg-5       |
|            |              |            | tg-6       |
|            |              |            | tg-7       |
|            |              |            | tg-8       |
+------------+--------------+------------+------------+
`,
			wantErr: false,
		},
		{
			name: "struct_markdown",
			args: args{
				opts: []Option{WithFormat(MarkdownFormat)},
				v:    basicTestStructSlice,
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
			wantErr: false,
		},
		{
			name: "struct_backlog",
			args: args{
				opts: []Option{WithFormat(BacklogFormat)},
				v:    basicTestStructSlice,
			},
			want: `| InstanceID | InstanceName | AttachedLB   | AttachedTG                   |h
| i-1        | server-1     | lb-1         | tg-1                         |
| i-2        | server-2     | lb-2&br;lb-3 | tg-2                         |
| i-3        | server-3     | lb-4         | tg-3&br;tg-4                 |
| i-4        | server-4     | -            | -                            |
| i-5        | server-5     | lb-5         | -                            |
| i-6        | server-6     | -            | tg-5&br;tg-6&br;tg-7&br;tg-8 |
`,
			wantErr: false,
		},
		{
			name: "struct_disable_header",
			args: args{
				opts: []Option{WithHeader(false)},
				v:    basicTestStructSlice,
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
			wantErr: false,
		},
		{
			name: "struct_emptyFieldPlaceholder",
			args: args{
				opts: []Option{WithEmptyFieldPlaceholder("")},
				v:    basicTestStructSlice,
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
| i-4        | server-4     |            |            |
+------------+--------------+------------+------------+
| i-5        | server-5     | lb-5       |            |
+------------+--------------+------------+------------+
| i-6        | server-6     |            | tg-5       |
|            |              |            | tg-6       |
|            |              |            | tg-7       |
|            |              |            | tg-8       |
+------------+--------------+------------+------------+
`,
			wantErr: false,
		},
		{
			name: "struct_wordDelimiter",
			args: args{
				opts: []Option{WithWordDelimiter(",")},
				v:    basicTestStructSlice,
			},
			want: `+------------+--------------+------------+---------------------+
| InstanceID | InstanceName | AttachedLB | AttachedTG          |
+------------+--------------+------------+---------------------+
| i-1        | server-1     | lb-1       | tg-1                |
+------------+--------------+------------+---------------------+
| i-2        | server-2     | lb-2,lb-3  | tg-2                |
+------------+--------------+------------+---------------------+
| i-3        | server-3     | lb-4       | tg-3,tg-4           |
+------------+--------------+------------+---------------------+
| i-4        | server-4     | -          | -                   |
+------------+--------------+------------+---------------------+
| i-5        | server-5     | lb-5       | -                   |
+------------+--------------+------------+---------------------+
| i-6        | server-6     | -          | tg-5,tg-6,tg-7,tg-8 |
+------------+--------------+------------+---------------------+
`,
			wantErr: false,
		},
		{
			name: "struct_mergeFields_off",
			args: args{
				opts: []Option{},
				v:    mergedTestStructSlice,
			},
			want: `+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| InstanceID | InstanceName | VPCID | SecurityGroupID | FlowDirection | IPProtocol | FromPort | ToPort | AddressType   | CidrBlock     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-1        | server-1     | vpc-1 | sg-1            | Ingress       | tcp        |       22 |     22 | SecurityGroup | sg-10         |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-1        | server-1     | vpc-1 | sg-1            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-1        | server-1     | vpc-1 | sg-2            | Ingress       | tcp        |      443 |    443 | Ipv4          | 0.0.0.0/0     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-1        | server-1     | vpc-1 | sg-2            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-2        | server-2     | vpc-1 | sg-3            | Ingress       | icmp       |       -1 |     -1 | SecurityGroup | sg-11         |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-2        | server-2     | vpc-1 | sg-3            | Ingress       | tcp        |     3389 |   3389 | Ipv4          | 10.1.0.0/16   |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-2        | server-2     | vpc-1 | sg-3            | Ingress       | tcp        |        0 |  65535 | PrefixList    | pl-id/pl-name |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-2        | server-2     | vpc-1 | sg-3            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
`,
			wantErr: false,
		},
		{
			name: "struct_mergeFields_on",
			args: args{
				opts: []Option{WithMergeFields([]int{0, 1, 2, 3})},
				v:    mergedTestStructSlice,
			},
			want: `+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| InstanceID | InstanceName | VPCID | SecurityGroupID | FlowDirection | IPProtocol | FromPort | ToPort | AddressType   | CidrBlock     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-1        | server-1     | vpc-1 | sg-1            | Ingress       | tcp        |       22 |     22 | SecurityGroup | sg-10         |
+            +              +       +                 +---------------+------------+----------+--------+---------------+---------------+
|            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
+            +              +       +-----------------+---------------+------------+----------+--------+---------------+---------------+
|            |              |       | sg-2            | Ingress       | tcp        |      443 |    443 | Ipv4          | 0.0.0.0/0     |
+            +              +       +                 +---------------+------------+----------+--------+---------------+---------------+
|            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
| i-2        | server-2     | vpc-1 | sg-3            | Ingress       | icmp       |       -1 |     -1 | SecurityGroup | sg-11         |
+            +              +       +                 +---------------+------------+----------+--------+---------------+---------------+
|            |              |       |                 | Ingress       | tcp        |     3389 |   3389 | Ipv4          | 10.1.0.0/16   |
+            +              +       +                 +---------------+------------+----------+--------+---------------+---------------+
|            |              |       |                 | Ingress       | tcp        |        0 |  65535 | PrefixList    | pl-id/pl-name |
+            +              +       +                 +---------------+------------+----------+--------+---------------+---------------+
|            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
`,
			wantErr: false,
		},
		{
			name: "struct_mergeFields_compressed",
			args: args{
				opts: []Option{WithFormat(CompressedTextFormat), WithMergeFields([]int{0, 1, 2, 3})},
				v:    mergedTestStructSlice,
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
			wantErr: false,
		},
		{
			name: "struct_ignoreFields",
			args: args{
				opts: []Option{WithIgnoreFields([]int{1, 2, 3})},
				v:    basicTestStructSlice,
			},
			want: `+------------+
| InstanceID |
+------------+
| i-1        |
+------------+
| i-2        |
+------------+
| i-3        |
+------------+
| i-4        |
+------------+
| i-5        |
+------------+
| i-6        |
+------------+
`,
			wantErr: false,
		},
		{
			name: "struct_multiline",
			args: args{
				opts: []Option{},
				v:    escapedTestStructSlice,
			},
			want: `+-------------------------+-----------------------------------------+
| Name                    | Value                                   |
+-------------------------+-----------------------------------------+
| wildcard domain         | *.example.com                           |
+-------------------------+-----------------------------------------+
| empty field placeholder | -                                       |
+-------------------------+-----------------------------------------+
| html tag                | <span style="color:#d70910;">red</span> |
+-------------------------+-----------------------------------------+
| JSON                    | {                                       |
|                         |   "key": [                              |
|                         |     "value1",                           |
|                         |     "value2",                           |
|                         |     "value3",                           |
|                         |   ]                                     |
|                         | }                                       |
+-------------------------+-----------------------------------------+
`,
			wantErr: false,
		},
		{
			name: "struct_escaped",
			args: args{
				opts: []Option{WithFormat(MarkdownFormat), WithEscape(true)},
				v:    escapedTestStructSlice,
			},
			want: `| Name                              | Value                                                                                                                                                                                                       |
|-----------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| wildcard&nbsp;domain              | &#42;.example.com                                                                                                                                                                                           |
| empty&nbsp;field&nbsp;placeholder | \-                                                                                                                                                                                                          |
| html&nbsp;tag                     | &lt;span&nbsp;style=&quot;color:#d70910;&quot;&gt;red&lt;/span&gt;                                                                                                                                          |
| JSON                              | {<br>&nbsp;&nbsp;&quot;key&quot;:&nbsp;[<br>&nbsp;&nbsp;&nbsp;&nbsp;&quot;value1&quot;,<br>&nbsp;&nbsp;&nbsp;&nbsp;&quot;value2&quot;,<br>&nbsp;&nbsp;&nbsp;&nbsp;&quot;value3&quot;,<br>&nbsp;&nbsp;]<br>} |
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			table := New(buf, tt.args.opts...)
			if err := table.Load(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", err, tt.wantErr)
				return
			}
			table.Render()
			if !reflect.DeepEqual(buf.String(), tt.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v\n", buf.String(), tt.want)
			}
			if diff := cmp.Diff(buf.String(), tt.want); diff != "" {
				t.Error(diff)
			}
		})
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
				w:                     &bytes.Buffer{},
				b:                     strings.Builder{},
				header:                nil,
				format:                TextFormat,
				newLine:               textNewLine,
				border:                "",
				marginWidth:           1,
				marginWidthBothSides:  2,
				margin:                " ",
				emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
				wordDelimiter:         TextDefaultWordDelimiter,
				mergedFields:          nil,
				ignoredFields:         nil,
				colWidths:             nil,
				hasHeader:             true,
				isEscape:              false,
			},
		},
		{
			name: "not-default",
			args: args{
				opts: []Option{
					WithFormat(MarkdownFormat),
					WithHeader(false),
					WithMargin(2),
					WithEmptyFieldPlaceholder(MarkdownDefaultEmptyFieldPlaceholder),
					WithWordDelimiter(MarkdownDefaultWordDelimiter),
					WithMergeFields([]int{0}),
					WithIgnoreFields([]int{0}),
					WithEscape(true),
				},
			},
			want: &Table{
				w:                     &bytes.Buffer{},
				b:                     strings.Builder{},
				header:                nil,
				format:                MarkdownFormat,
				newLine:               textNewLine, // change after setFormat()
				border:                "",
				marginWidth:           2,
				marginWidthBothSides:  4,
				margin:                "  ",
				emptyFieldPlaceholder: MarkdownDefaultEmptyFieldPlaceholder,
				wordDelimiter:         MarkdownDefaultWordDelimiter,
				mergedFields:          []int{0},
				ignoredFields:         []int{0},
				colWidths:             nil,
				hasHeader:             false,
				isEscape:              true,
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

func TestWithMargin(t *testing.T) {
	tests := []struct {
		name  string
		width int
		want  int
	}{
		{
			name:  "default",
			width: 1,
			want:  1,
		},
		{
			name:  "change",
			width: 2,
			want:  2,
		},
		{
			name:  "signed",
			width: -1,
			want:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{}
			opt := WithMargin(tt.width)
			opt(table)
			if table.marginWidth != tt.want {
				t.Errorf("\ngot\n%v\nset\n%v\nwant\n%v\n", tt.width, table.marginWidth, tt.want)
			}
		})
	}
}

func TestWithEmptyFieldPlacehplder(t *testing.T) {
	tests := []struct {
		name                  string
		emptyFieldPlaceholder string
		want                  string
	}{
		{
			name:                  "default",
			emptyFieldPlaceholder: TextDefaultEmptyFieldPlaceholder,
			want:                  TextDefaultEmptyFieldPlaceholder,
		},
		{
			name:                  "blank",
			emptyFieldPlaceholder: "",
			want:                  " ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{}
			opt := WithEmptyFieldPlaceholder(tt.emptyFieldPlaceholder)
			opt(table)
			if table.emptyFieldPlaceholder != tt.want {
				t.Errorf("\ngot\n%v\nset\n%v\nwant\n%v\n", tt.emptyFieldPlaceholder, table.emptyFieldPlaceholder, tt.want)
			}
		})
	}
}
