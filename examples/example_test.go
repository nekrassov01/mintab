package examples

import (
	"log"
	"os"

	"github.com/nekrassov01/mintab"
)

// Pass by struct `Input`
func ExampleTable_Load_input() {
	data := mintab.Input{
		Header: []string{"Instance", "SG", "Direction", "Protocol", "From", "To", "AddressType", "CidrBlock"},
		Data: [][]any{
			{"i-1", "sg-1", "Ingress", "tcp", 22, 22, "SG", "sg-10"},
			{"i-1", "sg-1", "Egress", "-1", 0, 0, "Ipv4", "0.0.0.0/0"},
			{"i-1", "sg-2", "Ingress", "tcp", 443, 443, "Ipv4", "0.0.0.0/0"},
			{"i-1", "sg-2", "Egress", "-1", 0, 0, "Ipv4", "0.0.0.0/0"},
			{"i-2", "sg-3", "Ingress", "icmp", -1, -1, "SG", "sg-11"},
			{"i-2", "sg-3", "Ingress", "tcp", 3389, 3389, "Ipv4", "10.1.0.0/16"},
			{"i-2", "sg-3", "Ingress", "tcp", 0, 65535, "PrefixList", "pl-id/pl-name"},
			{"i-2", "sg-3", "Egress", "-1", 0, 0, "Ipv4", "0.0.0.0/0"},
		},
	}
	table := mintab.New(os.Stdout, mintab.WithMergeFields([]int{0, 1, 2, 3}))
	if err := table.Load(data); err != nil {
		log.Fatal(err)
	}
	table.Render()

	// Output:
	// +----------+------+-----------+----------+------+-------+-------------+---------------+
	// | Instance | SG   | Direction | Protocol | From | To    | AddressType | CidrBlock     |
	// +----------+------+-----------+----------+------+-------+-------------+---------------+
	// | i-1      | sg-1 | Ingress   | tcp      |   22 |    22 | SG          | sg-10         |
	// +          +      +-----------+----------+------+-------+-------------+---------------+
	// |          |      | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
	// +          +------+-----------+----------+------+-------+-------------+---------------+
	// |          | sg-2 | Ingress   | tcp      |  443 |   443 | Ipv4        | 0.0.0.0/0     |
	// +          +      +-----------+----------+------+-------+-------------+---------------+
	// |          |      | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
	// +----------+------+-----------+----------+------+-------+-------------+---------------+
	// | i-2      | sg-3 | Ingress   | icmp     |   -1 |    -1 | SG          | sg-11         |
	// +          +      +           +----------+------+-------+-------------+---------------+
	// |          |      |           | tcp      | 3389 |  3389 | Ipv4        | 10.1.0.0/16   |
	// +          +      +           +          +------+-------+-------------+---------------+
	// |          |      |           |          |    0 | 65535 | PrefixList  | pl-id/pl-name |
	// +          +      +-----------+----------+------+-------+-------------+---------------+
	// |          |      | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
	// +----------+------+-----------+----------+------+-------+-------------+---------------+
}

// Pass by any struct slices (with `CompressedTextFormat`)
func ExampleTable_Load_struct() {
	data := []struct {
		Instance    string
		SG          string
		Direction   string
		Protocol    string
		From        int
		To          int
		AddressType string
		CidrBlock   string
	}{
		{Instance: "i-1", SG: "sg-1", Direction: "Ingress", Protocol: "tcp", From: 22, To: 22, AddressType: "SG", CidrBlock: "sg-10"},
		{Instance: "i-1", SG: "sg-1", Direction: "Egress", Protocol: "-1", From: 0, To: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{Instance: "i-1", SG: "sg-2", Direction: "Ingress", Protocol: "tcp", From: 443, To: 443, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{Instance: "i-1", SG: "sg-2", Direction: "Egress", Protocol: "-1", From: 0, To: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{Instance: "i-2", SG: "sg-3", Direction: "Ingress", Protocol: "icmp", From: -1, To: -1, AddressType: "SG", CidrBlock: "sg-11"},
		{Instance: "i-2", SG: "sg-3", Direction: "Ingress", Protocol: "tcp", From: 3389, To: 3389, AddressType: "Ipv4", CidrBlock: "10.1.0.0/16"},
		{Instance: "i-2", SG: "sg-3", Direction: "Ingress", Protocol: "tcp", From: 0, To: 65535, AddressType: "PrefixList", CidrBlock: "pl-id/pl-name"},
		{Instance: "i-2", SG: "sg-3", Direction: "Egress", Protocol: "-1", From: 0, To: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
	}
	table := mintab.New(os.Stdout, mintab.WithFormat(mintab.CompressedTextFormat), mintab.WithMergeFields([]int{0, 1, 2, 3}))
	if err := table.Load(data); err != nil {
		log.Fatal(err)
	}
	table.Render()

	// Output:
	// +----------+------+-----------+----------+------+-------+-------------+---------------+
	// | Instance | SG   | Direction | Protocol | From | To    | AddressType | CidrBlock     |
	// +----------+------+-----------+----------+------+-------+-------------+---------------+
	// | i-1      | sg-1 | Ingress   | tcp      |   22 |    22 | SG          | sg-10         |
	// |          |      | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
	// |          | sg-2 | Ingress   | tcp      |  443 |   443 | Ipv4        | 0.0.0.0/0     |
	// |          |      | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
	// +----------+------+-----------+----------+------+-------+-------------+---------------+
	// | i-2      | sg-3 | Ingress   | icmp     |   -1 |    -1 | SG          | sg-11         |
	// |          |      |           | tcp      | 3389 |  3389 | Ipv4        | 10.1.0.0/16   |
	// |          |      |           |          |    0 | 65535 | PrefixList  | pl-id/pl-name |
	// |          |      | Egress    |       -1 |    0 |     0 | Ipv4        | 0.0.0.0/0     |
	// +----------+------+-----------+----------+------+-------+-------------+---------------+
}
