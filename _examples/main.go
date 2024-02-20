package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nekrassov01/mintab"
)

func main() {
	/*
		Basic
	*/

	fmt.Println("/* Basic */")
	fmt.Println()

	type sample1 struct {
		InstanceID   string
		InstanceName string
		AttachedLB   []string
		AttachedTG   []string
	}

	s1 := []sample1{
		{InstanceID: "i-1", InstanceName: "server-1", AttachedLB: []string{"lb-1"}, AttachedTG: []string{"tg-1"}},
		{InstanceID: "i-2", InstanceName: "server-2", AttachedLB: []string{"lb-2", "lb-3"}, AttachedTG: []string{"tg-2"}},
		{InstanceID: "i-3", InstanceName: "server-3", AttachedLB: []string{"lb-4"}, AttachedTG: []string{"tg-3", "tg-4"}},
		{InstanceID: "i-4", InstanceName: "server-4", AttachedLB: []string{}, AttachedTG: []string{}},
		{InstanceID: "i-5", InstanceName: "server-5", AttachedLB: []string{"lb-5"}, AttachedTG: []string{}},
		{InstanceID: "i-6", InstanceName: "server-6", AttachedLB: []string{}, AttachedTG: []string{"tg-5", "tg-6", "tg-7", "tg-8"}},
	}

	var table *mintab.Table

	fmt.Println("// format: text")
	fmt.Println()
	table = mintab.New(os.Stdout)
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	table.Out()

	/*
		+------------+--------------+------------+------------+
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
	*/

	fmt.Println("// format: markdown")
	fmt.Println()
	table = mintab.New(os.Stdout, mintab.WithFormat(mintab.FormatMarkdown))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	table.Out()

	/*
		| InstanceID | InstanceName | AttachedLB   | AttachedTG                   |
		|------------|--------------|--------------|------------------------------|
		| i-1        | server-1     | lb-1         | tg-1                         |
		| i-2        | server-2     | lb-2<br>lb-3 | tg-2                         |
		| i-3        | server-3     | lb-4         | tg-3<br>tg-4                 |
		| i-4        | server-4     | \-           | \-                           |
		| i-5        | server-5     | lb-5         | \-                           |
		| i-6        | server-6     | \-           | tg-5<br>tg-6<br>tg-7<br>tg-8 |
	*/

	fmt.Println("// format: backlog")
	fmt.Println()
	table = mintab.New(os.Stdout, mintab.WithFormat(mintab.FormatBacklog))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	table.Out()

	/*
		| InstanceID | InstanceName | AttachedLB   | AttachedTG                   |h
		| i-1        | server-1     | lb-1         | tg-1                         |
		| i-2        | server-2     | lb-2&br;lb-3 | tg-2                         |
		| i-3        | server-3     | lb-4         | tg-3&br;tg-4                 |
		| i-4        | server-4     | -            | -                            |
		| i-5        | server-5     | lb-5         | -                            |
		| i-6        | server-6     | -            | tg-5&br;tg-6&br;tg-7&br;tg-8 |
	*/

	fmt.Println("// format: text")
	fmt.Println("// header: false")
	fmt.Println()
	table = mintab.New(os.Stdout, mintab.WithHeader(false))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	table.Out()

	/*
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
	*/

	fmt.Println("// format: text")
	fmt.Println("// emptyFieldPlaceholder: \"NULL\"")
	fmt.Println()
	table = mintab.New(os.Stdout, mintab.WithEmptyFieldPlaceholder("NULL"))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	table.Out()

	/*
		+------------+--------------+------------+------------+
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
		| i-4        | server-4     | NULL       | NULL       |
		+------------+--------------+------------+------------+
		| i-5        | server-5     | lb-5       | NULL       |
		+------------+--------------+------------+------------+
		| i-6        | server-6     | NULL       | tg-5       |
		|            |              |            | tg-6       |
		|            |              |            | tg-7       |
		|            |              |            | tg-8       |
		+------------+--------------+------------+------------+
	*/

	fmt.Println("// format: text")
	fmt.Println("// wordDelimiter: \",\"")
	fmt.Println()
	table = mintab.New(os.Stdout, mintab.WithWordDelimiter(","))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	table.Out()

	/*
		+------------+--------------+------------+---------------------+
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
	*/

	fmt.Println("// format: text")
	fmt.Println("// margin: 3")
	fmt.Println()
	table = mintab.New(os.Stdout, mintab.WithMargin(3))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	table.Out()

	/*
		+----------------+------------------+----------------+----------------+
		|   InstanceID   |   InstanceName   |   AttachedLB   |   AttachedTG   |
		+----------------+------------------+----------------+----------------+
		|   i-1          |   server-1       |   lb-1         |   tg-1         |
		+----------------+------------------+----------------+----------------+
		|   i-2          |   server-2       |   lb-2         |   tg-2         |
		|                |                  |   lb-3         |                |
		+----------------+------------------+----------------+----------------+
		|   i-3          |   server-3       |   lb-4         |   tg-3         |
		|                |                  |                |   tg-4         |
		+----------------+------------------+----------------+----------------+
		|   i-4          |   server-4       |   -            |   -            |
		+----------------+------------------+----------------+----------------+
		|   i-5          |   server-5       |   lb-5         |   -            |
		+----------------+------------------+----------------+----------------+
		|   i-6          |   server-6       |   -            |   tg-5         |
		|                |                  |                |   tg-6         |
		|                |                  |                |   tg-7         |
		|                |                  |                |   tg-8         |
		+----------------+------------------+----------------+----------------+
	*/

	/*
		Merge fields
	*/

	fmt.Println("/* Merge fields */")
	fmt.Println()

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

	s2 := []sample2{
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-1", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 22, ToPort: 22, AddressType: "SecurityGroup", CidrBlock: "sg-10"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-1", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-2", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 443, ToPort: 443, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-2", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "icmp", FromPort: -1, ToPort: -1, AddressType: "SecurityGroup", CidrBlock: "sg-11"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 3389, ToPort: 3389, AddressType: "Ipv4", CidrBlock: "10.1.0.0/16"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 0, ToPort: 65535, AddressType: "PrefixList", CidrBlock: "pl-id/pl-name"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
	}

	fmt.Println("// format: text")
	fmt.Println("// mergeFields: nil")
	fmt.Println()
	table = mintab.New(os.Stdout)
	if err := table.Load(s2); err != nil {
		log.Fatal(err)
	}
	table.Out()

	/*
		+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
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
	*/

	fmt.Println("// format: text")
	fmt.Println("// mergeFields: []int{0, 1, 2}")
	fmt.Println()
	table = mintab.New(os.Stdout, mintab.WithMergeFields([]int{0, 1, 2}))
	if err := table.Load(s2); err != nil {
		log.Fatal(err)
	}
	table.Out()

	/*
		+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
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
		+            +              +       +                 +               +------------+----------+--------+---------------+---------------+
		|            |              |       |                 |               | tcp        |     3389 |   3389 | Ipv4          | 10.1.0.0/16   |
		+            +              +       +                 +               +            +----------+--------+---------------+---------------+
		|            |              |       |                 |               |            |        0 |  65535 | PrefixList    | pl-id/pl-name |
		+            +              +       +                 +---------------+------------+----------+--------+---------------+---------------+
		|            |              |       |                 | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
		+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
	*/

	fmt.Println("// format: text")
	fmt.Println("// mergeFields: []int{0, 1, 2, 3}")
	fmt.Println()
	table = mintab.New(os.Stdout, mintab.WithMergeFields([]int{0, 1, 2, 3}))
	if err := table.Load(s2); err != nil {
		log.Fatal(err)
	}
	table.Out()

	/*
		+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
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
	*/

	fmt.Println("// format: compressed")
	fmt.Println("// mergeFields: []int{0, 1, 2, 3}")
	fmt.Println()
	table = mintab.New(os.Stdout, mintab.WithFormat(mintab.FormatCompressedText), mintab.WithMergeFields([]int{0, 1, 2, 3}))
	if err := table.Load(s2); err != nil {
		log.Fatal(err)
	}
	table.Out()

	/*
		+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
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
	*/

	/*
		Ignore fields
	*/

	fmt.Println("/* Ignore fields */")
	fmt.Println()

	type sample3 struct {
		RequiredField string
		UnneededField string
	}

	s3 := []sample3{
		{RequiredField: "v1", UnneededField: "v1"},
		{RequiredField: "v2", UnneededField: "v2"},
	}

	fmt.Println("// format: text")
	fmt.Println("// ignoredFields: nil")
	fmt.Println()
	table = mintab.New(os.Stdout, mintab.WithIgnoreFields(nil))
	if err := table.Load(s3); err != nil {
		log.Fatal(err)
	}
	table.Out()

	/*
		+---------------+---------------+
		| RequiredField | UnneededField |
		+---------------+---------------+
		| v1            | v1            |
		+---------------+---------------+
		| v2            | v2            |
		+---------------+---------------+
	*/

	fmt.Println("// format: text")
	fmt.Println("// ignoredFields: []int{1}")
	fmt.Println()
	table = mintab.New(os.Stdout, mintab.WithIgnoreFields([]int{1}))
	if err := table.Load(s3); err != nil {
		log.Fatal(err)
	}
	table.Out()

	/*
		+---------------+
		| RequiredField |
		+---------------+
		| v1            |
		+---------------+
		| v2            |
		+---------------+
	*/

	/*
		Escaping html special chars
	*/

	fmt.Println("/* Escaping html special chars */")
	fmt.Println()

	type sample4 struct {
		Name           string
		EscatableValue string
	}

	s4 := []sample4{
		{Name: "wildcard domain", EscatableValue: "*.example.com"},
		{Name: "empty field placeholder", EscatableValue: ""},
		{Name: "html tag", EscatableValue: "<span style=\"color:#d70910;\">red</span>"},
		{Name: "JSON", EscatableValue: `{
  "key": [
    "value1",
    "value2",
    "value3",
  ]
}
		`},
	}

	fmt.Println("// format: text")
	fmt.Println("// escape: false")
	fmt.Println()
	table = mintab.New(os.Stdout, mintab.WithEscape(false))
	if err := table.Load(s4); err != nil {
		log.Fatal(err)
	}
	table.Out()

	/*
		+-------------------------+-----------------------------------------+
		| Name                    | EscatableValue                          |
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
	*/

	fmt.Println("// format: markdown")
	fmt.Println("// escape: false")
	fmt.Println()
	table = mintab.New(os.Stdout, mintab.WithFormat(mintab.FormatMarkdown), mintab.WithEscape(false))
	if err := table.Load(s4); err != nil {
		log.Fatal(err)
	}
	table.Out()

	/*
		| Name                    | EscatableValue                                                                     |
		|-------------------------|------------------------------------------------------------------------------------|
		| wildcard domain         | \*.example.com                                                                     |
		| empty field placeholder | \-                                                                                 |
		| html tag                | <span style="color:#d70910;">red</span>                                            |
		| JSON                    | {<br>  "key": [<br>    "value1",<br>    "value2",<br>    "value3",<br>  ]<br>}<br> |
	*/

	fmt.Println("// format: markdown")
	fmt.Println("// escape: true")
	fmt.Println()
	table = mintab.New(os.Stdout, mintab.WithFormat(mintab.FormatMarkdown), mintab.WithEscape(true))
	if err := table.Load(s4); err != nil {
		log.Fatal(err)
	}
	table.Out()

	/*
		| Name                              | EscatableValue                                                                                                                                                                                                  |
		|-----------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
		| wildcard&nbsp;domain              | &#42;.example.com                                                                                                                                                                                               |
		| empty&nbsp;field&nbsp;placeholder | \-                                                                                                                                                                                                              |
		| html&nbsp;tag                     | &lt;span&nbsp;style=&quot;color:#d70910;&quot;&gt;red&lt;/span&gt;                                                                                                                                              |
		| JSON                              | {<br>&nbsp;&nbsp;&quot;key&quot;:&nbsp;[<br>&nbsp;&nbsp;&nbsp;&nbsp;&quot;value1&quot;,<br>&nbsp;&nbsp;&nbsp;&nbsp;&quot;value2&quot;,<br>&nbsp;&nbsp;&nbsp;&nbsp;&quot;value3&quot;,<br>&nbsp;&nbsp;]<br>}<br> |
	*/

	/*
		Using strings.Builder
	*/

	fmt.Println("/* Using strings.Builder */")
	fmt.Println()

	fmt.Println("// format: text")
	fmt.Println()
	var builder strings.Builder
	table = mintab.New(&builder)
	if err := table.Load(s2); err != nil {
		log.Fatal(err)
	}
	table.Out()
	fmt.Println(builder.String())

	/*
		+------------+--------------+-------+-----------------+---------------+------------+----------+--------+---------------+---------------+
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
	*/

	table = mintab.New(os.Stdout, mintab.WithFormat(mintab.FormatMarkdown), mintab.WithMergeFields([]int{0, 1, 2, 3}))
	if err := table.Load(s2); err != nil {
		log.Fatal(err)
	}
	table.Out()

	table = mintab.New(os.Stdout, mintab.WithFormat(mintab.FormatBacklog), mintab.WithMergeFields([]int{0, 1, 2, 3}))
	if err := table.Load(s2); err != nil {
		log.Fatal(err)
	}
	table.Out()
}
