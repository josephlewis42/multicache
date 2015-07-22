package multicache

import "math/rand"

/**
This file is part of multicache, a library for handling caches with multiple
keys and replacement algorithms.

Copyright 2015 Joseph Lewis <joseph@josephlewis.net>
Licensed under the MIT license
**/

/**
Replaces a random item in the cache. This is implemented in some ARM processors
because it does fairly well despite being fast and simple.
**/
type Random struct {
}

func (rof *Random) Reset(multicache *Multicache) {
}

func (rof *Random) GetNextReplacement(multicache *Multicache) *MulticacheItem {
	location := uint64(rand.Uint32()) % multicache.cacheSize
	return multicache.itemList[location]
}

func (rof *Random) UpdatesOnRetrieved() bool {
	return false
}

func (rof *Random) ItemRetrieved(item *MulticacheItem) bool {
	// Push this item to the head of the queue
	return true
}
