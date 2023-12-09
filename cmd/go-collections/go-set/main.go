package main

import (
	"fmt"
	"github.com/PavloVM7/go-collections/pkg/collections"
	"runtime"
)

func main() {
	set := collections.NewSetCapacity[int](3)
	using := func(funcs string) {
		fmt.Printf("=== using %s\n", funcs)
	}
	showSet := func() {
		fmt.Printf(">>> set capacity: %d, size: %d, elements: %v\n", set.Capacity(), set.Size(), set.ToSlice())
	}
	isSetEmpty := func() {
		fmt.Println("~~~ is set empty? -", set.IsEmpty())
	}
	showSet()
	isSetEmpty()
	using("AddAll()")
	if set.AddAll(1, 2, 3) {
		showSet()
		isSetEmpty()
	}
	if !set.AddAll(2, 3) {
		fmt.Println("- the set already contains values 2 and 3")
	}
	using("Add()")
	if set.Add(4) {
		fmt.Println("+ ", 4, "was added to the set")
		showSet()
	}
	if !set.Add(1) {
		fmt.Println("- the set already contains the value 1")
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
	const number = 1_000_000
	for i := 1; i <= number; i++ {
		set.Add(i)
	}

	getMemStats := func() runtime.MemStats {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		return mem
	}

	memToString := func(mem runtime.MemStats) string { return fmt.Sprintf("%d Kb", mem.Alloc/1024) }

	runtime.GC()

	fmt.Printf(">>> set capacity: %d, size: %d, memory usage: %s\n", set.Capacity(), set.Size(), memToString(getMemStats()))
	for i := 21; i <= number; i++ {
		set.Remove(i)
	}

	runtime.GC()

	fmt.Printf("after removing memory usage: %s, set size: %d\n", memToString(getMemStats()), set.Size())
	showSet()

	set.TrimToSize()

	runtime.GC()

	fmt.Printf("after TrimToSize() memory usage: %s, set size: %d\n", memToString(getMemStats()), set.Size())
	showSet()
}
