package multicache

/**
This file is part of multicache, a library for handling caches with multiple
keys and replacement algorithms.

Copyright 2015 Joseph Lewis <joseph@josephlewis.net>
Licensed under the MIT license
**/

/** ReplacementAlgorithm is the interface for all replacement algorithms to
follow.

When implementing a ReplacementAlgorithm DO NOT call any functions of the
Multicache or any external functions referencing the multicache. The Multicache
uses a multiple reader single writer lock and may deadlock if you do that.

MulticacheItem provides the field Tag for your use in implementing the
ReplacementAlgorithm. For example, you could use Tag to store an increasing
number if you were creating an LRU style replacement.
**/
type ReplacementAlgorithm interface {
	// Resets the items in the multicache
	Reset(multicache *Multicache)
	// Gets the next item to replace
	GetNextReplacement(multicache *Multicache) *MulticacheItem
	// True if the replacement algorithm makes modifications to itself
	// or the multicacheitem in the ItemRetrieved function.
	// this will cause a singular lock to happen rather than a readlock to
	// ensure safety of the datastructures.
	UpdatesOnRetrieved() bool
	// ItemRetrieved is called when an item is going to be returned by Get()
	// this can help update your algorithm like updating the time in a LRU.
	// It returns true if the item can be returned or false if nil should be
	// returned to the caller instead, for example in a time based cache.
	ItemRetrieved(item *MulticacheItem) bool
}
