package main

/**
This file is part of go-multicache, a library for handling caches with multiple
keys and replacement algorithms.

Copyright 2015 Joseph Lewis <joseph@josephlewis.net>
Licensed under the MIT license

This file generates a bunch of tests and compares all algorithms with each other
and the optimal caching algorithm.
**/

import (
	"fmt"
	"math/rand"
	"multicache"
	"sort"
	"strconv"
)

const (
	CacheSize      = 3
	TestSize       = 10000
	UniqueElements = 15
)

func main() {
	// Create some random data
	tests := map[string]*[]string{}

	// Build tests from 0 chance of repeat to 100
	for repeatChance := 0.0; repeatChance <= 1; repeatChance += .1 {
		data := []string{}

		for i := 0; i < TestSize; i++ {
			if i > CacheSize && rand.Float64() < repeatChance {
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

		testName := fmt.Sprintf("%3.0f%% Chance of Repeat", repeatChance*100)
		tests[testName] = &data
	}

	// Set up the algorithms we're going to test
	algs := map[string]multicache.ReplacementAlgorithm{"Round Robin": &multicache.RoundRobin{},
		"LRU":           &multicache.LeastRecentlyUsed{},
		"Random":        &multicache.Random{},
		"Second Chance": &multicache.SecondChance{}}

	algnames := getSortedKeys(algs)

	for testname, data := range tests {
		fmt.Println(testname)

		optimal := multicache.CalculateOptimalHitMiss(*data, CacheSize)
		fmt.Printf("Optimal Hit Percentage: %f\n\n", optimal)
		fmt.Printf("%-15s %-10s %-12s\n", "Algorithm", "Hit %", "% of optimal")

		for _, algname := range algnames {
			alg := algs[algname]

			amt := multicache.CalculateHitMiss(*data, CacheSize, alg)

			// Change everything to percentages
			amt *= 100
			pctOfOptimal := amt / optimal

			fmt.Printf("%-15s %3.5f   %3.5f\n", algname, amt, pctOfOptimal)
		}
		fmt.Println()
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
