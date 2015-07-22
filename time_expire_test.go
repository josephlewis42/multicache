package multicache

import "testing"

/**
This file is part of multicache, a library for handling caches with multiple
keys and replacement algorithms.

Copyright 2015 Joseph Lewis <joseph@josephlewis.net>
Licensed under the MIT license
**/

var timedExpireTestCases = []ReplacementAlgorithmTestcase{
	// Miss all items because they instantly expire
	{&TimedExpire{0}, 2, []string{"a", "a", "a", "b", "b"}, []bool{false, false, false, false, false}, 1},
	// Hit all items because they don't expire
	{&TimedExpire{100}, 2, []string{"a", "a", "a", "b", "b"}, []bool{false, true, true, false, true}, 1},
	// Overwrite the first element and try it again
	{&TimedExpire{5}, 2, []string{"a", "b", "c", "a", "b"}, []bool{false, false, false, false, true}, 1}}

func TestTimedExpire(t *testing.T) {

	for _, testcase := range roundRobinTestCases {
		testcase.RunTest(t)
	}
}
