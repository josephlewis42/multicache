package multicache

import "testing"

/**
This file is part of go-multicache, a library for handling caches with multiple
keys and replacement algorithms.

Copyright 2015 Joseph Lewis <joseph@josephlewis.net>
Licensed under the MIT license
**/

type FurthestFuncTestcase struct {
	index    int
	cache    map[string]bool
	items    []string
	expected string
}

var furthestFuncTestcases = []FurthestFuncTestcase{
	{0, map[string]bool{"a": true, "b": true}, []string{"a", "a", "b"}, "b"},
	{1, map[string]bool{"a": true, "b": true}, []string{"a", "a", "b"}, "b"},
	// b not included, so we ignore it
	{0, map[string]bool{"a": true, "b": false}, []string{"a", "a", "b"}, "a"},
	// c never used
	{0, map[string]bool{"a": true, "b": true, "c": true}, []string{"a", "a", "b"}, "c"}}

func TestFurthestFunc(t *testing.T) {
	for index, test := range furthestFuncTestcases {
		result := furthestFunc(test.index, &test.cache, &test.items)

		if result != test.expected {
			t.Error("Unexpected result, index:", index, "result:", result, "testcase:", test)
		}
	}
}

type OptimalTestcase struct {
	items          []string
	cacheSize      uint64
	expectedResult float64
}

var optimalTestcases = []OptimalTestcase{
	// Nothing in cache, = 0
	{[]string{}, 0, 0.0},
	// cache all misses
	{[]string{"a", "b"}, 42, 0.0},
	// don't replace the b optimally, it'll never be used again
	{[]string{"a", "b", "a"}, 1, 0},
	// Hit twice
	{[]string{"a", "b", "a", "a"}, 2, 0.5},
	// randomly generated,
	{[]string{"3", // Miss, (3)
		"3",  // Hit, (3)
		"3",  // Hit, (3)
		"2",  // Miss, (3,2)
		"1",  // Miss, furthest: 3, (1, 2)
		"1",  // Hit, (1, 2)
		"2",  // Hit, (1, 2)
		"3",  // Miss, furthest: 2, (1,3)
		"1",  // Hit, (1, 3)
		"1",  // Hit, (1, 3)
		"3",  // Hit, (1, 3)
		"2",  // Miss, furthest: 3, (1, 3)
		"1",  // Miss, (1, 2)
		"3",  // Hit, (1, 3)
		"1"}, // Hit, (1, 3)
		2, 9.0 / 15.0}}

func TestOptimal(t *testing.T) {
	for index, test := range optimalTestcases {
		result := CalculateOptimalHitMiss(test.items, test.cacheSize)

		if result != test.expectedResult {
			t.Error("Unexpected result, index:", index, "result:", result, "testcase:", test)
		}
	}
}

type NormalTestcase struct {
	items          []string
	cacheSize      uint64
	expectedResult float64
}

//
// // This function will find the cache items that are going to be replaced farthest
// // in the future
// func furthestFunc(tmpIndex int, cache *map[string]bool, items *[]string) string {
//
// 	// Holds all the items in the cache
// 	cacheCopy := make(map[string]bool)
// 	for key, held := range *cache {
// 		if held == true {
// 			cacheCopy[key] = true
// 		}
// 	}
// 	// scan the list of items, removing each until we reach the end of the
// 	// list, or we get the item used farthest in the future.
// 	for ; tmpIndex < len(*items) || len(cacheCopy) == 1; tmpIndex++ {
// 		usedItem := (*items)[tmpIndex]
// 		delete(cacheCopy, usedItem)
// 	}
//
// 	// Return an item in the map
// 	for k, _ := range cacheCopy {
// 		return k
// 	}
//
// 	fatal("no items in cache")
// }
