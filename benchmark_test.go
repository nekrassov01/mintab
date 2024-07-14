package mintab

import (
	"bytes"
	"testing"

	"github.com/olekukonko/tablewriter"
)

func BenchmarkMintab(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		table := New(&bytes.Buffer{})
		if err := table.Load(s1); err != nil {
			b.Fatal(err)
		}
		table.Render()
	}
}

func BenchmarkMintabSimple(b *testing.B) {
	data := []struct {
		InstanceID    string
		InstanecName  string
		InstanceState string
	}{
		{InstanceID: "i-1", InstanecName: "server-1", InstanceState: "running"},
		{InstanceID: "i-2", InstanecName: "server-2", InstanceState: "stopped"},
		{InstanceID: "i-3", InstanecName: "server-3", InstanceState: "pending"},
		{InstanceID: "i-4", InstanecName: "server-4", InstanceState: "terminated"},
		{InstanceID: "i-5", InstanecName: "server-5", InstanceState: "stopping"},
		{InstanceID: "i-6", InstanecName: "server-6", InstanceState: "shutting-down"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		table := New(&bytes.Buffer{})
		if err := table.Load(data); err != nil {
			b.Fatal(err)
		}
		table.Render()
	}
}

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
