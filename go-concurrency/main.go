package main

import (
	"fmt"
	"github.com/PavloVM7/go-concurrency/collections"
	"slices"
	"sync"
	"sync/atomic"
)

const (
	threads = 1_000
	number  = 10_000
)

var barrier int32

type counter struct {
	id    int
	count uint32
}

func (c *counter) String() string {
	return fmt.Sprintf("counter{id: %d, count: %d}", c.id, c.count)
}

func main() {
	println("Example of using a ConcurrencyMap.")
	fmt.Printf("Threads number: %d, number of counters: %d\n", threads, number)
	cache := collections.NewConcurrentMapCapacity[int, *counter](number)
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			loop(cache)
			wg.Done()
		}()
	}
	start()
	wg.Wait()
	keys := cache.Keys()
	slices.Sort(keys)
	for _, key := range keys {
		c, _ := cache.Get(key)
		if c.count != threads {
			panic(fmt.Sprintf("Attention! key: %d, %s, the counter not equals %d!!!", key, c, threads))
		}
	}
	println("All counters equal", threads)
}

func loop(cache *collections.ConcurrentMap[int, *counter]) {
	wait()
	for i := 0; i < number; i++ {
		_, c := cache.PutIfNotExistsDoubleCheck(i, &counter{id: i, count: 0})
		atomic.AddUint32(&c.count, 1)
	}
}

func wait() {
	for atomic.LoadInt32(&barrier) == 0 {
	}
}
func start() {
	atomic.StoreInt32(&barrier, 1)
}
