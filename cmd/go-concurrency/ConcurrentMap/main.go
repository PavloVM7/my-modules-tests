// Package main contains an example of using the ConcurrentMap
package main

import (
	"fmt"
	"github.com/PavloVM7/go-concurrency/collections"
	"runtime"
	"strings"
	"sync"
	"time"
)

func showMap[K comparable, V any](mp map[K]V) string {
	var res strings.Builder
	res.WriteRune('[')
	for k, v := range mp {
		if res.Len() > 1 {
			res.WriteString(", ")
		}
		res.WriteString(fmt.Sprintf("['%v' => '%v']", k, v))
	}
	res.WriteRune(']')
	return res.String()
}

func createStringValue(i int) string {
	return fmt.Sprintf("value %d", i)
}

func main() {
	println("👉 Example of using ConcurrentMap")
	using := func(funcs string) {
		fmt.Println("=== using ", funcs)
	}

	cmp := collections.NewConcurrentMapCapacity[int, string](10)

	showKeys := func() {
		fmt.Println("map keys:", cmp.Keys())
	}
	showCurMap := func() {
		its := showMap(cmp.Copy())
		fmt.Printf(">>> ConcurrentMap size: %d, entities: %v\n", cmp.Size(), its)
	}
	showCurMap()
	isMapEmpty := func() {
		fmt.Println("~~~ is map empty? -", cmp.IsEmpty())
	}
	isMapEmpty()

	using("Put() and Get()")
	key := 1
	cmp.Put(key, "value 1")
	value, ok := cmp.Get(key)
	fmt.Printf("+ %d => '%s', exists: %t\n", key, value, ok)
	showCurMap()
	cmp.Put(key, "other value 1")
	value, ok = cmp.Get(key)
	fmt.Printf("+ %d => '%s', exists: %t\n", key, value, ok)
	showCurMap()
	isMapEmpty()

	using("PutIfNotExists() and Keys()")
	key = 2
	ok, value = cmp.PutIfNotExists(key, "value 2")
	fmt.Printf("+ %d => '%s', added: %t\n", key, value, ok)
	ok, value = cmp.PutIfNotExists(key, "other value 2")
	fmt.Printf("- %d => '%s', added: %t\n", key, value, ok)
	for _, key = range []int{3, 4, 5} {
		cmp.PutIfNotExists(key, createStringValue(key))
	}
	showCurMap()
	fmt.Printf("keys: %v\n", cmp.Keys())

	using("Remove()")
	cmp.Remove(4)
	cmp.Remove(123)
	showKeys()
	showCurMap()

	using("RemoveIfExists()")
	ok, value = cmp.RemoveIfExists(5)
	fmt.Printf("+ key: %d, value: '%s', removed: %t\n", 5, value, ok)
	showKeys()
	ok, value = cmp.RemoveIfExists(5)
	fmt.Printf("- key: %d, value: '%s', removed: %t\n", 5, value, ok)
	showKeys()
	showCurMap()

	using("ForEachRead()")
	sum := 0
	var sb strings.Builder
	cmp.ForEachRead(func(key int, value string) {
		sum += key
		if sb.Len() > 0 {
			sb.WriteRune(',')
		}
		sb.WriteRune('\'')
		sb.WriteString(value)
		sb.WriteRune('\'')
	})
	fmt.Printf("sum of keys: %d, all values: \"%s\"\n", sum, sb.String())

	using("Clear()")
	fmt.Println("= before clearing")
	showCurMap()
	showKeys()
	isMapEmpty()
	cmp.Clear()
	fmt.Println("= after clearing")
	showCurMap()
	showKeys()
	isMapEmpty()

	using("TrimToSize()")
	const amount = 1_000_000
	fillMap(cmp, amount, 3)
	fmt.Println(">>> map size:", cmp.Size())

	getMemStats := func() runtime.MemStats {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		return mem
	}

	memToString := func(mem runtime.MemStats) string { return fmt.Sprintf("%d Kb", mem.Alloc/1024) }

	runtime.GC()

	fmt.Printf(">>> map size: %d, memory usage: %s\n", cmp.Size(), memToString(getMemStats()))

	removeValues(cmp, 6, amount, 3)

	runtime.GC()

	fmt.Printf("after removing memory usage: %s, map size: %d\n", memToString(getMemStats()), cmp.Size())
	showKeys()

	cmp.TrimToSize()

	runtime.GC()

	fmt.Printf("after TrimToSize() memory usage: %s, map size: %d\n", memToString(getMemStats()), cmp.Size())
	showKeys()
	showCurMap()
}

func fillMap(mp *collections.ConcurrentMap[int, string], amount, threads int) {
	fmt.Printf("* filling map, amount: %d, threads: %d\n", amount, threads)
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
			_, ok := mp.Get(amount)
			for !ok {
				n++
				if aok, _ := mp.PutIfNotExists(n, createStringValue(n)); aok {
					adds[num]++
				}
				_, ok = mp.Get(amount)
			}
		}(i)
	}
	close(chStart)
	wg.Wait()
	fmt.Printf(">>> the map was filled, duration: %v, amount: %d, threads: %d, each thread added: %v\n",
		time.Since(start), mp.Size(), threads, adds)
}
func removeValues(mp *collections.ConcurrentMap[int, string], start, end, threads int) {
	fmt.Printf("* remove values from map, from %d to %d , threads: %d\n", start, end, threads)
	st := time.Now()
	chStart := make(chan struct{})
	var wg sync.WaitGroup
	adds := make([]int, threads)
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			<-chStart
			for key := start; key <= end; key++ {
				ok, _ := mp.RemoveIfExistsDoubleCheck(key)
				if ok {
					adds[num]++
				}
			}
		}(i)
	}
	close(chStart)
	wg.Wait()
	fmt.Printf(">>> values were removed, duration: %v, map size: %d, threads: %d, each thread removed: %v\n",
		time.Since(st), mp.Size(), threads, adds)
}
