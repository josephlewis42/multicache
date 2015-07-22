package main

/**
This file is part of multicache, a library for handling caches with multiple
keys and replacement algorithms.

Copyright 2015 Joseph Lewis <joseph@josephlewis.net>
Licensed under the MIT license
**/

import (
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"strconv"
	"testing"

	"github.com/josephlewis42/multicache"

	"github.com/dkumor/golang-lru"
)

const (
	// The number of elements that the cache can hold, choose a number that
	// represents the percentage of unique elements you can fit in your cache
	// compared to the UniqueElements variable.
	CacheSize = 3

	// The number of items that we're going to test over
	TestSize = 10000

	// Number of elements in the entire "database" that we're caching
	UniqueElements = 15

	// Probability that an element is repeated from the last few items
	// this emulates "real" access patterns, e.g. items that were recently
	// accessed will be accessed again with this probability
	RepeatChance = .8

	// Max number of threads to try under
	MaxThreads = 8
)

func main() {
	// The main test data
	data := []string{}

	for i := 0; i < TestSize; i++ {
		if i > CacheSize && rand.Float64() < RepeatChance {
			// Give a chance to recently used values, likely to happen in the
			// real world
			index := (rand.Int() % (CacheSize - 1)) + 1
			chosen := data[i-index]
			data = append(data, chosen)

		} else {
			// rnadomly change to 1 of 15 items
			randval := rand.Int() % UniqueElements
			data = append(data, strconv.Itoa(randval))
		}
	}

	// Set up the algorithms we're going to test
	algs := map[string]multicache.ReplacementAlgorithm{"Round Robin": &multicache.RoundRobin{},
		"LRU":           &multicache.LeastRecentlyUsed{},
		"Random":        &multicache.Random{},
		"Second Chance": &multicache.SecondChance{},
		"Timed Cache":   multicache.CreateTimeExpireAlgorithm(1000)}

	algnames := getSortedKeys(algs)

	for thread := 1; thread <= MaxThreads; thread++ {
		fmt.Printf("%d Threads\n", thread)
		runtime.GOMAXPROCS(thread)
		for _, algname := range algnames {
			alg := algs[algname]

			wrapped := wrapForParallelBenchmarking(data, alg, CacheSize, thread)

			result := testing.Benchmark(wrapped)
			showResults(algname, result)
		}

		// Test golang-lru
		wrapped := golangLruParallelBenchmarking(data, CacheSize)
		result := testing.Benchmark(wrapped)
		showResults("golang-lru", result)

		// Now test with a cache that holds all items
		wrapped = wrapForBenchmarking(data, &multicache.Random{}, UniqueElements)
		result = testing.Benchmark(wrapped)
		showResults("NONE", result)

		fmt.Println()
	}
}

func showResults(algname string, result testing.BenchmarkResult) {
	fmt.Printf("%-15s Allocated: %5d kB Operations/Second: %7.0f\n",
		algname,
		result.MemBytes/1024,
		float64(result.N)/result.T.Seconds())

}

func wrapForBenchmarking(data []string, alg multicache.ReplacementAlgorithm, CacheSize uint64) func(b *testing.B) {
	return func(b *testing.B) {
		cache, _ := multicache.NewMulticache(CacheSize, alg)
		datLen := len(data)
		var ok bool
		// run the Fib function b.N times
		for n := 0; n < b.N; n++ {
			dat := data[n%datLen]
			_, ok = cache.Get(dat)
			if !ok {
				cache.Add(dat, dat)
			}
		}
	}
}

func wrapForParallelBenchmarking(data []string, alg multicache.ReplacementAlgorithm, CacheSize uint64, Threads int) func(b *testing.B) {
	return func(b *testing.B) {
		cache, _ := multicache.NewMulticache(CacheSize, alg)
		datLen := len(data)

		b.RunParallel(func(pb *testing.PB) {
			// Each goroutine has its own bytes.Buffer.
			var ok bool
			n := 0
			for pb.Next() {
				dat := data[n%datLen]
				_, ok = cache.Get(dat)
				if !ok {
					cache.Add(dat, dat)
				}
				n++
			}
		})
	}
}

func golangLruParallelBenchmarking(data []string, CacheSize uint64) func(b *testing.B) {
	return func(b *testing.B) {
		cache, _ := lru.New(int(CacheSize))
		datLen := len(data)

		b.RunParallel(func(pb *testing.PB) {
			// Each goroutine has its own bytes.Buffer.
			var ok bool
			n := 0
			for pb.Next() {
				dat := data[n%datLen]
				_, ok = cache.Get(dat)
				if !ok {
					cache.Add(dat, dat)
				}
				n++
			}
		})
	}
}

func getSortedKeys(m map[string]multicache.ReplacementAlgorithm) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}
