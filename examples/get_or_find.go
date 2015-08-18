package main

import (
	"fmt"

	"github.com/josephlewis42/multicache"
)

func myMissFunction(searchKey string) (item interface{}, keys []string, err error) {
	// this may be a heavy DB operation or something; instead we'll just
	// concatenate a string. You have to return the searchKey in the
	// results, otherwise it won't get cached. It also won't get cached if
	// an error is returned.

	item = interface{}(searchKey + "'s value")

	return item, []string{searchKey}, nil
}

func main() {
	// creates a cache that can hold 10 values
	cache, _ := multicache.NewDefaultMulticache(10)

	// GetOrFind returns the item if it's in the cache or calls the given
	// function to "find" it, good if you use caching a lot and don't want to
	// keep using the if ! ok { fill_the_cache(); handle_error()}
	item, err := cache.GetOrFind("myKey", myMissFunction)

	fmt.Printf("%v, %v\n", item, err) // myKey's value, <nil>

}
