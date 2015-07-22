multicache
===========

A go caching library that supports multiple keys and various replacement algorithms.

Exapmles
========



Benchmarks
==========

Speed
-----

You can play around with the benchmark params in examples/speedtest/speedtest.go
to fit your exact system.

These results are for a cache that holds 20% of the dataset, averaged over 10000
trials, with 80% chance that one of the items in the current cache is going to
be accessed again in the next `Get()`--this simulates an access pattern where
you're working with an item over and over before switching to a different one,
like users logged into a system.

Results (sorted for ops/sec):

		NONE            Allocated:     5 kB Operations/Second: 4164141
		Round Robin     Allocated: 24851 kB Operations/Second: 2895037
		LRU             Allocated: 23772 kB Operations/Second: 2758429
		Random          Allocated: 27642 kB Operations/Second: 2727461
		Second Chance   Allocated: 24335 kB Operations/Second: 2689030
		golang-lru      Allocated: 96809 kB Operations/Second: 1741016

		Timed Cache     Allocated: 62517 kB Operations/Second: 1081585

* *NONE* is a cache that hits every time.
* *golang-lru* is an external lru cache for comparison with this system. We use
a static number of pre-allocated elements so the overhead of malloc()s is
greatly reduced.
* *Timed Cache* checks the time of each element coming out and nixes it if
it is old; this causes many calls to `time()` which slows it considerably.

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

The optimality metric is based on Bélády's optimal ratio.
