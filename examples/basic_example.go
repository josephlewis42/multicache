package main

import (
	"fmt"

	"github.com/josephlewis42/multicache"
)

func main() {
	// creates a cache that can hold 10 values
	cache, _ := multicache.NewDefaultMulticache(10)

	// Nothing yet in cache
	value, ok := cache.Get("foo")
	fmt.Printf("%v %v\n", value, ok) // <nil> false

	// Add a key value pair
	cache.Add("foo", 42)
	value, ok = cache.Get("foo")
	fmt.Printf("%v %v\n", value, ok) // 42 true

	// Add a multiple key-value pair, note that the value goes first.
	cache.AddMany(44.009, "bar", "baz")
	value, ok = cache.Get("bar")
	fmt.Printf("%v %v\n", value, ok) // 44.009 true

	// Delete one key and all of the "multikeys" get removed
	cache.Remove("baz")
	value, ok = cache.Get("bar")
	fmt.Printf("%v %v\n", value, ok) // <nil> false

	// Create a multicache with your desired replacement algorithm
	// (you can even make your own)
	multicache.NewMulticache(10, &multicache.SecondChance{})

	// We even have time expiring caches, items expire after the given number
	// of ms. (10 items, 1000ms)
	multicache.CreateTimeExpireMulticache(10, 1000)
}
