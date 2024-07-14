package mintab

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type sample1 struct {
	InstanceID   string
	InstanceName string
	AttachedLB   []string
	AttachedTG   []string
}

var s1 = []sample1{
	{InstanceID: "i-1", InstanceName: "server-1", AttachedLB: []string{"lb-1"}, AttachedTG: []string{"tg-1"}},
	{InstanceID: "i-2", InstanceName: "server-2", AttachedLB: []string{"lb-2", "lb-3"}, AttachedTG: []string{"tg-2"}},
	{InstanceID: "i-3", InstanceName: "server-3", AttachedLB: []string{"lb-4"}, AttachedTG: []string{"tg-3", "tg-4"}},
	{InstanceID: "i-4", InstanceName: "server-4", AttachedLB: []string{}, AttachedTG: []string{}},
	{InstanceID: "i-5", InstanceName: "server-5", AttachedLB: []string{"lb-5"}, AttachedTG: []string{}},
	{InstanceID: "i-6", InstanceName: "server-6", AttachedLB: []string{}, AttachedTG: []string{"tg-5", "tg-6", "tg-7", "tg-8"}},
}

type sample2 struct {
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

var s2 = []sample2{
	{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-1", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 22, ToPort: 22, AddressType: "SecurityGroup", CidrBlock: "sg-10"},
	{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-1", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
	{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-2", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 443, ToPort: 443, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
	{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-2", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
	{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "icmp", FromPort: -1, ToPort: -1, AddressType: "SecurityGroup", CidrBlock: "sg-11"},
	{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 3389, ToPort: 3389, AddressType: "Ipv4", CidrBlock: "10.1.0.0/16"},
	{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 0, ToPort: 65535, AddressType: "PrefixList", CidrBlock: "pl-id/pl-name"},
	{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
}

type sample3 struct {
	RequiredField string
	UnneededField string
}

var s3 = []sample3{
	{RequiredField: "v1", UnneededField: "v1"},
	{RequiredField: "v2", UnneededField: "v2"},
}

type sample4 struct {
	Name  string
	Value string
}

var jsonSample = `{
  "key": [
    "value1",
    "value2",
    "value3",
  ]
}
`

var s4 = []sample4{
	{Name: "wildcard domain", Value: "*.example.com"},
	{Name: "empty field placeholder", Value: ""},
	{Name: "html tag", Value: "<span style=\"color:#d70910;\">red</span>"},
	{Name: "JSON", Value: jsonSample},
}

func ExampleTable_Load_basic() {
	table := New(os.Stdout)
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	table.Render()

	// Output:
	// +------------+--------------+------------+------------+
	// | InstanceID | InstanceName | AttachedLB | AttachedTG |
	// +------------+--------------+------------+------------+
	// | i-1        | server-1     | lb-1       | tg-1       |
	// +------------+--------------+------------+------------+
	// | i-2        | server-2     | lb-2       | tg-2       |
	// |            |              | lb-3       |            |
	// +------------+--------------+------------+------------+
	// | i-3        | server-3     | lb-4       | tg-3       |
	// |            |              |            | tg-4       |
	// +------------+--------------+------------+------------+
	// | i-4        | server-4     | -          | -          |
	// +------------+--------------+------------+------------+
	// | i-5        | server-5     | lb-5       | -          |
	// +------------+--------------+------------+------------+
	// | i-6        | server-6     | -          | tg-5       |
	// |            |              |            | tg-6       |
	// |            |              |            | tg-7       |
	// |            |              |            | tg-8       |
	// +------------+--------------+------------+------------+
}

func ExampleTable_Load_markdown() {
	table := New(os.Stdout, WithFormat(MarkdownFormat))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	table.Render()

	// Output:
	// | InstanceID | InstanceName | AttachedLB   | AttachedTG                   |
	// |------------|--------------|--------------|------------------------------|
	// | i-1        | server-1     | lb-1         | tg-1                         |
	// | i-2        | server-2     | lb-2<br>lb-3 | tg-2                         |
	// | i-3        | server-3     | lb-4         | tg-3<br>tg-4                 |
	// | i-4        | server-4     | \-           | \-                           |
	// | i-5        | server-5     | lb-5         | \-                           |
	// | i-6        | server-6     | \-           | tg-5<br>tg-6<br>tg-7<br>tg-8 |
}

func ExampleTable_Load_backlog() {
	table := New(os.Stdout, WithFormat(BacklogFormat))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	table.Render()

	// Output:
	// | InstanceID | InstanceName | AttachedLB   | AttachedTG                   |h
	// | i-1        | server-1     | lb-1         | tg-1                         |
	// | i-2        | server-2     | lb-2&br;lb-3 | tg-2                         |
	// | i-3        | server-3     | lb-4         | tg-3&br;tg-4                 |
	// | i-4        | server-4     | -            | -                            |
	// | i-5        | server-5     | lb-5         | -                            |
	// | i-6        | server-6     | -            | tg-5&br;tg-6&br;tg-7&br;tg-8 |
}

func ExampleTable_Load_disableheader() {
	table := New(os.Stdout, WithHeader(false))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	table.Render()

	// Output:
	// +------------+--------------+------------+------------+
	// | i-1        | server-1     | lb-1       | tg-1       |
	// +------------+--------------+------------+------------+
	// | i-2        | server-2     | lb-2       | tg-2       |
	// |            |              | lb-3       |            |
	// +------------+--------------+------------+------------+
	// | i-3        | server-3     | lb-4       | tg-3       |
	// |            |              |            | tg-4       |
	// +------------+--------------+------------+------------+
	// | i-4        | server-4     | -          | -          |
	// +------------+--------------+------------+------------+
	// | i-5        | server-5     | lb-5       | -          |
	// +------------+--------------+------------+------------+
	// | i-6        | server-6     | -          | tg-5       |
	// |            |              |            | tg-6       |
	// |            |              |            | tg-7       |
	// |            |              |            | tg-8       |
	// +------------+--------------+------------+------------+
}

func ExampleTable_Load_emptyfieldplaceholder() {
	table := New(os.Stdout, WithEmptyFieldPlaceholder("NULL"))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	table.Render()

	// Output:
	// +------------+--------------+------------+------------+
	// | InstanceID | InstanceName | AttachedLB | AttachedTG |
	// +------------+--------------+------------+------------+
	// | i-1        | server-1     | lb-1       | tg-1       |
	// +------------+--------------+------------+------------+
	// | i-2        | server-2     | lb-2       | tg-2       |
	// |            |              | lb-3       |            |
	// +------------+--------------+------------+------------+
	// | i-3        | server-3     | lb-4       | tg-3       |
	// |            |              |            | tg-4       |
	// +------------+--------------+------------+------------+
	// | i-4        | server-4     | NULL       | NULL       |
	// +------------+--------------+------------+------------+
	// | i-5        | server-5     | lb-5       | NULL       |
	// +------------+--------------+------------+------------+
	// | i-6        | server-6     | NULL       | tg-5       |
	// |            |              |            | tg-6       |
	// |            |              |            | tg-7       |
	// |            |              |            | tg-8       |
	// +------------+--------------+------------+------------+
}

func ExampleTable_Load_worddelimiter() {
	table := New(os.Stdout, WithWordDelimiter(","))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	table.Render()

	// Output:
	// +------------+--------------+------------+---------------------+
	// | InstanceID | InstanceName | AttachedLB | AttachedTG          |
	// +------------+--------------+------------+---------------------+
	// | i-1        | server-1     | lb-1       | tg-1                |
	// +------------+--------------+------------+---------------------+
	// | i-2        | server-2     | lb-2,lb-3  | tg-2                |
	// +------------+--------------+------------+---------------------+
	// | i-3        | server-3     | lb-4       | tg-3,tg-4           |
	// +------------+--------------+------------+---------------------+
	// | i-4        | server-4     | -          | -                   |
	// +------------+--------------+------------+---------------------+
	// | i-5        | server-5     | lb-5       | -                   |
	// +------------+--------------+------------+---------------------+
	// | i-6        | server-6     | -          | tg-5,tg-6,tg-7,tg-8 |
	// +------------+--------------+------------+---------------------+
}

func ExampleTable_Load_margin() {
	table := New(os.Stdout, WithMargin(3))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	table.Render()

	// Output:
	// +----------------+------------------+----------------+----------------+
	// |   InstanceID   |   InstanceName   |   AttachedLB   |   AttachedTG   |
	// +----------------+------------------+----------------+----------------+
	// |   i-1          |   server-1       |   lb-1         |   tg-1         |
	// +----------------+------------------+----------------+----------------+
	// |   i-2          |   server-2       |   lb-2         |   tg-2         |
	// |                |                  |   lb-3         |                |
	// +----------------+------------------+----------------+----------------+
	// |   i-3          |   server-3       |   lb-4         |   tg-3         |
	// |                |                  |                |   tg-4         |
	// +----------------+------------------+----------------+----------------+
	// |   i-4          |   server-4       |   -            |   -            |
	// +----------------+------------------+----------------+----------------+
	// |   i-5          |   server-5       |   lb-5         |   -            |
	// +----------------+------------------+----------------+----------------+
	// |   i-6          |   server-6       |   -            |   tg-5         |
	// |                |                  |                |   tg-6         |
	// |                |                  |                |   tg-7         |
	// |                |                  |                |   tg-8         |
	// +----------------+------------------+----------------+----------------+
}

func ExampleTable_Load_mergefields1() {
	table := New(os.Stdout)
	if err := table.Load(s2); err != nil {
		log.Fatal(err)
	}
	table.Render()

	// Output:
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | InstanceID | InstanceName | VPCID | SecurityGroupID | FlowDirection | IPProtocol | FromPort | ToPort | AddressType   | CidrBlock     |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-1        | server-1     | vpc-1 | sg-1            | Ingress       | tcp        |       22 |     22 | SecurityGroup | sg-10         |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-1        | server-1     | vpc-1 | sg-1            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-1        | server-1     | vpc-1 | sg-2            | Ingress       | tcp        |      443 |    443 | Ipv4          | 0.0.0.0/0     |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-1        | server-1     | vpc-1 | sg-2            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-2        | server-2     | vpc-1 | sg-3            | Ingress       | icmp       |       -1 |     -1 | SecurityGroup | sg-11         |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-2        | server-2     | vpc-1 | sg-3            | Ingress       | tcp        |     3389 |   3389 | Ipv4          | 10.1.0.0/16   |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-2        | server-2     | vpc-1 | sg-3            | Ingress       | tcp        |        0 |  65535 | PrefixList    | pl-id/pl-name |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-2        | server-2     | vpc-1 | sg-3            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
}

func ExampleTable_Load_mergefields2() {
	table := New(os.Stdout, WithMergeFields([]int{0, 1, 2, 3}))
	if err := table.Load(s2); err != nil {
		log.Fatal(err)
	}
	table.Render()

	// Output:
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | InstanceID | InstanceName | VPCID | SecurityGroupID | FlowDirection | IPProtocol | FromPort | ToPort | AddressType   | CidrBlock     |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-1        | server-1     | vpc-1 | sg-1            | Ingress       | tcp        |       22 |     22 | SecurityGroup | sg-10         |
	// +            +              +       +                 +---------------+------------+----------+--------+---------------+---------------+
	// |            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
	// +            +              +       +-----------------+---------------+------------+----------+--------+---------------+---------------+
	// |            |              |       | sg-2            | Ingress       | tcp        |      443 |    443 | Ipv4          | 0.0.0.0/0     |
	// +            +              +       +                 +---------------+------------+----------+--------+---------------+---------------+
	// |            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-2        | server-2     | vpc-1 | sg-3            | Ingress       | icmp       |       -1 |     -1 | SecurityGroup | sg-11         |
	// +            +              +       +                 +---------------+------------+----------+--------+---------------+---------------+
	// |            |              |       |                 | Ingress       | tcp        |     3389 |   3389 | Ipv4          | 10.1.0.0/16   |
	// +            +              +       +                 +---------------+------------+----------+--------+---------------+---------------+
	// |            |              |       |                 | Ingress       | tcp        |        0 |  65535 | PrefixList    | pl-id/pl-name |
	// +            +              +       +                 +---------------+------------+----------+--------+---------------+---------------+
	// |            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
}

func ExampleTable_Load_mergefields3() {
	table := New(os.Stdout, WithFormat(CompressedTextFormat), WithMergeFields([]int{0, 1, 2, 3}))
	if err := table.Load(s2); err != nil {
		log.Fatal(err)
	}
	table.Render()

	// Output:
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | InstanceID | InstanceName | VPCID | SecurityGroupID | FlowDirection | IPProtocol | FromPort | ToPort | AddressType   | CidrBlock     |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-1        | server-1     | vpc-1 | sg-1            | Ingress       | tcp        |       22 |     22 | SecurityGroup | sg-10         |
	// |            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
	// |            |              |       | sg-2            | Ingress       | tcp        |      443 |    443 | Ipv4          | 0.0.0.0/0     |
	// |            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-2        | server-2     | vpc-1 | sg-3            | Ingress       | icmp       |       -1 |     -1 | SecurityGroup | sg-11         |
	// |            |              |       |                 | Ingress       | tcp        |     3389 |   3389 | Ipv4          | 10.1.0.0/16   |
	// |            |              |       |                 | Ingress       | tcp        |        0 |  65535 | PrefixList    | pl-id/pl-name |
	// |            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
}

func ExampleTable_Load_ignorefields1() {
	table := New(os.Stdout, WithIgnoreFields(nil))
	if err := table.Load(s3); err != nil {
		log.Fatal(err)
	}
	table.Render()

	// Output:
	// +---------------+---------------+
	// | RequiredField | UnneededField |
	// +---------------+---------------+
	// | v1            | v1            |
	// +---------------+---------------+
	// | v2            | v2            |
	// +---------------+---------------+
}

func ExampleTable_Load_ignorefields2() {
	table := New(os.Stdout, WithIgnoreFields([]int{1}))
	if err := table.Load(s3); err != nil {
		log.Fatal(err)
	}
	table.Render()

	// Output:
	// +---------------+
	// | RequiredField |
	// +---------------+
	// | v1            |
	// +---------------+
	// | v2            |
	// +---------------+
}

func ExampleTable_Load_escape1() {
	table := New(os.Stdout, WithEscape(false))
	if err := table.Load(s4); err != nil {
		log.Fatal(err)
	}
	table.Render()

	// Output:
	// +-------------------------+-----------------------------------------+
	// | Name                    | Value                                   |
	// +-------------------------+-----------------------------------------+
	// | wildcard domain         | *.example.com                           |
	// +-------------------------+-----------------------------------------+
	// | empty field placeholder | -                                       |
	// +-------------------------+-----------------------------------------+
	// | html tag                | <span style="color:#d70910;">red</span> |
	// +-------------------------+-----------------------------------------+
	// | JSON                    | {                                       |
	// |                         |   "key": [                              |
	// |                         |     "value1",                           |
	// |                         |     "value2",                           |
	// |                         |     "value3",                           |
	// |                         |   ]                                     |
	// |                         | }                                       |
	// +-------------------------+-----------------------------------------+
}

func ExampleTable_Load_escape2() {
	table := New(os.Stdout, WithFormat(MarkdownFormat), WithEscape(false))
	if err := table.Load(s4); err != nil {
		log.Fatal(err)
	}
	table.Render()

	// Output:
	// | Name                    | Value                                                                              |
	// |-------------------------|------------------------------------------------------------------------------------|
	// | wildcard domain         | \*.example.com                                                                     |
	// | empty field placeholder | \-                                                                                 |
	// | html tag                | <span style="color:#d70910;">red</span>                                            |
	// | JSON                    | {<br>  "key": [<br>    "value1",<br>    "value2",<br>    "value3",<br>  ]<br>}<br> |
}

func ExampleTable_Load_escape3() {
	table := New(os.Stdout, WithFormat(MarkdownFormat), WithEscape(true))
	if err := table.Load(s4); err != nil {
		log.Fatal(err)
	}
	table.Render()

	// Output:
	// | Name                              | Value                                                                                                                                                                                                           |
	// |-----------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
	// | wildcard&nbsp;domain              | &#42;.example.com                                                                                                                                                                                               |
	// | empty&nbsp;field&nbsp;placeholder | \-                                                                                                                                                                                                              |
	// | html&nbsp;tag                     | &lt;span&nbsp;style=&quot;color:#d70910;&quot;&gt;red&lt;/span&gt;                                                                                                                                              |
	// | JSON                              | {<br>&nbsp;&nbsp;&quot;key&quot;:&nbsp;[<br>&nbsp;&nbsp;&nbsp;&nbsp;&quot;value1&quot;,<br>&nbsp;&nbsp;&nbsp;&nbsp;&quot;value2&quot;,<br>&nbsp;&nbsp;&nbsp;&nbsp;&quot;value3&quot;,<br>&nbsp;&nbsp;]<br>}<br> |
}

func ExampleTable_Load_string() {
	builder := &strings.Builder{}
	table := New(builder)
	if err := table.Load(s2); err != nil {
		log.Fatal(err)
	}
	table.Render()
	fmt.Println(builder.String())

	// Output:
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | InstanceID | InstanceName | VPCID | SecurityGroupID | FlowDirection | IPProtocol | FromPort | ToPort | AddressType   | CidrBlock     |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-1        | server-1     | vpc-1 | sg-1            | Ingress       | tcp        |       22 |     22 | SecurityGroup | sg-10         |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-1        | server-1     | vpc-1 | sg-1            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-1        | server-1     | vpc-1 | sg-2            | Ingress       | tcp        |      443 |    443 | Ipv4          | 0.0.0.0/0     |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-1        | server-1     | vpc-1 | sg-2            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-2        | server-2     | vpc-1 | sg-3            | Ingress       | icmp       |       -1 |     -1 | SecurityGroup | sg-11         |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-2        | server-2     | vpc-1 | sg-3            | Ingress       | tcp        |     3389 |   3389 | Ipv4          | 10.1.0.0/16   |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-2        | server-2     | vpc-1 | sg-3            | Ingress       | tcp        |        0 |  65535 | PrefixList    | pl-id/pl-name |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	// | i-2        | server-2     | vpc-1 | sg-3            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
	// +------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
}
