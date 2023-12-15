package main

import (
	"fmt"
	"github.com/PavloVM7/go-concurrency/caches"
	"strings"
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

func main() {
	println("👉 Example of using LRU cache")

	lruCache := caches.NewLRU[int, string](5)
	fmt.Println("cache:", lruCache)
	showLruCache := func() {
		items := showMap(lruCache.Copy())
		fmt.Printf(">>> cache: %v; entities: %s\n", lruCache, items)
	}
	using := func(funcs string) {
		fmt.Printf("=== using %s\n", funcs)
	}

	using("Put() and Get()")

	lruCache.Put(1, "value1")
	lruCache.Put(2, "value2")
	lruCache.Put(3, "value3")
	lruCache.Put(4, "value4")
	lruCache.Put(5, "value5")
	showLruCache()

	key := 1
	ok, val := lruCache.Get(key)
	fmt.Printf("%d => %s, exists: %t\n", key, val, ok)

	lruCache.Put(6, "value6")
	ok, val = lruCache.Get(key)
	fmt.Printf("%d => %s, exists: %t\n", key, val, ok)
	showLruCache()
	fmt.Println("Replaced the key 6 value:")
	lruCache.Put(6, "other6")
	showLruCache()

	using("PutIfAbsent()")
	key = 3
	ok, val = lruCache.PutIfAbsent(key, "other3")
	fmt.Printf("key: %d; old: %s; replaced: %t\n", key, val, ok)
	key = 7
	ok, val = lruCache.PutIfAbsent(key, "value7")
	fmt.Printf("key: %d; value: %s; added: %t\n", key, val, ok)
	showLruCache()

	using("Evict()")
	key = 4
	ok, val = lruCache.Get(key)
	fmt.Printf("%d => %s, exists: %t\n", key, val, ok)
	ok, val = lruCache.Evict(key)
	fmt.Printf("%d => %s, evicted: %t\n", key, val, ok)
	showLruCache()

	using("Clear()")
	lruCache.Clear()
	showLruCache()
}
