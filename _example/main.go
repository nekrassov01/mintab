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
	samplesPtr := make([]*sample, 0, len(samples))
	for i := range samples {
		samplesPtr = append(samplesPtr, &samples[i])
	}
	slicePtr := &samples

	type num struct {
		Number      int
		NumberSlice []int
	}
	nums := []num{
		{Number: 0, NumberSlice: []int{0, 1, 2}},
		{Number: 1, NumberSlice: []int{}},
		{Number: -1, NumberSlice: []int{-1, 0, 1}},
	}

	type escaped struct {
		Domain string
	}
	escapes := []escaped{
		{Domain: "*.example.com"},
	}

	var table *mintab.Table
	table = mintab.NewTable()
	if err := table.Load(samples); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	if err := table.Load(samplesPtr); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	if err := table.Load(slicePtr); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*

		| InstanceName | SecurityGroupName | CidrBlock                |
		|--------------|-------------------|--------------------------|
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

	table = mintab.NewTable(mintab.WithFormat(mintab.BacklogFormat))
	if err := table.Load(samplesPtr); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	table = mintab.NewTable(mintab.WithFormat(mintab.BacklogFormat), mintab.WithTheme(mintab.DarkTheme))
	if err := table.Load(samples); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*

		| InstanceName | SecurityGroupName | CidrBlock                |h
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

	table = mintab.NewTable(mintab.WithHeader(false))
	if err := table.Load(samples); err != nil {
		log.Fatal(err)
	}
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

	table = mintab.NewTable(mintab.WithEmptyFieldPlaceholder("NULL"))
	if err := table.Load(samples); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*

		| InstanceName | SecurityGroupName | CidrBlock                |
		|--------------|-------------------|--------------------------|
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

	table = mintab.NewTable(mintab.WithWordDelimiter(","))
	if err := table.Load(samples); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*

		| InstanceName | SecurityGroupName | CidrBlock             |
		|--------------|-------------------|-----------------------|
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

	table = mintab.NewTable(mintab.WithMergeFields([]int{0, 1}))
	if err := table.Load(samples); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*

		| InstanceName | SecurityGroupName | CidrBlock                |
		|--------------|-------------------|--------------------------|
		| i-1          | sg-1              | 10.0.0.0/16              |
		|              |                   | 10.1.0.0/16              |
		|              | sg-2              | 10.2.0.0/16              |
		|              |                   | 10.3.0.0/16              |
		| i-2          | sg-1              | 10.0.0.0/16<br>0.0.0.0/0 |
		|              |                   | 10.1.0.0/16<br>0.0.0.0/0 |
		|              | sg-2              | 10.2.0.0/16<br>0.0.0.0/0 |
		|              |                   | 10.3.0.0/16<br>0.0.0.0/0 |
		| i-3          | N/A               | 10.0.0.0/16<br>0.0.0.0/0 |
		| i-4          | sg-4              | N/A                      |

	*/

	table = mintab.NewTable(mintab.WithMergeFields([]int{0, 1}), mintab.WithTheme(mintab.LightTheme))
	if err := table.Load(samples); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*

		| InstanceName | SecurityGroupName | CidrBlock                |
		|--------------|-------------------|--------------------------|
		| i-1          | sg-1              | 10.0.0.0/16              |
		|              |                   | 10.1.0.0/16              |
		|              | sg-2              | 10.2.0.0/16              |
		|              |                   | 10.3.0.0/16              |
		| i-2          | sg-1              | 10.0.0.0/16<br>0.0.0.0/0 |
		|              |                   | 10.1.0.0/16<br>0.0.0.0/0 |
		|              | sg-2              | 10.2.0.0/16<br>0.0.0.0/0 |
		|              |                   | 10.3.0.0/16<br>0.0.0.0/0 |
		| i-3          | N/A               | 10.0.0.0/16<br>0.0.0.0/0 |
		| i-4          | sg-4              | N/A                      |

	*/

	table = mintab.NewTable(mintab.WithIgnoreFields([]int{2}))
	if err := table.Load(samples); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*

		| InstanceName | SecurityGroupName |
		|--------------|-------------------|
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

	table = mintab.NewTable()
	if err := table.Load(nums); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*

		| Number | NumberSlice  |
		|--------|--------------|
		|      0 | 0<br>1<br>2  |
		|      1 | N/A          |
		|     -1 | -1<br>0<br>1 |

	*/

	table = mintab.NewTable(mintab.WithFormat(mintab.BacklogFormat))
	if err := table.Load(nums); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*

		| Number | NumberSlice  |h
		|      0 | 0&br;1&br;2  |
		|      1 | N/A          |
		|     -1 | -1&br;0&br;1 |

	*/

	table = mintab.NewTable(mintab.WithEscapeTargets([]string{"*"}))
	if err := table.Load(escapes); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*

		| Domain         |
		|----------------|
		| \*.example.com |

	*/
}
