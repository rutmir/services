package memcache

import (
	"fmt"
	"os"
	"testing"
)

var (
	mct         MemCache
	err         error
	testPackage *Item
)

func TestMain(m *testing.M) {
	url := os.Getenv("MEMCACHE_URL")
	mct, err = GetInstance("memcached", "test", url)
	testPackage = &Item{
		Key:        "Test",
		Value:      []byte("HELLO"),
		Expiration: 50}
	m.Run()
}

func TestInitial(t *testing.T) {
	if err != nil {
		t.Fatalf("failed to Set MemCache: %v", err)
	}
}

func TestSet(t *testing.T) {
	if err != nil {
		t.Fatalf("failed to Set MemCache: %v", err)
	}
	if err := mct.Set(testPackage); err != nil {
		t.Errorf("Error in Set function: %v", err)
	}
}

func TestGet(t *testing.T) {
	if err != nil {
		t.Fatalf("failed to Set MemCache: %v", err)
	}
	item, err := mct.Get(testPackage.Key)
	if err != nil {
		t.Errorf("Error in Set function: %v", err)
	} else {
		fmt.Println(item)
	}
}
