package store

import (
	"fmt"
	"reflect"
	"testing"
)

func Example() {
	fmt.Println("this is an example")
}

func TestNewStore(t *testing.T) {
	store, _ := NewStore("test")
	storeType := reflect.TypeOf(store)
	expected := "*store.Store"
	actual := storeType.String()
	if actual != expected {
		t.Errorf("expected NewStore to be of type %v; got %v", expected, actual)
	}

	if store.data == nil {
		t.Errorf("expected data map to be instantiated")
	}
}

func TestCreateStoreCache(t *testing.T) {
	store, _ := NewStore("test")
	key := "test"
	_, err := NewCache[int](store, StoreKey(key))
	if err != nil {
		t.Error(err)
	}
	_, err = NewCache[int](store, StoreKey(key))
	if err == nil {
		t.Error("expected duplicate key error")
	}
}

// TODO: Update with t.Run flags
func TestUseStoreCache(t *testing.T) {
	//TODO: Need TestMain to set up base store from init in store.go
	store, _ := NewStore("test")
	key := "test"
	cache, err := NewCache[int](store, StoreKey(key))
	if err != nil {
		t.Error(err)
	}

	data := 123
	cache.Cache(&data, key)
	cache, err = UseCache[int]("test", StoreKey(key))
	if err != nil {
		t.Error(err)
	}
	if len(cache.GetAll()) != 1 {
		t.Error("cache not retrieved")
	}

	_, err = UseCache[int]("test", StoreKey("test2"))
	if err == nil {
		t.Error("expected no data with key error")
	}

	_, err = UseCache[string]("test", StoreKey("test"))
	if err == nil {
		t.Error("expected invalid type for cache with key error")
	}
}
