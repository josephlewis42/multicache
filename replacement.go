package multicache

/**
This file is part of go-multicache, a library for handling caches with multiple
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
	Reset(multicache *Multicache)
	GetNextReplacement(multicache *Multicache) *MulticacheItem
	UpdatesOnRetrieved() bool
	ItemRetrieved(item *MulticacheItem)
}
