package multicache

/**
This file is part of go-multicache, a library for handling caches with multiple
keys and replacement algorithms.

Copyright 2015 Joseph Lewis <joseph@josephlewis.net>
Licensed under the MIT license
**/

/**
The SecondChance cache operates on a round-robin principle but varies in the
following manner:

When an item is accessed, a bit is set on the item. When another item is
inserted, the algorithm looks at the next item in the round-robin order. If the
accessed bit is set then the item is passed over and the accessed bit is
removed giving the item a "second chance". This algorithm has all the efficiency
of round robin but has some properties of LRU without needing O(N) time per
access.
**/
type SecondChance struct {
	position uint64
}

func (rof *SecondChance) Reset(multicache *MultiCache) {
	rof.position = 0
}

func (rof *SecondChance) GetNextReplacement(multicache *MultiCache) *MultiCacheItem {
	for {
		// Increment our position, falling over if necessary
		rof.position = uint64(rof.position+1) % multicache.cacheSize

		currentItem := multicache.itemList[rof.position]

		// This item hasn't been referenced since the last sweep
		if currentItem.Tag == 0 {
			return currentItem
		}

		// Mark this for being swept next time.
		currentItem.Tag = 0
	}
}

func (rof *SecondChance) UpdatesOnRetrieved() bool {
	return true
}

func (rof *SecondChance) ItemRetrieved(item *MultiCacheItem) {
	// Set a flag representing the item being referenced.
	item.Tag = 1
}
