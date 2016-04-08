package memcache

import (
	"fmt"
//	"time"
)

// Item is an item to be got or stored in a memcached server.
type Item struct {
	// Key is the Item's key (250 bytes maximum).
	Key        string

	// Value is the Item's value.
	Value      []byte

	// Flags are server-opaque flags whose semantics are entirely up to the app.
	Flags      uint32

	// Expiration is the cache expiration time, in seconds: either a relative
	// time from now (up to 1 month), or an absolute Unix epoch time.
	// Zero means the Item has no expiration time.
	Expiration int32

}

// ToString stringify Item object
func (f *Item) ToString() string {
	//return fmt.Sprintf("Item@ Key: %s, Value: %v, Expiration: %v", f.Key, f.Value, time.Unix(f.Expiration, 0).UTC())
	return fmt.Sprintf("Item@ Key: %s, Value: %v, Expiration: %v", f.Key, f.Value, f.Expiration)
}
