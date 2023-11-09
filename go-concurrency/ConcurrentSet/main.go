// Package main contains an example of using the ConcurrentSet
package main

import (
	"fmt"
	"github.com/PavloVM7/go-concurrency/collections"
	"strconv"
	"sync"
	"time"
)

const (
	count   = 1_000_000
	threads = 10
	prefix  = "string number "
)

func main() {
	println("Example for using the ConcurrentSet.")
	stopString := createString(count)
	fmt.Printf("Thread number: %d, count: %d, stop string: '%s'\n", threads, count, stopString)

	chStart := make(chan struct{})
	chEnd := make(chan struct{})
	var wg sync.WaitGroup
	adds := make([]int, threads)
	set := collections.NewConcurrentSetCapacity[string](count)

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(n int) {
			work(n, set, adds, chStart, chEnd)
			wg.Done()
		}(i)
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	go tick(ticker, set, stopString, chEnd)

	start := time.Now()
	close(chStart)
	wg.Wait()
	dur := time.Since(start)

	sum := 0
	for i := 0; i < len(adds); i++ {
		sum += adds[i]
	}
	if sum != count {
		panic(fmt.Sprintf("The correct number of strings were not added to the set. Need %d strings, got %d",
			count, sum))
	}
	fmt.Printf("%d threads added %d strings, duration: %v\n", threads, sum, dur)
	fmt.Printf("every thread added: %v\n", adds)
}

func work(num int, set *collections.ConcurrentSet[string], adds []int, startCh, endCh <-chan struct{}) {
	<-startCh
	for i := 1; i <= count; i++ {
		str := createString(i)
		if set.Add(str) {
			adds[num]++
		}
	}
	<-endCh
}
func tick(ticker *time.Ticker, set *collections.ConcurrentSet[string], stopString string, chEnd chan struct{}) {
	var once sync.Once
	for range ticker.C {
		if set.Contains(stopString) {
			once.Do(func() {
				close(chEnd)
			})
		}
	}
}
func createString(num int) string {
	return prefix + strconv.Itoa(num)
}
