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
	mc := NewDefaultMulticache(100)

	assert(t, mc.cacheSize == 100, "Cache size not set")
	assert(t, len(mc.kvStore) == 0, "Wrong kvstore size")
	assert(t, mc.replace != nil, "Algorithm missing")
}

func TestNewMulticacheAddGet(t *testing.T) {
	mc := NewDefaultMulticache(100)
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

// /** If GetOrFind misses the cache, this function is called. It should get the
// item for the given string and return it, the item's keys and optionally an error.
//
// If an error is returned, saving the item and keys is skipped and the error is
// passed on to the caller, otherwise the returned item is passed on and the error
// will be nil.
//
// searchKey is the key that we looked up that didn't exist.
// **/
// type GetOrFindMiss func(searchKey string) (item interface{}, keys []string, err error)
//
// /** GetOrFind checks to see if the given item is in the cache. If the item is
// in the cache, it returns the item and a nil error. If the item is not in the
// cache replaceFunc is called to get the requested item along with its keys; this
// item will be stored in the cache if err is nil. If err is not nil, GetOrFind
// will return a nil item and the error returned by GetOrFindMiss.
//
// **/
// func (mc *Multicache) GetOrFind(key string, replaceFunc GetOrFindMiss) (item interface{}, err error) {
// 	// Do a full write lock because we don't want a race condition in case we
// 	// need to write.
// 	c.lock.Lock()
// 	defer c.lock.Unlock()
//
// 	// Try to get the item, on success return it
// 	item, ok := mc.get(key)
// 	if ok {
// 		return item, nil
// 	}
//
// 	// Call replaceFunc to see if it can get the item instead.
// 	item, keys, err := replaceFunc(key)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// If replaceFunc was a success, add and return
// 	add(item, keys...)
// 	return item, nil
// }
//

func TestGetOrFind(t *testing.T) {
	mc := NewDefaultMulticache(100)

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
	mc := NewDefaultMulticache(100)

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

func TestNewMulticachePurge(t *testing.T) {
	mc := NewDefaultMulticache(100)

	// Add and get a single value
	mc.Add("key", "value")
	val, ok := mc.Get("key")
	assert(t, ok, "Didn't get inserted value")
	assert(t, val == "value", "Returned value was incorrect")
	mc.Purge()
	_, ok = mc.Get("key")
	assert(t, ok == false, "Got key after purge")
}
