package multicache

import "testing"

/**
This file is part of go-multicache, a library for handling caches with multiple
keys and replacement algorithms.

Copyright 2015 Joseph Lewis <joseph@josephlewis.net>
Licensed under the MIT license
**/

type ReplacementAlgorithmTestcase struct {
	// The algorithm this test is used for
	ra ReplacementAlgorithm

	// The size of the cache for this test
	cacheSize uint64

	// The queries that are passed to Get(), if they are not found they are
	// then added.
	queries []string

	// The expected results from Get() for each item in queries
	findExpected []bool
}

func (rat *ReplacementAlgorithmTestcase) RunTest(t *testing.T) {
	mc := NewMulticache(rat.cacheSize, rat.ra)

	for index, query := range rat.queries {
		_, found := mc.Get(query)
		if !found {
			mc.Add(query, query)
		}

		if found != rat.findExpected[index] {
			t.Error("Unexpected result, index:", index, rat)

		}
	}

}
