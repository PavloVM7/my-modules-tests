package main

import (
	"fmt"
	"github.com/PavloVM7/go-collections/pkg/collections"
	"sort"
)

func main() {
	src := map[int]string{1: "value 1", 2: "value 2", 3: "value 3"}
	fmt.Println("source map:")
	showIntMap(src)

	cpy := collections.CopyMap(src)
	fmt.Println("copy of the map:")
	showIntMap(cpy)

	src[4] = "value 4"
	fmt.Println("source map after adding new value:")
	showIntMap(src)

	fmt.Println("copy of the map:")
	showIntMap(cpy)
}
func showIntMap(mp map[int]string) {
	keys := make([]int, 0, len(mp))
	for k := range mp {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fmt.Printf("\t%d => %v", k, mp[k])
	}
	fmt.Println()
}
