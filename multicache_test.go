package multicache

import "testing"

/**
This file is part of go-multicache, a library for handling caches with multiple
keys and replacement algorithms.

Copyright 2015 Joseph Lewis <joseph@josephlewis.net>
Licensed under the MIT license
**/

func assert(t *testing.T, assertion bool, errinfo string) {
	if !assertion {
		t.Error(errinfo)
	}
}

func TestNewDefaultMulticache(t *testing.T) {
	mc, _ := NewDefaultMulticache(100)

	assert(t, mc.cacheSize == 100, "Cache size not set")
	assert(t, len(mc.kvStore) == 0, "Wrong kvstore size")
	assert(t, mc.replace != nil, "Algorithm missing")
}

func TestNewMulticacheAddGet(t *testing.T) {
	mc, _ := NewDefaultMulticache(100)
	_, ok := mc.Get("not there")
	assert(t, ok == false, "Got non existant item")

	// Add and get a single value
	mc.Add("key", "value")
	val, ok := mc.Get("key")
	assert(t, ok, "Didn't get inserted value")
	assert(t, val == "value", "Returned value was incorrect")

	// Try adding a second value
	mc.Add("key2", "value2")
	val, ok = mc.Get("key2")
	assert(t, ok, "Second item not ok")
	assert(t, val == "value2", "second item correct value")

	// Make sure everything else is the same
	val, ok = mc.Get("key")
	assert(t, ok, "Original Value Gone")
	assert(t, val == "value", "Original value changed")

	// Make sure we overwrite old values
	mc.Add("key", "value3")
	val, ok = mc.Get("key")
	assert(t, ok, "Didn't get replacement key")
	assert(t, val == "value3", "Didn't get correct value")
}

func TestGetOrFind(t *testing.T) {
	mc, _ := NewDefaultMulticache(100)

	// Test a missing key being added
	didAdd := false
	mc.GetOrFind("foo", func(key string) (item interface{}, keys []string, err error) {
		didAdd = true
		return "value", []string{"bar", "foo"}, nil
	})

	assert(t, didAdd, "Find wasn't called")
	val, ok := mc.Get("foo")
	assert(t, ok, "GetOrFind didn't add the key")
	assert(t, val == "value", "GetOrFind didn't add the proper value")

	val, ok = mc.Get("bar")
	assert(t, ok, "GetOrFind didn't add the second given key")
	assert(t, val == "value", "GetOrFind didn't add the proper value for the second key")

	// Test that an existing key won't be added

	didAdd = false
	mc.GetOrFind("foo", func(key string) (item interface{}, keys []string, err error) {
		didAdd = true
		return "value", []string{"bar", "foo"}, nil
	})

	assert(t, didAdd == false, "Find was called for existing key")

}

func TestNewMulticacheRemove(t *testing.T) {
	mc, _ := NewDefaultMulticache(100)

	mc.Add("key", "value")
	mc.AddMany("value2", "key2-1", "key2-2")

	mc.Remove("key")
	_, ok := mc.Get("key")
	assert(t, ok == false, "Got removed value")

	// Test multikey removal, should get rid of all
	// references
	mc.Remove("key2-1")
	_, ok = mc.Get("key2-2")
	assert(t, ok == false, "Didn't remove all multikey references")
}

func TestNewMulticacheRemoveItemFunc(t *testing.T) {
	mc, _ := NewDefaultMulticache(100)

	// Set up our items, we're going to remove the items that are equal
	// to a
	mc.Add("foo", "a")
	mc.Add("bar", "a")
	mc.Add("baz", "a")
	mc.Add("oof", "b")

	// Remove all items with value a
	mc.RemoveManyFunc(func(item interface{}) bool {
		if item.(string) == "a" {
			return true
		}

		return false
	})

	{
		_, ok := mc.Get("foo")
		assert(t, ok == false, "Got removed value foo")
	}

	{
		_, ok := mc.Get("bar")
		assert(t, ok == false, "Got removed value bar")
	}

	{
		_, ok := mc.Get("baz")
		assert(t, ok == false, "Got removed value baz")
	}

	{
		_, ok := mc.Get("oof")
		assert(t, ok == true, "Didn't get saved value oof")
	}
}

func TestNewMulticachePurge(t *testing.T) {
	mc, _ := NewDefaultMulticache(100)

	// Add and get a single value
	mc.Add("key", "value")
	val, ok := mc.Get("key")
	assert(t, ok, "Didn't get inserted value")
	assert(t, val == "value", "Returned value was incorrect")
	mc.Purge()
	_, ok = mc.Get("key")
	assert(t, ok == false, "Got key after purge")
}
