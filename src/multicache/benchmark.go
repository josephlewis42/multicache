package multicache

/**
This file is part of go-multicache, a library for handling caches with multiple
keys and replacement algorithms.

Copyright 2015 Joseph Lewis <joseph@josephlewis.net>
Licensed under the MIT license
**/

type HitMissTester struct {
	hits   int
	misses int
	items  []string
	mc     MultiCache
}

/** Given a list of strings, calculates the hit/miss ratio for a given algorithm
and cache size.

Returns the hit/miss ratio along with Bélády's optimal ratio for the given input
with a cache of the same size.

To minimize the warmup penalty, make sure the number of items is an order of
magnitude greater than the cache size.
**/
func CalculateHitMiss(items []string, cacheSize uint64, algorithm ReplacementAlgorithm) (ratio float64, perfectRatio float64) {
	if cacheSize <= 0 {
		return 0, 0
	}

	mc := NewMultiCache(cacheSize, algorithm)
	hits := 0
	misses := 0

	for _, item := range items {
		_, ok := mc.Get(item)
		if ok {
			hits++
		} else {
			misses++
			mc.Add(item, item)
		}

	}

	algorithmHitMissRatio := float64(hits) / float64(misses+hits)
	optimal := CalculateOptimalHitMiss(items, cacheSize)

	return algorithmHitMissRatio, optimal
}

/**
	Returns Bélády's optimal ratio for the given input with a cache of the same size.
**/
func CalculateOptimalHitMiss(items []string, cacheSize uint64) (ratio float64) {
	//log.Printf("doing %v size %v\n", items, cacheSize)
	// If we fit everything in cache, ratio = 0 because we have to insert
	// everything.
	if cacheSize <= 0 {
		return 0.0
	}

	// The "cache" holds the item key and a boolean of whether or not the key
	// is in the cache.
	cache := make(map[string]bool)

	// Initially nothing is in the cache
	for _, item := range items {
		cache[item] = false
	}

	hits := 0
	misses := 0

	// Insert our original number of cache items which are all misses
	for _, key := range items {
		cache[key] = false
	}

	// while the cache is not full yet and we haven't run out of items
	// add the items to the cache to get to capacity if they aren't there
	// already
	index := 0
	for ; uint64(misses) < cacheSize && index < len(items); index++ {
		key := items[index]
		res, _ := cache[key]
		if res == true {
			hits++
			continue
		}

		cache[key] = true
		misses++
	}
	//log.Printf("pre iteration hits: %v, misses: %v\n", hits, misses)

	for ; index < len(items); index++ {
		nextToAdd := items[index]
		nextToRemove := ""
		val := cache[nextToAdd]
		if val == true {
			// found
			hits++
		} else {
			// not found, remove the furthest item and return the closest
			nextToRemove = furthestFunc(index, &cache, &items)
			if nextToRemove != nextToAdd && nextToRemove != "" {
				cache[nextToAdd] = true
				cache[nextToRemove] = false
			}
			misses++
		}

		//log.Printf("val: %v, wasHit: %v, nextToRemove: %v cache: %v\n", nextToAdd, val, nextToRemove, cache)
	}

	//log.Printf("post iteration hits: %v, misses: %v\n", hits, misses)

	// avoid division errors
	if misses == 0 {
		return 1
	}

	// calculate the result
	algorithmHitMissRatio := float64(hits) / float64(misses+hits)
	return algorithmHitMissRatio
}

// This function will find the cache items that are going to be replaced farthest
// in the future
func furthestFunc(tmpIndex int, cache *map[string]bool, items *[]string) string {

	// Holds all the items in the cache
	cacheCopy := make(map[string]bool)
	for key, held := range *cache {
		if held == true {
			cacheCopy[key] = true
		}
	}

	// scan the list of items, removing each until we reach the end of the
	// list, or we get the item used farthest in the future.
	for ; tmpIndex < len(*items) && len(cacheCopy) > 1; tmpIndex++ {
		usedItem := (*items)[tmpIndex]
		delete(cacheCopy, usedItem)
	}

	furthestKey := ""

	// Return an item in the map
	for k, _ := range cacheCopy {
		furthestKey = k
		break
	}

	return furthestKey
}
