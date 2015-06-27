package multicache

import "sync"

/**
This file is part of go-multicache, a library for handling caches with multiple
keys and replacement algorithms.

Copyright 2015 Joseph Lewis <joseph@josephlewis.net>
Licensed under the MIT license
**/

type MultiCache struct {
	kvStore         map[string]*MultiCacheItem
	itemList        []*MultiCacheItem
	cacheSize       uint64
	replace         ReplacementAlgorithm
	lock            sync.RWMutex
	retrieveUpdates bool
}

// Creates a new multicache that can hold the given number of items.
// The default algorithm used is SecondChance
func NewDefaultMultiCache(numItems uint64) *MultiCache {
	var defaultAlgorithm SecondChance
	return NewMultiCache(numItems, &defaultAlgorithm)
}

// Creates a multicache that can hold the given number of items using the given
// replacement algorithm. You should use CalculateHitMiss to look for the best
// ReplacementAlgorithm for your specific data.
func NewMultiCache(numItems uint64, algorithm ReplacementAlgorithm) *MultiCache {
	var mc MultiCache
	mc.kvStore = make(map[string]*MultiCacheItem)
	mc.itemList = make([]*MultiCacheItem, numItems)

	for i, _ := range mc.itemList {
		mc.itemList[i] = new(MultiCacheItem)
	}

	mc.cacheSize = numItems
	mc.replace = algorithm
	mc.retrieveUpdates = algorithm.UpdatesOnRetrieved()

	mc.Purge()

	return &mc
}

// Adds an item to the cache with the given key
func (mc *MultiCache) Add(key string, value interface{}) {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	mc.add(value, key)
}

/* Adds an item to the cache with the given keys

NOTE: do not include duplicate keys in AddMany e.g. AddMany("foo", "bar", "baz", "bar")
this will caused undefined results.
*/
func (mc *MultiCache) AddMany(value interface{}, keys ...string) {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	mc.add(value, keys...)
}

// Adds an item to the cache with the given keys
func (mc *MultiCache) add(value interface{}, keys ...string) {
	// Do nothing on empty key
	if len(keys) == 0 {
		return
	}

	cacheItem := mc.getItem()

	cacheItem.value = value
	cacheItem.keys = keys

	for _, key := range keys {
		// Remove old references if they exist.
		item, ok := mc.kvStore[key]
		if ok {
			mc.removeItem(item)
		}

		mc.kvStore[key] = cacheItem
	}
}

// Fetches an item from the cache
func (mc *MultiCache) Get(key string) (value interface{}, ok bool) {
	// If the caching algorithm updates some state when a get is done
	// do a normal lock, otherwise do a multiple reader lock for speed.
	if mc.retrieveUpdates {
		mc.lock.Lock()
		defer mc.lock.Unlock()
	} else {
		mc.lock.RLock()
		defer mc.lock.RUnlock()
	}

	return mc.get(key)
}

// This get function does no locking so it can be used elsewhere.
func (mc *MultiCache) get(key string) (value interface{}, ok bool) {
	v, ok := mc.kvStore[key]
	if !ok {
		return nil, false
	}

	mc.replace.ItemRetrieved(v)

	return v.value, true
}

/** If GetOrFind misses the cache, this function is called. It should get the
item for the given string and return it, the item's keys and optionally an error.

If an error is returned, saving the item and keys is skipped and the error is
passed on to the caller, otherwise the returned item is passed on and the error
will be nil.

searchKey is the key that we looked up that didn't exist.

NOTE: you must add searchKey to keys when returning without an error, if you do
not the resulting cache is undefined.

**/
type GetOrFindMiss func(searchKey string) (item interface{}, keys []string, err error)

/** GetOrFind checks to see if the given item is in the cache. If the item is
in the cache, it returns the item and a nil error. If the item is not in the
cache replaceFunc is called to get the requested item along with its keys; this
item will be stored in the cache if err is nil. If err is not nil, GetOrFind
will return a nil item and the error returned by GetOrFindMiss.

**/
func (mc *MultiCache) GetOrFind(key string, replaceFunc GetOrFindMiss) (item interface{}, err error) {
	// Do a full write lock because we don't want a race condition in case we
	// need to write.
	mc.lock.Lock()
	defer mc.lock.Unlock()

	// Try to get the item, on success return it
	item, ok := mc.get(key)
	if ok {
		return item, nil
	}

	// Call replaceFunc to see if it can get the item instead.
	item, keys, err := replaceFunc(key)
	if err != nil {
		return nil, err
	}

	// If replaceFunc was a success, add and return
	mc.add(item, keys...)
	return item, nil
}

// Removes an item from the multicache
func (mc *MultiCache) Remove(key string) {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	item, ok := mc.kvStore[key]
	if ok {
		mc.removeItem(item)
	}
}

// Removes all items from the cache.
func (mc *MultiCache) Purge() {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	mc.kvStore = make(map[string]*MultiCacheItem)

	for _, item := range mc.itemList {
		item.reset()
	}

	mc.replace.Reset(mc)
}

// Removes an item from the cache.
func (mc *MultiCache) removeItem(item *MultiCacheItem) {
	// Remove all references to this item.
	for _, v := range item.keys {
		delete(mc.kvStore, v)
	}

	item.softReset()
}

// Grabs and clears an item to be filled according to the replacement algorithm
func (mc *MultiCache) getItem() *MultiCacheItem {
	item := mc.replace.GetNextReplacement(mc)

	// Remove all references to this item.
	mc.removeItem(item)

	return item
}
