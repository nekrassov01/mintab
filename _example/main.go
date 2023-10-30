package main

import (
	"fmt"
	"log"

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
	var str string
	var err error

	table, err = mintab.New(samples)
	if err != nil {
		log.Fatal(err)
	}
	str, err = table.Out()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)

	/*
		| InstanceName | SecurityGroupName | CidrBlock                |
		| ------------ | ----------------- | ------------------------ |
		| i-1          | sg-1              | 10.0.0.0/16              |
		| i-1          | sg-1              | 10.1.0.0/16              |
		| i-1          | sg-2              | 10.2.0.0/16              |
		| i-1          | sg-2              | 10.3.0.0/16              |
		| i-2          | sg-1              | 10.0.0.0/16<br>0.0.0.0/0 |
		| i-2          | sg-1              | 10.1.0.0/16<br>0.0.0.0/0 |
		| i-2          | sg-2              | 10.2.0.0/16<br>0.0.0.0/0 |
		| i-2          | sg-2              | 10.3.0.0/16<br>0.0.0.0/0 |
		| i-3          | N/A               | 10.0.0.0/16<br>0.0.0.0/0 |
		| i-4          | sg-4              | N/A                      |
	*/

	table, err = mintab.New(samples, mintab.WithTableFormat(mintab.Backlog))
	if err != nil {
		log.Fatal(err)
	}
	str, err = table.Out()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)

	/*
		| InstanceName | SecurityGroupName |        CidrBlock         |h
		| i-1          | sg-1              | 10.0.0.0/16              |
		| i-1          | sg-1              | 10.1.0.0/16              |
		| i-1          | sg-2              | 10.2.0.0/16              |
		| i-1          | sg-2              | 10.3.0.0/16              |
		| i-2          | sg-1              | 10.0.0.0/16&br;0.0.0.0/0 |
		| i-2          | sg-1              | 10.1.0.0/16&br;0.0.0.0/0 |
		| i-2          | sg-2              | 10.2.0.0/16&br;0.0.0.0/0 |
		| i-2          | sg-2              | 10.3.0.0/16&br;0.0.0.0/0 |
		| i-3          | N/A               | 10.0.0.0/16&br;0.0.0.0/0 |
		| i-4          | sg-4              | N/A                      |
	*/

	table, err = mintab.New(samples, mintab.WithTableHeader(false))
	if err != nil {
		log.Fatal(err)
	}
	str, err = table.Out()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)

	/*
		| i-1 | sg-1 | 10.0.0.0/16              |
		| i-1 | sg-1 | 10.1.0.0/16              |
		| i-1 | sg-2 | 10.2.0.0/16              |
		| i-1 | sg-2 | 10.3.0.0/16              |
		| i-2 | sg-1 | 10.0.0.0/16<br>0.0.0.0/0 |
		| i-2 | sg-1 | 10.1.0.0/16<br>0.0.0.0/0 |
		| i-2 | sg-2 | 10.2.0.0/16<br>0.0.0.0/0 |
		| i-2 | sg-2 | 10.3.0.0/16<br>0.0.0.0/0 |
		| i-3 | N/A  | 10.0.0.0/16<br>0.0.0.0/0 |
		| i-4 | sg-4 | N/A                      |
	*/

	table, err = mintab.New(samples, mintab.WithEmptyFieldPlaceholder("NULL"))
	if err != nil {
		log.Fatal(err)
	}
	str, err = table.Out()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)

	/*
		| InstanceName | SecurityGroupName | CidrBlock                |
		| ------------ | ----------------- | ------------------------ |
		| i-1          | sg-1              | 10.0.0.0/16              |
		| i-1          | sg-1              | 10.1.0.0/16              |
		| i-1          | sg-2              | 10.2.0.0/16              |
		| i-1          | sg-2              | 10.3.0.0/16              |
		| i-2          | sg-1              | 10.0.0.0/16<br>0.0.0.0/0 |
		| i-2          | sg-1              | 10.1.0.0/16<br>0.0.0.0/0 |
		| i-2          | sg-2              | 10.2.0.0/16<br>0.0.0.0/0 |
		| i-2          | sg-2              | 10.3.0.0/16<br>0.0.0.0/0 |
		| i-3          | NULL              | 10.0.0.0/16<br>0.0.0.0/0 |
		| i-4          | sg-4              | NULL                     |
	*/

	table, err = mintab.New(samples, mintab.WithWordDelimiter(","))
	if err != nil {
		log.Fatal(err)
	}
	str, err = table.Out()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)

	/*
		| InstanceName | SecurityGroupName | CidrBlock             |
		| ------------ | ----------------- | --------------------- |
		| i-1          | sg-1              | 10.0.0.0/16           |
		| i-1          | sg-1              | 10.1.0.0/16           |
		| i-1          | sg-2              | 10.2.0.0/16           |
		| i-1          | sg-2              | 10.3.0.0/16           |
		| i-2          | sg-1              | 10.0.0.0/16,0.0.0.0/0 |
		| i-2          | sg-1              | 10.1.0.0/16,0.0.0.0/0 |
		| i-2          | sg-2              | 10.2.0.0/16,0.0.0.0/0 |
		| i-2          | sg-2              | 10.3.0.0/16,0.0.0.0/0 |
		| i-3          | N/A               | 10.0.0.0/16,0.0.0.0/0 |
		| i-4          | sg-4              | N/A                   |
	*/

	table, err = mintab.New(samples, mintab.WithMergeFields([]int{0, 1}), mintab.WithTableTheme(mintab.DarkTheme))
	if err != nil {
		log.Fatal(err)
	}
	str, err = table.Out()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)

	/*
		| InstanceName | SecurityGroupName | CidrBlock                |
		| ------------ | ----------------- | ------------------------ |
		| i-1          | sg-1              | 10.0.0.0/16              |
		| i-1          |                   | 10.1.0.0/16              |
		| i-1          | sg-2              | 10.2.0.0/16              |
		| i-1          |                   | 10.3.0.0/16              |
		| i-2          | sg-1              | 10.0.0.0/16<br>0.0.0.0/0 |
		| i-2          |                   | 10.1.0.0/16<br>0.0.0.0/0 |
		| i-2          | sg-2              | 10.2.0.0/16<br>0.0.0.0/0 |
		| i-2          |                   | 10.3.0.0/16<br>0.0.0.0/0 |
		| i-3          | N/A               | 10.0.0.0/16<br>0.0.0.0/0 |
		| i-4          | sg-4              | N/A                      |
	*/

	table, err = mintab.New(samples, mintab.WithIgnoreFields([]int{2}))
	if err != nil {
		log.Fatal(err)
	}
	str, err = table.Out()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)

	/*
		| InstanceName | SecurityGroupName |
		| ------------ | ----------------- |
		| i-1          | sg-1              |
		| i-1          | sg-1              |
		| i-1          | sg-2              |
		| i-1          | sg-2              |
		| i-2          | sg-1              |
		| i-2          | sg-1              |
		| i-2          | sg-2              |
		| i-2          | sg-2              |
		| i-3          | N/A               |
		| i-4          | sg-4              |
	*/
}
