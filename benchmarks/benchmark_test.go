package benchmarks

import (
	"bytes"
	"testing"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/nekrassov01/mintab"
	"github.com/olekukonko/tablewriter"
)

func BenchmarkMintabInput(b *testing.B) {
	data := mintab.Input{
		Header: []string{"InstanceID", "InstanceName", "InstanceState"},
		Data: [][]any{
			{"i-1", "server-1", "running"},
			{"i-2", "server-2", "stopped"},
			{"i-3", "server-3", "pending"},
			{"i-4", "server-4", "terminated"},
			{"i-5", "server-5", "stopping"},
			{"i-6", "server-6", "shutting-down"},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t := mintab.New(&bytes.Buffer{})
		if err := t.Load(data); err != nil {
			b.Fatal(err)
		}
		t.Render()
	}
}

func BenchmarkMintabStruct(b *testing.B) {
	data := []struct {
		InstanceID    string
		InstaneceName string
		InstanceState string
	}{
		{InstanceID: "i-1", InstaneceName: "server-1", InstanceState: "running"},
		{InstanceID: "i-2", InstaneceName: "server-2", InstanceState: "stopped"},
		{InstanceID: "i-3", InstaneceName: "server-3", InstanceState: "pending"},
		{InstanceID: "i-4", InstaneceName: "server-4", InstanceState: "terminated"},
		{InstanceID: "i-5", InstaneceName: "server-5", InstanceState: "stopping"},
		{InstanceID: "i-6", InstaneceName: "server-6", InstanceState: "shutting-down"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t := mintab.New(&bytes.Buffer{})
		if err := t.Load(data); err != nil {
			b.Fatal(err)
		}
		t.Render()
	}
}

/*

func BenchmarkMintabInputLarge(b *testing.B) {
	data := mintab.Input{
		Header: []string{
			"InstanceID",
			"InstanceName",
			"VPCID",
			"SecurityGroupID",
			"FlowDirection",
			"IPProtocol",
			"FromPort",
			"ToPort",
			"AddressType",
			"CidrBlock",
		},
		Data: [][]any{
			{"i-1", "server-1", "vpc-1", "sg-1", "Ingress", "tcp", 22, 22, "SecurityGroup", "sg-10"},
			{"i-1", "server-1", "vpc-1", "sg-1", "Egress", "-1", 0, 0, "Ipv4", "0.0.0.0/0"},
			{"i-1", "server-1", "vpc-1", "sg-2", "Ingress", "tcp", 443, 443, "Ipv4", "0.0.0.0/0"},
			{"i-1", "server-1", "vpc-1", "sg-2", "Egress", "-1", 0, 0, "Ipv4", "0.0.0.0/0"},
			{"i-2", "server-2", "vpc-1", "sg-3", "Ingress", "icmp", -1, -1, "SecurityGroup", "sg-11"},
			{"i-2", "server-2", "vpc-1", "sg-3", "Ingress", "tcp", 3389, 3389, "Ipv4", "10.1.0.0/16"},
			{"i-2", "server-2", "vpc-1", "sg-3", "Ingress", "tcp", 0, 65535, "PrefixList", "pl-id/pl-name"},
			{"i-2", "server-2", "vpc-1", "sg-3", "Egress", "-1", 0, 0, "Ipv4", "0.0.0.0/0"},
			{"i-1", "server-1", "vpc-1", "sg-1", "Ingress", "tcp", 22, 22, "SecurityGroup", "sg-10"},
			{"i-1", "server-1", "vpc-1", "sg-1", "Egress", "-1", 0, 0, "Ipv4", "0.0.0.0/0"},
			{"i-1", "server-1", "vpc-1", "sg-2", "Ingress", "tcp", 443, 443, "Ipv4", "0.0.0.0/0"},
			{"i-1", "server-1", "vpc-1", "sg-2", "Egress", "-1", 0, 0, "Ipv4", "0.0.0.0/0"},
			{"i-2", "server-2", "vpc-1", "sg-3", "Ingress", "icmp", -1, -1, "SecurityGroup", "sg-11"},
			{"i-2", "server-2", "vpc-1", "sg-3", "Ingress", "tcp", 3389, 3389, "Ipv4", "10.1.0.0/16"},
			{"i-2", "server-2", "vpc-1", "sg-3", "Ingress", "tcp", 0, 65535, "PrefixList", "pl-id/pl-name"},
			{"i-2", "server-2", "vpc-1", "sg-3", "Egress", "-1", 0, 0, "Ipv4", "0.0.0.0/0"},
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t := mintab.New(&bytes.Buffer{}, mintab.WithMergeFields([]int{1, 2, 3}))
		if err := t.Load(data); err != nil {
			b.Fatal(err)
		}
		t.Render()
	}
}

func BenchmarkMintabStructLarge(b *testing.B) {
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
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-1", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 22, ToPort: 22, AddressType: "SecurityGroup", CidrBlock: "sg-10"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-1", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-2", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 443, ToPort: 443, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-2", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "icmp", FromPort: -1, ToPort: -1, AddressType: "SecurityGroup", CidrBlock: "sg-11"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 3389, ToPort: 3389, AddressType: "Ipv4", CidrBlock: "10.1.0.0/16"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 0, ToPort: 65535, AddressType: "PrefixList", CidrBlock: "pl-id/pl-name"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-1", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 22, ToPort: 22, AddressType: "SecurityGroup", CidrBlock: "sg-10"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-1", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-2", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 443, ToPort: 443, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-1", InstanceName: "server-1", VPCID: "vpc-1", SecurityGroupID: "sg-2", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "icmp", FromPort: -1, ToPort: -1, AddressType: "SecurityGroup", CidrBlock: "sg-11"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 3389, ToPort: 3389, AddressType: "Ipv4", CidrBlock: "10.1.0.0/16"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Ingress", IPProtocol: "tcp", FromPort: 0, ToPort: 65535, AddressType: "PrefixList", CidrBlock: "pl-id/pl-name"},
		{InstanceID: "i-2", InstanceName: "server-2", VPCID: "vpc-1", SecurityGroupID: "sg-3", FlowDirection: "Egress", IPProtocol: "-1", FromPort: 0, ToPort: 0, AddressType: "Ipv4", CidrBlock: "0.0.0.0/0"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t := mintab.New(&bytes.Buffer{}, mintab.WithMergeFields([]int{1, 2, 3}))
		if err := t.Load(data); err != nil {
			b.Fatal(err)
		}
		t.Render()
	}
}

*/

func BenchmarkTableWriter(b *testing.B) {
	data := [][]string{
		{"i-1", "server-1", "running"},
		{"i-2", "server-2", "stopped"},
		{"i-3", "server-3", "pending"},
		{"i-4", "server-4", "terminated"},
		{"i-5", "server-5", "stopping"},
		{"i-6", "server-6", "shutting-down"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		table := tablewriter.NewWriter(&bytes.Buffer{})
		table.AppendBulk(data)
		table.Render()
	}
}

func BenchmarkGoPrettyTable(b *testing.B) {
	data := []table.Row{
		{"i-1", "server-1", "running"},
		{"i-2", "server-2", "stopped"},
		{"i-3", "server-3", "pending"},
		{"i-4", "server-4", "terminated"},
		{"i-5", "server-5", "stopping"},
		{"i-6", "server-6", "shutting-down"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t := table.NewWriter()
		t.SetOutputMirror(&bytes.Buffer{})
		t.AppendRows(data)
		t.Render()
	}
}
