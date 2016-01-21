package multicache

/**
This file is part of multicache, a library for handling caches with multiple
keys and replacement algorithms.

Copyright 2015 Joseph Lewis <joseph@josephlewis.net>
Licensed under the MIT license
**/

/**
Replaces the oldest item in the cache first, then the next oldest, etc.

This is an extremely simple (and thus fast) replacement strategy.
**/
type RoundRobin struct {
	position uint64
}

func (rof *RoundRobin) InitItem(item *MulticacheItem) {}

func (rof *RoundRobin) Reset(multicache *Multicache) {
	rof.position = 0
}

func (rof *RoundRobin) GetNextReplacement(multicache *Multicache) *MulticacheItem {
	rof.position = uint64(rof.position+1) % multicache.cacheSize
	return multicache.itemList[rof.position]
}

func (rof *RoundRobin) UpdatesOnRetrieved() bool {
	return false
}

func (rof *RoundRobin) ItemRetrieved(item *MulticacheItem) bool {
	// We don't update anything here, no need.
	return true
}
