package main

import (
	"fmt"
	"log"

	"github.com/nekrassov01/mintab"
)

func main() {
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

	table = mintab.NewTable(mintab.WithFormat(mintab.MarkdownFormat))
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
		| i-4        | server-4     | \-                         | \-           |
		| i-5        | server-5     | lb-domain-5                | \-           |
		| i-5        | server-5     | \-                         | tg-5<br>tg-6 |
	*/

	table = mintab.NewTable(mintab.WithFormat(mintab.BacklogFormat))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*
		| InstanceID | InstanceName | AttachedLB                 | AttachedTG   |h
		| i-1        | server-1     | lb-domain-1                | tg-1         |
		| i-2        | server-2     | lb-doamin-2&br;lb-doamin-3 | tg-2         |
		| i-3        | server-3     | lb-doamin-4                | tg-3&br;tg-4 |
		| i-4        | server-4     | \-                         | \-           |
		| i-5        | server-5     | lb-domain-5                | \-           |
		| i-5        | server-5     | \-                         | tg-5&br;tg-6 |
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

	table = mintab.NewTable(mintab.WithMergeFields([]int{0, 1}), mintab.WithTheme(mintab.DarkTheme))
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

	table = mintab.NewTable(mintab.WithIgnoreFields([]int{2}))
	if err := table.Load(s1); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*
		| InstanceID | InstanceName | AttachedTG |
		|------------|--------------|------------|
		| i-1        | server-1     | tg-1       |
		| i-2        | server-2     | tg-2       |
		| i-3        | server-3     | tg-3       |
		|            |              | tg-4       |
		| i-4        | server-4     | -          |
		| i-5        | server-5     | -          |
		| i-5        | server-5     | tg-5       |
		|            |              | tg-6       |
	*/

	type sample2 struct {
		Domain string
	}

	s2 := []sample2{
		{Domain: "*.example.com"},
	}

	table = mintab.NewTable(mintab.WithEscapeTargets([]string{"*"}))
	if err := table.Load(s2); err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Out())

	/*
		| Domain         |
		|----------------|
		| \*.example.com |
	*/
}
