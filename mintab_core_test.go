package mintab

import (
	"bytes"
	"net"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
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

type stringerSample struct {
	ElapsedTime time.Duration
	IPAddress   net.IP
}

type nestedStringerSample struct {
	ElapsedTime []time.Duration
	IPAddress   []net.IP
	NestedBytes [][]byte
}

type nonExportedSample struct {
	f1 string
	f2 string
}

var (
	basicsample              []basicSample
	basicsampleEmpty         []basicSample
	basicsampleNonSlice      basicSample
	basicsampleNonSliceEmpty basicSample
	nestedsample             []nestedSample
	mergedsample             []mergedSample
	stringersample           stringerSample
	nestedstringerSample     nestedStringerSample
	basicsamplePtr           []*basicSample
	basicsampleSlicePtr      *[]basicSample
	irregularsample          []interface{}
	nonexportedsample        []nonExportedSample
	basicinputsample         Input
	basicinputsamplePtr      *Input
	noheadersample           Input
	irregularinputsample1    Input
	irregularinputsample2    Input
	irregularinputsample3    Input
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
	basicsampleEmpty = []basicSample{}
	basicsampleNonSlice = basicSample{
		InstanceID:   "i-1",
		InstanceName: "server-1",
		AttachedLB:   []string{"lb-1"},
		AttachedTG:   []string{"tg-1"},
	}
	basicsampleNonSliceEmpty = basicSample{}
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
	stringersample = stringerSample{
		ElapsedTime: 123 * time.Hour,
		IPAddress:   net.IPv4allsys,
	}
	nestedstringerSample = nestedStringerSample{
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
	}
	basicsamplePtr = make([]*basicSample, 0, len(basicsample))
	for i := range basicsample {
		basicsamplePtr = append(basicsamplePtr, &basicsample[i])
	}
	basicsampleSlicePtr = &basicsample
	nonexportedsample = []nonExportedSample{
		{
			f1: "f1",
			f2: "f2",
		},
	}
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
	basicinputsample = Input{
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
	basicinputsamplePtr = &Input{
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
	noheadersample = Input{
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
	irregularinputsample1 = Input{
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
	irregularinputsample2 = Input{
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
	irregularinputsample3 = Input{
		Header: []string{"NestedField"},
		Data: [][]any{
			{basicsample},
		},
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
				data:                  nil,
				header:                nil,
				format:                TextFormat,
				newLine:               textNewLine,
				border:                "",
				marginWidth:           1,
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
				data:                  nil,
				header:                nil,
				format:                MarkdownFormat,
				newLine:               textNewLine, // change after setFormat()
				border:                "",
				marginWidth:           2,
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
