multicache
===========

[![GoDoc](http://godoc.org/github.com/josephlewis42/multicache?status.svg)](http://godoc.org/github.com/josephlewis42/multicache)

A speedy go caching library that supports multiple keys and various replacement algorithms.

Features
--------

* Support for caching items with multiple keys
* Lots of common out of the box replacement algorithms
	* LRU
	* Time Expiration
	* Round Robin
	* Random Replace
	* Second Chance
* Easily benchmark your application's access patterns to find the optimal configuration
* Custom replacement algorithms supported
* Very fast (see benchmarks below)
* Pull and feature requests are welcome!

What makes it fast?
-------------------

* String keys make for fast comparisons
* No dynamic memory allocation (except your keys and values)
* Cache elements are stored in an array so they are paged together and fit in an L2/L3 cache (unlike linked lists)
* Algorithms use basic comparisons and integer math
* Few (if any) external API calls depending on the caching algorithm

Examples
--------

Bare minimal example:

	package main

	import (
		"fmt"
		"github.com/josephlewis42/multicache"
	)

	func main() {
		// creates a cache that can hold 10 values
		cache, _ := multicache.NewDefaultMulticache(10)

		cache.Add("foo", 42)                // Add a key value pair
		cache.AddMany(44.009, "bar", "baz") // one value for many keys (value first)

		value, ok := cache.Get("foo")
		fmt.Printf("%v %v\n", value, ok) // 42 true

		value, ok = cache.Get("bar")
		fmt.Printf("%v %v\n", value, ok) // 44.009 true

		// Delete one key and all of the "multikeys" get removed
		cache.Remove("baz")
		value, ok = cache.Get("bar")
		fmt.Printf("%v %v\n", value, ok) // <nil> false
	}

Other cool features are described in the `examples` directory such as:

* custom removal of elements
* function callbacks to fill in data when there is a cache miss
* creating caches with other replacement algorithms
* creating a time expiring cache

Benchmarks
==========

Speed
-----

All benchmarks are BS; therefore you can modify the parameters in `examples/speedtest/speedtest.go` to fit your exact system.

Results (sorted for ops/sec):

		NONE            Allocated:     5 kB Operations/Second: 4164141
		Round Robin     Allocated: 24851 kB Operations/Second: 2895037
		LRU             Allocated: 23772 kB Operations/Second: 2758429
		Random          Allocated: 27642 kB Operations/Second: 2727461
		Second Chance   Allocated: 24335 kB Operations/Second: 2689030
		golang-lru      Allocated: 96809 kB Operations/Second: 1741016

		Timed Cache     Allocated: 62517 kB Operations/Second: 1081585

* *NONE* is a cache that hits every time.
* *golang-lru* is an external lru cache for comparison with this system.
* *Timed Cache* checks the time of each element coming out and nixes it if it is old; this causes many calls to `time()` which slows it considerably.


Results are for a cache that holds 20% of the data, averaged over 10,000
trials, with 80% chance that one of the items in the current cache is going to
be accessed again in the next `Get()` similar to what users would do if they were logged into the system.


Optimality
----------

You can also test the optimality of various algorithms for your workflow using
`examples/comparison.go`. This will give you an idea of how the caches will
perform "in the wild". By default, they'll run through various repeat chances
(likelihood that an item already in the cache will be the next item gotten).

	 70% Chance of Repeat
	Optimal Hit Percentage: 0.816000

	Algorithm       Hit %      % of optimal
	LRU             76.58000   93.84804
	Random          73.13000   89.62010
	Round Robin     75.82000   92.91667
	Second Chance   76.24000   93.43137

The optimality metric is based on Bélády's optimal ratio i.e. the optimal thing to dump from a cache is the one being used farthest from now.
