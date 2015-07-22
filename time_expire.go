package multicache

import (
	"math"
	"time"
)

/**
This file is part of multicache, a library for handling caches with multiple
keys and replacement algorithms.

Copyright 2015 Joseph Lewis <joseph@josephlewis.net>
Licensed under the MIT license
**/

const (
	nsToMs float64 = 1e-6
)

/**

TimedExpire is a caching algorithm that gives each item an expiration based on
the time inserted. After that time the item will not be available.

If the cache is full (i.e. no items are expired) then the item with the least
amount of time remaining will be replaced.

**/
type TimedExpire struct {
	timeExpireMs float64
}

func (this *TimedExpire) Reset(multicache *Multicache) {
}

func (this *TimedExpire) GetNextReplacement(multicache *Multicache) *MulticacheItem {

	currentTimeMs := float64(time.Now().Nanosecond()) * nsToMs

	difference := math.MaxFloat64
	var smallestItem *MulticacheItem

	for _, item := range multicache.itemList {
		itemTime := float64(item.Tag)
		timeDelta := currentTimeMs - itemTime

		// short circuit the rest of the items.
		if timeDelta <= this.timeExpireMs {
			return item
		}

		if timeDelta < difference {
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

	currentTimeMs := float64(time.Now().Nanosecond()) * nsToMs
	itemTime := float64(item.Tag)
	timeDelta := currentTimeMs - itemTime

	isValid := timeDelta >= this.timeExpireMs

	return isValid
}

/**
Creates a new multicache that removes items inserted before expireTimeMs milliseconds
ago with numItems slots for items.
**/
func CreateTimeExpireMulticache(numItems, expireTimeMs uint64) (*Multicache, error) {
	return NewMulticache(numItems, CreateTimeExpireAlgorithm(expireTimeMs))
}

/**
Creates a timed expire algorithm.
**/
func CreateTimeExpireAlgorithm(expireTimeMs uint64) *TimedExpire {
	return &TimedExpire{float64(expireTimeMs)}
}
