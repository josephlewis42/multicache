package multicache

import "time"

/**
This file is part of multicache, a library for handling caches with multiple
keys and replacement algorithms.

Copyright 2015 Joseph Lewis <joseph@josephlewis.net>
Licensed under the MIT license
**/

/**

TimedExpire is a caching algorithm that gives each item an expiration based on
the time inserted. After that time the item will not be available.

If the cache is full (i.e. no items are expired) then the item with the least
amount of time remaining will be replaced.

**/
type TimedExpire struct {
	timeExpireMs int64
}

func (this *TimedExpire) InitItem(item *MulticacheItem) {
	item.Tag = time.Now().UnixNano() / int64(time.Millisecond)
}

func (this *TimedExpire) Reset(multicache *Multicache) {
}

func (this *TimedExpire) GetNextReplacement(multicache *Multicache) *MulticacheItem {

	currentTimeMs := time.Now().UnixNano() / int64(time.Millisecond)

	var difference int64 = 0
	var smallestItem *MulticacheItem

	for _, item := range multicache.itemList {
		// item.Tag has the creation time of this item
		timeDelta := currentTimeMs - item.Tag

		// short circuit the rest of the items.
		if timeDelta >= this.timeExpireMs {
			return item
		}

		// stored older item
		if timeDelta > difference {
			difference = timeDelta
			smallestItem = item
		}
	}

	return smallestItem
}

func (this *TimedExpire) UpdatesOnRetrieved() bool {
	return false
}

func (this *TimedExpire) ItemRetrieved(item *MulticacheItem) bool {
	// Make sure the item is still valid.

	currentTimeMs := time.Now().UnixNano() / int64(time.Millisecond)
	timeDelta := currentTimeMs - item.Tag

	isValid := timeDelta <= this.timeExpireMs

	return isValid
}

/**
Creates a new multicache that removes items inserted before expireTimeMs milliseconds
ago with numItems slots for items.
**/
func CreateTimeExpireMulticache(numItems uint64, expireTimeMs int64) (*Multicache, error) {
	return NewMulticache(numItems, CreateTimeExpireAlgorithm(expireTimeMs))
}

/**
Creates a timed expire algorithm.
**/
func CreateTimeExpireAlgorithm(expireTimeMs int64) *TimedExpire {
	return &TimedExpire{expireTimeMs}
}
