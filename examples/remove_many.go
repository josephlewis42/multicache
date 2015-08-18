package main

import (
	"fmt"

	"github.com/josephlewis42/multicache"
)

func main() {
	// creates a cache that can hold 10 values
	cache, _ := multicache.NewDefaultMulticache(10)

	// session id and it's user
	cache.Add("1", "trillian")
	cache.Add("2", "zaphod")
	cache.Add("3", "trillian")
	cache.Add("4", "trillian")
	cache.Add("5", "ford_prefect")

	// Removes many items from the cache using a specified function that looks
	// at each item.
	// For example, if a user wants to log out of all sessions
	cache.RemoveManyFunc(func(item interface{}) bool {
		val, ok := item.(string)
		if !ok {
			return false
		}

		// true means remove the item, false means keep it.
		return val == "trillian"
	})

	// trillian is logged out
	val, ok := cache.Get("1")
	fmt.Printf("%v %v\n", val, ok) // <nil> false
	val, ok = cache.Get("3")
	fmt.Printf("%v %v\n", val, ok) // <nil> false

	// zaphod is still logged in
	val, ok = cache.Get("2")
	fmt.Printf("%v %v\n", val, ok) // zaphod true
}
