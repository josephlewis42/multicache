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

// We're just testing that this function works, so we use a basic known
// algorithm rather than something more complex. The algorithms themselves are
// tested in their respective testing functions.

type FullTestcase struct {
	items           []string
	cacheSize       uint64
	algorithm       ReplacementAlgorithm
	expectedResult  float64
	expectedOptimal float64
}

// We use the same testcases as the optimal to make sure it comes out right too
var hitmissTestcases = []FullTestcase{
	// Nothing in cache, = 0
	{[]string{}, 0, &RoundRobin{}, 0.0, 0.0},
	// cache all misses
	{[]string{"a", "b"}, 42, &RoundRobin{}, 0.0, 0.0},
	// don't replace the b optimally, it'll never be used again
	{[]string{"a", "b", "a"}, 1, &RoundRobin{}, 0.0, 0},
	// Hit twice
	{[]string{"a", "b", "a", "a"}, 2, &RoundRobin{}, 0.5, 0.5},
	// randomly generated, these calculations are for round robin
	{[]string{"3", // Miss, (3)
		"3",  // Hit,  	(3)
		"3",  // Hit,  	(3)
		"2",  // Miss, 	(3,2)
		"1",  // Miss, 	(1, 2)
		"1",  // Hit,  	(1, 2)
		"2",  // Hit,  	(1, 2)
		"3",  // Miss, 	(3, 1)
		"1",  // Hit,  	(3, 1)
		"1",  // Hit, 	(3, 1)
		"3",  // Hit, 	(3, 1)
		"2",  // Miss,	(2, 3)
		"1",  // Miss, 	(1, 2)
		"3",  // Miss, 	(3, 1)
		"1"}, // Hit, 	(3, 1)
		2, &RoundRobin{}, 8.0 / 15.0, 9.0 / 15.0}}

func TestCalculateHitMiss(t *testing.T) {
	for index, test := range hitmissTestcases {
		result, optimal := CalculateHitMiss(test.items, test.cacheSize, test.algorithm)

		if result != test.expectedResult {
			t.Error("Unexpected result, index:", index, "result:", result, "testcase:", test)
		}

		if optimal != test.expectedOptimal {
			t.Error("Unexpected result, index:", index, "optimal_result:", optimal, "testcase:", test)
		}
	}
}
