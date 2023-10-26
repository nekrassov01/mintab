mintab
======

[![CI](https://github.com/nekrassov01/mintab/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/nekrassov01/mintab/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/nekrassov01/mintab/graph/badge.svg?token=RIV62CQILM)](https://codecov.io/gh/nekrassov01/mintab)
[![Go Reference](https://pkg.go.dev/badge/github.com/nekrassov01/mintab.svg)](https://pkg.go.dev/github.com/nekrassov01/mintab)
[![Go Report Card](https://goreportcard.com/badge/github.com/nekrassov01/mintab)](https://goreportcard.com/report/github.com/nekrassov01/mintab)

mintab is a minimum ASCII table utilities using [tablewriter](https://github.com/olekukonko/tablewriter)

![terminal](_assets/terminal.png)

Support
---------

- Markdown table format
- Backlog table format
- Group columns based on first field value
- Color rows based on first field value
- Ignore specified columns

Notes
-----

- Only slice of struct is accepted
- Using reflect

Usage
-----

```go
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

	table = mintab.New(samples, mintab.WithTableFormat(mintab.Backlog))
	fmt.Println(table.Out())

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

	table = mintab.New(samples, mintab.WithDisableHeader())
	fmt.Println(table.Out())

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

	table = mintab.New(samples, mintab.WithEmptyFieldPlaceholder("NULL"))
	fmt.Println(table.Out())

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

	table = mintab.New(samples, mintab.WithWordDelimitter(","))
	fmt.Println(table.Out())

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

	table = mintab.New(samples, mintab.WithMergeFields([]int{0, 1}), mintab.WithTableTheme(mintab.DarkTheme))
	fmt.Println(table.Out())

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

	table = mintab.New(samples, mintab.WithIgnoreFields([]int{2}))
	fmt.Println(table.Out())

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
```

Author
------

[nekrassov01](https://github.com/nekrassov01)

License
-------

[MIT](https://github.com/nekrassov01/mintab/blob/main/LICENSE)
