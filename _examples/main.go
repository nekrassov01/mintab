package main

import (
	"fmt"
	"log"

	"github.com/nekrassov01/mintab"
)

func main() {
	/*
		Basic
	*/

	type sample1 struct {
		InstanceID   string
		InstanceName string
		AttachedLB   []string
		AttachedTG   []string
	}

	s1 := []sample1{
		{InstanceID: "i-1", InstanceName: "server-1", AttachedLB: []string{"lb-domain-1"}, AttachedTG: []string{"tg-1"}},
		{InstanceID: "i-2", InstanceName: "server-2", AttachedLB: []string{"lb-doamin-2", "lb-doamin-3"}, AttachedTG: []string{"tg-2"}},
		{InstanceID: "i-3", InstanceName: "server-3", AttachedLB: []string{"lb-doamin-4"}, AttachedTG: []string{"tg-3", "tg-4"}},
		{InstanceID: "i-4", InstanceName: "server-4", AttachedLB: []string{}, AttachedTG: []string{}},
		{InstanceID: "i-5", InstanceName: "server-5", AttachedLB: []string{"lb-domain-5"}, AttachedTG: []string{}},
		{InstanceID: "i-5", InstanceName: "server-5", AttachedLB: []string{}, AttachedTG: []string{"tg-5", "tg-6"}},
	}

	var table *mintab.Table

	table = mintab.NewTable()
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*
		| InstanceID | InstanceName | AttachedLB  | AttachedTG |
		|------------|--------------|-------------|------------|
		| i-1        | server-1     | lb-domain-1 | tg-1       |
		| i-2        | server-2     | lb-doamin-2 | tg-2       |
		|            |              | lb-doamin-3 |            |
		| i-3        | server-3     | lb-doamin-4 | tg-3       |
		|            |              |             | tg-4       |
		| i-4        | server-4     | -           | -          |
		| i-5        | server-5     | lb-domain-5 | -          |
		| i-5        | server-5     | -           | tg-5       |
		|            |              |             | tg-6       |
	*/

	table = mintab.NewTable(mintab.WithFormat(mintab.FormatMarkdown))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*
		| InstanceID | InstanceName | AttachedLB                 | AttachedTG   |
		|------------|--------------|----------------------------|--------------|
		| i-1        | server-1     | lb-domain-1                | tg-1         |
		| i-2        | server-2     | lb-doamin-2<br>lb-doamin-3 | tg-2         |
		| i-3        | server-3     | lb-doamin-4                | tg-3<br>tg-4 |
		| i-4        | server-4     | &#45;                      | &#45;        |
		| i-5        | server-5     | lb-domain-5                | &#45;        |
		| i-5        | server-5     | &#45;                      | tg-5<br>tg-6 |
	*/

	table = mintab.NewTable(mintab.WithFormat(mintab.FormatBacklog))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*
		| InstanceID | InstanceName | AttachedLB                 | AttachedTG   |h
		| i-1        | server-1     | lb-domain-1                | tg-1         |
		| i-2        | server-2     | lb-doamin-2&br;lb-doamin-3 | tg-2         |
		| i-3        | server-3     | lb-doamin-4                | tg-3&br;tg-4 |
		| i-4        | server-4     | -                          | -            |
		| i-5        | server-5     | lb-domain-5                | -            |
		| i-5        | server-5     | -                          | tg-5&br;tg-6 |
	*/

	table = mintab.NewTable(mintab.WithHeader(false))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*
		| i-1 | server-1 | lb-domain-1 | tg-1 |
		| i-2 | server-2 | lb-doamin-2 | tg-2 |
		|     |          | lb-doamin-3 |      |
		| i-3 | server-3 | lb-doamin-4 | tg-3 |
		|     |          |             | tg-4 |
		| i-4 | server-4 | -           | -    |
		| i-5 | server-5 | lb-domain-5 | -    |
		| i-5 | server-5 | -           | tg-5 |
		|     |          |             | tg-6 |
	*/

	table = mintab.NewTable(mintab.WithEmptyFieldPlaceholder("NULL"))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*
		| InstanceID | InstanceName | AttachedLB  | AttachedTG |
		|------------|--------------|-------------|------------|
		| i-1        | server-1     | lb-domain-1 | tg-1       |
		| i-2        | server-2     | lb-doamin-2 | tg-2       |
		|            |              | lb-doamin-3 |            |
		| i-3        | server-3     | lb-doamin-4 | tg-3       |
		|            |              |             | tg-4       |
		| i-4        | server-4     | NULL        | NULL       |
		| i-5        | server-5     | lb-domain-5 | NULL       |
		| i-5        | server-5     | NULL        | tg-5       |
		|            |              |             | tg-6       |
	*/

	table = mintab.NewTable(mintab.WithWordDelimiter(","))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*
		| InstanceID | InstanceName | AttachedLB              | AttachedTG |
		|------------|--------------|-------------------------|------------|
		| i-1        | server-1     | lb-domain-1             | tg-1       |
		| i-2        | server-2     | lb-doamin-2,lb-doamin-3 | tg-2       |
		| i-3        | server-3     | lb-doamin-4             | tg-3,tg-4  |
		| i-4        | server-4     | -                       | -          |
		| i-5        | server-5     | lb-domain-5             | -          |
		| i-5        | server-5     | -                       | tg-5,tg-6  |
	*/

	table = mintab.NewTable(mintab.WithMergeFields([]int{0, 1}), mintab.WithTheme(mintab.ThemeDark))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*
		| InstanceID | InstanceName | AttachedLB  | AttachedTG |
		|------------|--------------|-------------|------------|
		| i-1        | server-1     | lb-domain-1 | tg-1       |
		| i-2        | server-2     | lb-doamin-2 | tg-2       |
		|            |              | lb-doamin-3 |            |
		| i-3        | server-3     | lb-doamin-4 | tg-3       |
		|            |              |             | tg-4       |
		| i-4        | server-4     | -           | -          |
		| i-5        | server-5     | lb-domain-5 | -          |
		|            |              | -           | tg-5       |
		|            |              |             | tg-6       |
	*/

	table = mintab.NewTable(mintab.WithIgnoreFields([]int{1}))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*
		| InstanceID | AttachedLB  | AttachedTG |
		|------------|-------------|------------|
		| i-1        | lb-domain-1 | tg-1       |
		| i-2        | lb-doamin-2 | tg-2       |
		|            | lb-doamin-3 |            |
		| i-3        | lb-doamin-4 | tg-3       |
		|            |             | tg-4       |
		| i-4        | -           | -          |
		| i-5        | lb-domain-5 | -          |
		| i-5        | -           | tg-5       |
		|            |             | tg-6       |
	*/

	/*
		Escaping when select markdown format
	*/

	type sample2 struct {
		Domain string
	}

	s2 := []sample2{
		{Domain: "*.example.com"},
		{Domain: "| _"},
	}

	table = mintab.NewTable(mintab.WithFormat(mintab.FormatMarkdown))
	if err := table.Load(s2); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*
		| Domain             |
		|--------------------|
		| &#42;.example.com  |
		| &#124;&nbsp;&#095; |
	*/

	/*
		Grouping
	*/

	type sample3 struct {
		InstanceID      string
		InstanceName    string
		SecurityGroupID string
		FlowDirection   string
		IPProtocol      string
		FromPort        int
		ToPort          int
		AddressType     string
		CidrBlock       string
	}

	s3 := []sample3{
		{InstanceID: "i-1", InstanceName: "server-1", SecurityGroupID: "sg-1", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 22, ToPort: 22, AddressType: "SecurityGroup", CidrBlock: "sg-10"},
		{InstanceID: "i-1", InstanceName: "server-1", SecurityGroupID: "sg-1", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-1", InstanceName: "server-1", SecurityGroupID: "sg-2", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 443, ToPort: 443, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-1", InstanceName: "server-1", SecurityGroupID: "sg-2", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-2", InstanceName: "server-2", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "icmp", FromPort: -1, ToPort: -1, AddressType: "SecurityGroup", CidrBlock: "sg-11"},
		{InstanceID: "i-2", InstanceName: "server-2", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 3389, ToPort: 3389, AddressType: "Ipv4", CidrBlock: "10.1.0.0/16"},
		{InstanceID: "i-2", InstanceName: "server-2", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 0, ToPort: 65535, AddressType: "PrefixList", CidrBlock: "pl-id/pl-name"},
		{InstanceID: "i-2", InstanceName: "server-2", SecurityGroupID: "sg-3", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
	}

	table = mintab.NewTable()
	if err := table.Load(s3); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*
		| InstanceID | InstanceName | SecurityGroupID | FlowDirection | IPProtocol | FromPort | ToPort | AddressType   | CidrBlock     |
		|------------|--------------|-----------------|---------------|------------|----------|--------|---------------|---------------|
		| i-1        | server-1     | sg-1            | Ingress       | tcp        |       22 |     22 | SecurityGroup | sg-10         |
		| i-1        | server-1     | sg-1            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
		| i-1        | server-1     | sg-2            | Ingress       | tcp        |      443 |    443 | Ipv4          | 0.0.0.0/0     |
		| i-1        | server-1     | sg-2            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
		| i-2        | server-2     | sg-3            | Ingress       | icmp       |       -1 |     -1 | SecurityGroup | sg-11         |
		| i-2        | server-2     | sg-3            | Ingress       | tcp        |     3389 |   3389 | Ipv4          | 10.1.0.0/16   |
		| i-2        | server-2     | sg-3            | Ingress       | tcp        |        0 |  65535 | PrefixList    | pl-id/pl-name |
		| i-2        | server-2     | sg-3            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
	*/

	table = mintab.NewTable(mintab.WithMergeFields([]int{0, 1}), mintab.WithTheme(mintab.ThemeDark))
	if err := table.Load(s3); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*
		| InstanceID | InstanceName | SecurityGroupID | FlowDirection | IPProtocol | FromPort | ToPort | AddressType   | CidrBlock     |
		|------------|--------------|-----------------|---------------|------------|----------|--------|---------------|---------------|
		| i-1        | server-1     | sg-1            | Ingress       | tcp        |       22 |     22 | SecurityGroup | sg-10         |
		|            |              | sg-1            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
		|            |              | sg-2            | Ingress       | tcp        |      443 |    443 | Ipv4          | 0.0.0.0/0     |
		|            |              | sg-2            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
		| i-2        | server-2     | sg-3            | Ingress       | icmp       |       -1 |     -1 | SecurityGroup | sg-11         |
		|            |              | sg-3            | Ingress       | tcp        |     3389 |   3389 | Ipv4          | 10.1.0.0/16   |
		|            |              | sg-3            | Ingress       | tcp        |        0 |  65535 | PrefixList    | pl-id/pl-name |
		|            |              | sg-3            | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0     |
	*/

	table = mintab.NewTable(mintab.WithIgnoreFields([]int{1}), mintab.WithFormat(mintab.FormatMarkdown))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*
		| InstanceID | AttachedLB  | AttachedTG |
		|------------|-------------|------------|
		| i-1        | lb-domain-1 | tg-1       |
		| i-2        | lb-doamin-2 | tg-2       |
		|            | lb-doamin-3 |            |
		| i-3        | lb-doamin-4 | tg-3       |
		|            |             | tg-4       |
		| i-4        | -           | -          |
		| i-5        | lb-domain-5 | -          |
		| i-5        | -           | tg-5       |
		|            |             | tg-6       |
	*/

}
