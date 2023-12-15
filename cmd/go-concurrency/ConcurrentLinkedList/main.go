package main

import (
	"fmt"
	"github.com/PavloVM7/go-concurrency/collections"
	"sync"
)

func main() {
	println("👉 Example of using ConcurrentLinkedList")
	list := collections.NewConcurrentLinkedList[int]()
	using := func(funcs string) {
		fmt.Printf("=== using %s\n", funcs)
	}
	var wg sync.WaitGroup
	using("AddLast() and AddFirst()")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 11; i <= 20; i++ {
			list.AddLast(i) // adds items to the end of the list
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 10; i > 0; i-- {
			list.AddFirst(i) // adds items to the head of the list
		}
	}()
	wg.Wait()
	showList := func() {
		fmt.Printf(">>> list size: %d, items: %v\n", list.Size(), list.ToArray())
	}
	showList()

	using("Get() and Remove()")
	item10, err := list.Get(10)
	fmt.Printf("before remove 10th item = %d, err = %v\n", item10, err)
	item10, err = list.Remove(10) // removes 10th item
	fmt.Printf("removed item10 = %d, err = %v\n", item10, err)
	item10, err = list.Get(10)
	fmt.Printf("after remove 10th item = %d, err = %v\n", item10, err)
	showList()

	using("GetFirst() and RemoveFirst()")
	first, firstOk := list.GetFirst()
	fmt.Printf("before remove first element: %d, exists: %t\n", first, firstOk)
	first, firstOk = list.RemoveFirst()
	fmt.Printf("first element: %d, removed: %t\n", first, firstOk)
	first, firstOk = list.GetFirst()
	fmt.Printf("current first element: %d, exists: %t\n", first, firstOk)
	showList()

	using("GetLast() and RemoveLast()")
	last, lastOk := list.GetLast()
	fmt.Printf("before remove last element: %d, exists: %t\n", last, lastOk)
	last, lastOk = list.RemoveLast()
	fmt.Printf("last element: %d, removed: %t\n", last, lastOk)
	last, lastOk = list.GetLast()
	fmt.Printf("current last element: %d, exists: %t\n", last, lastOk)
	showList()

	using("RemoveFirstOccurrence()")
	rFirst, fIndex := list.RemoveFirstOccurrence(func(value int) bool {
		return value%2 != 0
	})
	fmt.Printf("removed first odd value: %d, index: %d\n", rFirst, fIndex)
	showList()

	using("RemoveLastOccurrence()")
	rLast, lIndex := list.RemoveLastOccurrence(func(value int) bool {
		return value%2 == 0
	})
	fmt.Printf("removed last even value: %d, index: %d\n", rLast, lIndex)
	showList()

	using("RemoveAll()")
	count := list.RemoveAll(func(value int) bool {
		return value%3 == 0
	})
	fmt.Printf("%d elements that are dividable by 3 have been removed\n", count)
	showList()

	using("Clear()")
	list.Clear()
	showList()
}
