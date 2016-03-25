package memcache

import (
	"fmt"
	mc "github.com/bradfitz/gomemcache/memcache"
	"time"
)

// MemCache
type MemCache interface {
	Get(key string) (item *Item, err error)
	GetMulti(keys []string) (map[string]*Item, error)
	Touch(key string, seconds int32) (err error)
	Set(item *Item) error
	Delete(key string) error
}

type sourceType byte

type memcachedObject struct {
	Prefix string
	Source sourceType
	Client *mc.Client
}

const (
	MEM_CACHED sourceType = iota // 0
)

/* MemCached implementation */
// Get
func (f *memcachedObject) Get(key string) (*Item, error) {
	item, err := f.Client.Get(f.Prefix + key)
	if err != nil {
		return nil, err
	}

	return createItem(item, f.Prefix), nil
}

// GetMulti
func (f *memcachedObject) GetMulti(keys []string) (map[string]*Item, error) {
	newKeys := make([]string, len(keys))
	for _, item := range keys {
		newKeys = append(newKeys, f.Prefix + item)
	}
	items, err := f.Client.GetMulti(newKeys)
	if err != nil {
		return nil, err
	}
	result := make(map[string]*Item, len(items))

	for _, element := range items {
		result[element.Key[len(f.Prefix):]] = createItem(element, f.Prefix)
	}

	return result, nil
}

// Touch
func (f *memcachedObject) Touch(key string, seconds int32) error {
	return f.Client.Touch(key, seconds)
}

// Set
func (f *memcachedObject) Set(item *Item) error {
	return f.Client.Set(&mc.Item{Key: f.Prefix + item.Key, Value: item.Value, Expiration: item.Expiration})
}

// Delete
func (f *memcachedObject) Delete(key string) error {
	return f.Delete(f.Prefix + key)
}

// GetInstance
func GetInstance(cacheTarget, prefix string, v ...string) (MemCache, error) {
	switch cacheTarget {
	case "memcached":
		result := new(memcachedObject)
		result.Source = MEM_CACHED
		result.Prefix = prefix + "_"
		result.Client = mc.New(v...)
		result.Client.Timeout = time.Second * 5
		return result, nil
	}
	return nil, fmt.Errorf("%v not implemented", cacheTarget)
}

func createItem(item *mc.Item, pefix string) *Item {
	result := new(Item)
	result.Key = item.Key[len(pefix):]
	result.Value = item.Value
	result.Expiration = item.Expiration
	result.Flags = item.Flags

	return result
}
