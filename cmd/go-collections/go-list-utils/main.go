package main

import (
	"fmt"
	"github.com/PavloVM7/go-collections/pkg/collections/lists"
)

func main() {
	list := lists.NewLinkedListItems[int](10, 8, 6, 4, 2, 1, 3, 5, 7, 9)
	fmt.Printf("before sorting the list: %v\n", list.ToArray())
	lists.SortList(list, func(item1, item2 int) bool {
		return item1 < item2
	})
	fmt.Printf("after sorting the list:  %v\n", list.ToArray())
}
