package main

import (
	"fmt"

	"github.com/nekrassov01/mintab"
)

func main() {
	type sample struct {
		InstanceName      string
		SecurityGroupName string
		CidrBlock         []string
	}

	samples := []sample{
		{InstanceName: "i-1", SecurityGroupName: "sg-1", CidrBlock: []string{"10.0.0.0/16"}},
		{InstanceName: "i-1", SecurityGroupName: "sg-1", CidrBlock: []string{"10.1.0.0/16"}},
		{InstanceName: "i-1", SecurityGroupName: "sg-2", CidrBlock: []string{"10.2.0.0/16"}},
		{InstanceName: "i-1", SecurityGroupName: "sg-2", CidrBlock: []string{"10.3.0.0/16"}},
		{InstanceName: "i-2", SecurityGroupName: "sg-1", CidrBlock: []string{"10.0.0.0/16", "0.0.0.0/0"}},
		{InstanceName: "i-2", SecurityGroupName: "sg-1", CidrBlock: []string{"10.1.0.0/16", "0.0.0.0/0"}},
		{InstanceName: "i-2", SecurityGroupName: "sg-2", CidrBlock: []string{"10.2.0.0/16", "0.0.0.0/0"}},
		{InstanceName: "i-2", SecurityGroupName: "sg-2", CidrBlock: []string{"10.3.0.0/16", "0.0.0.0/0"}},
		{InstanceName: "i-3", SecurityGroupName: "", CidrBlock: []string{"10.0.0.0/16", "0.0.0.0/0"}},
		{InstanceName: "i-4", SecurityGroupName: "sg-4", CidrBlock: []string{}},
	}

	var table *mintab.Table
	table = mintab.New(samples)
	fmt.Println(table.Out())

	table = mintab.New(samples, mintab.WithTableFormat(mintab.Backlog))
	fmt.Println(table.Out())

	table = mintab.New(samples, mintab.WithTableHeader(false))
	fmt.Println(table.Out())

	table = mintab.New(samples, mintab.WithEmptyFieldPlaceholder("NULL"))
	fmt.Println(table.Out())

	table = mintab.New(samples, mintab.WithEmptyFieldPlaceholder(""))
	fmt.Println(table.Out())

	table = mintab.New(samples, mintab.WithWordDelimitter(","))
	fmt.Println(table.Out())

	table = mintab.New(samples, mintab.WithMergeFields([]int{0, 1}), mintab.WithTableTheme(mintab.DarkTheme))
	fmt.Println(table.Out())

	table = mintab.New(samples, mintab.WithMergeFields([]int{0, 1}), mintab.WithTableTheme(mintab.DarkTheme), mintab.WithEmptyFieldPlaceholder(""))
	fmt.Println(table.Out())

	table = mintab.New(samples, mintab.WithMergeFields([]int{0, 1}), mintab.WithTableTheme(mintab.LightTheme))
	fmt.Println(table.Out())

	table = mintab.New(samples, mintab.WithIgnoreFields([]int{2}))
	fmt.Println(table.Out())

	/* ignored except for slice of struct */

	table = mintab.New([]string{"aaa", "bbb", "ccc"})
	fmt.Println(table.Out())

	table = mintab.New([]int{1, 2, 3})
	fmt.Println(table.Out())

	table = mintab.New([]bool{true, false, true})
	fmt.Println(table.Out())

	table = mintab.New([]rune{'a', 'b', 'c'})
	fmt.Println(table.Out())

	table = mintab.New("aaa")
	fmt.Println(table.Out())

	table = mintab.New(1)
	fmt.Println(table.Out())

	table = mintab.New(true)
	fmt.Println(table.Out())

	table = mintab.New('a')
	fmt.Println(table.Out())
}
