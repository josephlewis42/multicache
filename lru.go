package multicache

/**
This file is part of go-multicache, a library for handling caches with multiple
keys and replacement algorithms.

Copyright 2015 Joseph Lewis <joseph@josephlewis.net>
Licensed under the MIT license
**/

/**
Replaces the least recently used element in the cache.

This is a common cache replacement scheme for smallish caches.

Note!

This takes O(n) time currently, similar to golang-lru. We don't specifically use a
min-heap because the items of MultiCache are all in an array and will be paged
together, using a min heap may actually cause cache thrashing whereas a
sequential scan should not and *should* be taken care of by the processor
automatically prefetching the next page before it is needed.
**/
type LeastRecentlyUsed struct {
	counter int32
}

func (rof *LeastRecentlyUsed) Reset(multicache *MultiCache) {
	rof.counter = 0
}

func (rof *LeastRecentlyUsed) GetNextReplacement(multicache *MultiCache) *MultiCacheItem {
	minItem := multicache.itemList[0]

	for _, item := range multicache.itemList {
		if item.Tag < minItem.Tag {
			minItem = item
		}
	}

	// Update the count on this item
	rof.ItemRetrieved(minItem)

	return minItem
}

func (rof *LeastRecentlyUsed) UpdatesOnRetrieved() bool {
	return true
}

func (rof *LeastRecentlyUsed) ItemRetrieved(item *MultiCacheItem) {
	rof.counter += 1
	item.Tag = rof.counter
}
