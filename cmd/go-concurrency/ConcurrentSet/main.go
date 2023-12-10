// Package main contains an example of using the ConcurrentSet
package main

import (
	"fmt"
	"github.com/PavloVM7/go-concurrency/collections"
	"runtime"
	"sync"
	"time"
)

func main() {
	println("👉 Example of using ConcurrentSet")
	using := func(funcs string) {
		fmt.Println("=== using ", funcs)
	}
	set := collections.NewConcurrentSetCapacity[int](10)
	showSet := func() {
		fmt.Printf(">>> ConcurrentSet size: %d, elements: %v\n", set.Size(), set.ToSlice())
	}
	isSetEmpty := func() {
		fmt.Println("~~~ is set empty? -", set.IsEmpty())
	}
	isSetEmpty()

	using("AddAll()")
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if set.AddAll(values...) {
		showSet()
		isSetEmpty()
	}
	subvalues := values[1:8]
	if !set.AddAll(subvalues...) {
		fmt.Printf("- the set already contains values %v\n", subvalues)
	}
	using("Add()")
	val := 11
	if set.Add(val) {
		fmt.Printf("value %d was added to the set\n", val)
		showSet()
	}
	if !set.Add(val) {
		fmt.Println("- the set already contains the value", val)
	}
	using("Contains()")
	showSet()
	if set.Contains(3) {
		fmt.Println("+ the set contains the value 3")
	}
	if set.Contains(4) {
		fmt.Println("+ the set contains the value 4")
	}
	if !set.Contains(123) {
		fmt.Println("- there is no value 123 in the set")
	}

	using("Remove()")
	if set.Remove(3) {
		fmt.Printf("+ the value %d was removed from the set\n", 3)
	}
	if set.Remove(4) {
		fmt.Printf("+ the value %d was removed from the set\n", 4)
	}
	if !set.Remove(123) {
		fmt.Printf("- the value %d was not removed from the set because the set did not contain it\n", 123)
	}
	showSet()

	using("Clear()")
	set.Clear()
	showSet()
	isSetEmpty()

	using("TrimToSize()")
	const amount = 1_000_000
	fillSet(set, amount, 2)
	fmt.Println(">>> set size =", set.Size())

	getMemStats := func() runtime.MemStats {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		return mem
	}

	memToString := func(mem runtime.MemStats) string { return fmt.Sprintf("%d Kb", mem.Alloc/1024) }

	runtime.GC()

	fmt.Printf(">>> set size: %d, memory usage: %s\n", set.Size(), memToString(getMemStats()))

	removeValues(set, 21, amount, 3)

	runtime.GC()

	fmt.Printf("after removing memory usage: %s, set size: %d\n", memToString(getMemStats()), set.Size())
	showSet()

	set.TrimToSize()

	runtime.GC()

	fmt.Printf("after TrimToSize() memory usage: %s, set size: %d\n", memToString(getMemStats()), set.Size())
	showSet()
}

func fillSet(set *collections.ConcurrentSet[int], amount, threads int) {
	fmt.Printf("* filling set, amount: %d, threads: %d\n", amount, threads)
	start := time.Now()
	chStart := make(chan struct{})
	var wg sync.WaitGroup
	adds := make([]int, threads)
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			<-chStart
			n := 0
			for !set.Contains(amount) {
				n++
				if set.Add(n) {
					adds[num]++
				}
			}
		}(i)
	}
	close(chStart)
	wg.Wait()
	fmt.Printf(">>> the set was filled, duration: %v, amount: %d, threads: %d, each thread added: %v\n",
		time.Since(start), set.Size(), threads, adds)
}
func removeValues(set *collections.ConcurrentSet[int], start, end, threads int) {
	fmt.Printf("* remove values from set, from %d to %d , threads: %d\n", start, end, threads)
	st := time.Now()
	chStart := make(chan struct{})
	var wg sync.WaitGroup
	adds := make([]int, threads)
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			<-chStart
			for val := start; val <= end; val++ {
				if set.Remove(val) {
					adds[num]++
				}
			}
		}(i)
	}
	close(chStart)
	wg.Wait()
	fmt.Printf(">>> values were removed, duration: %v, set size: %d, threads: %d, each thread removed: %v\n",
		time.Since(st), set.Size(), threads, adds)
}
