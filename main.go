package main

import (
	"fmt"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {
	fmt.Println("=== Testing Fair Queue ===")
	c := NewFairChan(50)
	exampleMode(c)

	fmt.Println("\n=== Testing FIFO Queue ===")
	c2 := NewFifoChan(50)
	exampleMode(c2)
}

func exampleMode(queue Queue) {
	showSteps(exampleSteps)
	processSteps(queue, exampleSteps, true)
}

func showSteps(steps [][]Message) {
	// 收集所有租户
	tenantSet := make(map[string]bool)
	for _, step := range steps {
		for _, msg := range step {
			tenantSet[msg.tenant] = true
		}
	}

	// 按字母顺序排序租户
	tenants := make([]string, 0, len(tenantSet))
	for tenant := range tenantSet {
		tenants = append(tenants, tenant)
	}
	for i := 0; i < len(tenants)-1; i++ {
		for j := i + 1; j < len(tenants); j++ {
			if tenants[i] > tenants[j] {
				tenants[i], tenants[j] = tenants[j], tenants[i]
			}
		}
	}

	// 计算每个step的最大写入次数
	maxWritesPerStep := make([]int, len(steps))
	for stepIdx, step := range steps {
		max := 0
		tenantCounts := make(map[string]int)
		for _, msg := range step {
			tenantCounts[msg.tenant]++
			if tenantCounts[msg.tenant] > max {
				max = tenantCounts[msg.tenant]
			}
		}
		maxWritesPerStep[stepIdx] = max
	}

	// 显示每个租户的输入情况
	fmt.Println("\nInput pattern:")
	for _, tenant := range tenants {
		fmt.Printf("%s: ", tenant)
		for stepIdx, step := range steps {
			// 计算该租户在当前step的写入次数
			writeCount := 0
			for _, msg := range step {
				if msg.tenant == tenant {
					writeCount++
				}
			}

			// 显示写入标记
			for i := 0; i < writeCount; i++ {
				fmt.Print("+")
			}
			// 显示空位标记
			for i := writeCount; i < maxWritesPerStep[stepIdx]; i++ {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func processSteps(queue Queue, steps [][]Message, continuous bool) {
	index := 0
	step := 0
	nextLog := make([]map[string]int, 0) // Track Next() calls per step

	for len(steps) > index || !queue.IsEmpty() {

		// put all steps in queue
		if index < len(steps) {
			item := steps[index]
			//fmt.Printf("Input: ")
			for _, message := range item {
				//fmt.Printf("%s:%d ", message.tenant, message.message)
				if err := queue.Put(message.tenant, message.message); err != nil {
					//fmt.Printf("(error: %v) ", err)
				}
			}
			//fmt.Println()
			index++
		}

		// show queue state before processing
		// Note: Graph() method not implemented in queue types

		// Initialize next log for this step
		stepNextLog := make(map[string]int)

		// read a next
		if !queue.IsEmpty() {
			t, message := queue.Next()
			_ = message
			stepNextLog[t]++
		} else {

		}

		nextLog = append(nextLog, stepNextLog)
		step++
	}

	// Display Next() call pattern similar to showSteps
	showNextPattern(nextLog)
}

func showNextPattern(nextLog []map[string]int) {
	fmt.Println("next call pattern")
	// Collect all tenants from Next() calls
	tenantSet := make(map[string]bool)
	for _, stepLog := range nextLog {
		for tenant := range stepLog {
			tenantSet[tenant] = true
		}
	}

	// Sort tenants alphabetically
	tentants := make([]string, 0, len(tenantSet))
	for tenant := range tenantSet {
		tentants = append(tentants, tenant)
	}
	for i := 0; i < len(tentants)-1; i++ {
		for j := i + 1; j < len(tentants); j++ {
			if tentants[i] > tentants[j] {
				tentants[i], tentants[j] = tentants[j], tentants[i]
			}
		}
	}

	// Find max Next() calls per step
	maxNextPerStep := make([]int, len(nextLog))
	for stepIdx, stepLog := range nextLog {
		max := 0
		for _, count := range stepLog {
			if count > max {
				max = count
			}
		}
		maxNextPerStep[stepIdx] = max
	}

	// Display Next() pattern for each tenant
	for _, tenant := range tentants {
		fmt.Printf("%s: ", tenant)
		for stepIdx, stepLog := range nextLog {
			nextCount := stepLog[tenant]

			// Display Next() calls
			for i := 0; i < nextCount; i++ {
				fmt.Print("+")
			}
			// Display empty slots
			for i := nextCount; i < maxNextPerStep[stepIdx]; i++ {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

var globalGenerator = Generator()

// exampleSteps output example
// showSteps result
// hello: ++++++
// xxx  : +...+.
//
// 对于一个 step 来说，每个租户都独立从0开始排序，由于有些租户在一个 step 中可以输入多次，其他租户后续则标记为没有输入。
var exampleSteps [][]Message = [][]Message{
	{
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t2"),
	},
	{
		globalGenerator("t1"),
		globalGenerator("t1"),
	},
	{}, // Empty step to test queue clearing
	{
		globalGenerator("t1"),
		globalGenerator("t2"),
		globalGenerator("t3"),
	},
	{
		globalGenerator("t3"),
	},
	{ // big input
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
	},
	{
		globalGenerator("t3"),
		globalGenerator("t3"),
		globalGenerator("t3"),
		globalGenerator("t2"),
		globalGenerator("t3"),
	},
	{
		globalGenerator("t1"),
		globalGenerator("t2"),
		globalGenerator("t3"),
	},
	{ // big input
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
		globalGenerator("t1"),
	},
}
