package multicache

import "testing"

/**
This file is part of multicache, a library for handling caches with multiple
keys and replacement algorithms.

Copyright 2015 Joseph Lewis <joseph@josephlewis.net>
Licensed under the MIT license
**/

var lruTestCases = []ReplacementAlgorithmTestcase{
	// Miss all items because they aren't in cache
	{&LeastRecentlyUsed{}, 4, []string{"a", "b", "c", "d", "e"}, []bool{false, false, false, false, false}, 0},
	// Overwrite the first element and try it again
	{&LeastRecentlyUsed{}, 3, []string{"a", "b", "c", "d", "a"}, []bool{false, false, false, false, false}, 0},
	// Overwrite multiple in a row
	{&LeastRecentlyUsed{}, 2, []string{"a", "b", "c", "a", "b", "c"}, []bool{false, false, false, false, false, false}, 0},
	// Hit things
	{&LeastRecentlyUsed{}, 2, []string{"a", "b", "a", "b"}, []bool{false, false, true, true}, 0},
	// Out of order access
	{&LeastRecentlyUsed{}, 2, []string{"a", "b", "b", "a"}, []bool{false, false, true, true}, 0},
	// Make sure we didn't get rid of the thing we've used it most recently
	{&LeastRecentlyUsed{}, 2, []string{"a", "b", "a", "c", "a", "c"}, []bool{false, false, true, false, true, true}, 0},
	// A more complex example, we should trash b
	{&LeastRecentlyUsed{}, 3, []string{"a", "b", "c", "c", "a", "d", "b"}, []bool{false, false, false, true, true, false, false}, 0}}

func TestLRU(t *testing.T) {

	for _, testcase := range lruTestCases {
		testcase.RunTest(t)
	}
}
