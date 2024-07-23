package mintab

import (
	"log"
	"os"
)

// Pass by struct `Input`
func ExampleTable_Load_input() {
	data := Input{
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
	table := New(os.Stdout, WithMergeFields([]int{0, 1, 2, 3}))
	if err := table.Load(data); err != nil {
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

// Pass by any struct slices (with `CompressedTextFormat`)
func ExampleTable_Load_struct() {
	data := []struct {
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
	}{
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-1", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 22, ToPort: 22, AddressType: "SecurityGroup", CidrBlock: "sg-10"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-1", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-2", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 443, ToPort: 443, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-2", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "icmp", FromPort: -1, ToPort: -1, AddressType: "SecurityGroup", CidrBlock: "sg-11"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 3389, ToPort: 3389, AddressType: "Ipv4", CidrBlock: "10.1.0.0/16"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 0, ToPort: 65535, AddressType: "PrefixList", CidrBlock: "pl-id/pl-name"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
	}
	table := New(os.Stdout, WithFormat(CompressedTextFormat), WithMergeFields([]int{0, 1, 2, 3}))
	if err := table.Load(data); err != nil {
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
