package mnemo

import (
	"fmt"
	"sync"
)

var strMgr = storeManager{
	stores: make(map[StoreKey]*Store),
}

type (
	storeManager struct {
		mu     sync.Mutex
		stores map[StoreKey]*Store
	}
	// Store is a collection of caches.
	Store struct {
		mu sync.Mutex
		*Server
		key      StoreKey
		data     map[CacheKey]any
		Commands Commands
	}
	// StoreKey is a unique identifier for a store.
	StoreKey string
)

// get returns the data for a given key or an error if the key does not exist.
func (s *Store) getCache(key CacheKey) (any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, ok := s.data[key]
	if !ok {
		return nil, NewError[Store](fmt.Sprintf("no cache with key '%v'", key))
	}
	return data, nil
}

// setCache sets the data for a given key.
func (s *Store) setCache(key CacheKey, data any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = data
}

// NewStore creates a new store or returns an error if a store with the same key already exists.
func NewStore(key StoreKey, opts ...Opt[Store]) (*Store, error) {
	strMgr.mu.Lock()
	defer strMgr.mu.Unlock()
	s := &Store{
		key:      key,
		data:     make(map[CacheKey]any),
		Commands: NewCommands(),
	}

	for _, o := range opts {
		o(s)
	}
	if _, ok := strMgr.stores[key]; ok {
		return nil, NewError[Store](fmt.Sprintf("store with key '%v' already exists", key))
	}
	strMgr.stores[key] = s
	return s, nil
}

// Key returns the store's key.
func (s *Store) Key() StoreKey {
	return s.key
}

// UseStore returns a store by key or an error if the store does not exist.
func UseStore(key StoreKey) (*Store, error) {
	strMgr.mu.Lock()
	defer strMgr.mu.Unlock()
	s, ok := strMgr.stores[key]
	if !ok {
		return nil, NewError[Store](fmt.Sprintf("no store with key '%v'", key))

	}
	return s, nil
}

// NewCache creates a new cache or returns an error if a cache with the same key already exists.
func NewCache[Cache any](s StoreKey, c CacheKey) (*cache[Cache], error) {
	store, err := UseStore(s)
	if err != nil {
		return nil, err
	}

	_, err = store.getCache(c)
	if err == nil {
		return nil, NewError[Cache](fmt.Sprintf("cache with key '%v' already exists", c))
	}

	nc := newCache[Cache]()
	store.setCache(c, nc)

	return nc, nil
}

// UseCache returns a cache by key or an error if the cache does not exist.
func UseCache[Cache any](s StoreKey, c CacheKey) (*cache[Cache], error) {
	store, err := UseStore(s)
	if err != nil {
		return nil, err
	}
	store.mu.Lock()
	defer store.mu.Unlock()

	// Check if cache exists
	data, err := store.getCache(c)
	if err != nil {
		return nil, err
	}

	// Check if cache is of type Cache
	cache, ok := data.(*cache[Cache])
	if !ok {
		return nil, NewError[Cache](
			fmt.Sprintf("cache with key '%v' is not of type '%v'", c, new(Cache)),
		)
	}
	return cache, nil
}
