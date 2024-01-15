// Package store blah blah blah
package mnemo

import (
	"fmt"
	"os"
	"sync"

	"github.com/labstack/gommon/log"
)

const ServerStore StoreKey = "servers"

var ServerCache = CacheKey("servers_cache")

var strMgr = storeManager{
	stores: make(map[StoreKey]*Store),
}

func init() {
	// TODO: put this initializing and server starting as an option
	serverStore, _ := NewStore(ServerStore, WithServer("/servers"))
	cache, err := NewCache[Server](ServerStore, ServerCache)
	if err != nil {
		log.Fatal(err)
	}
	cache.SetReducer(func(state Server) (mutation any) {
		return state.cfg
	})

	port := ":" + os.Getenv("MNEMO_PORT")
	if port == ":" {
		port = ":8080"
	}
	err = serverStore.Serve()
	if err != nil {
		log.Error(StoreError{
			Err: err,
		})
	}
}

type (
	storeManager struct {
		mu     sync.Mutex
		stores map[StoreKey]*Store
	}
	Store struct {
		mu sync.Mutex
		*Server
		key      StoreKey
		data     map[CacheKey]any
		Commands Commands
	}
	StoreKey string
)

func NewStore(key StoreKey, opts ...func(s *Store)) (*Store, error) {
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
		return nil, StoreError{
			Err: fmt.Errorf("store with key '%v' already exists", key),
		}
	}
	strMgr.stores[key] = s
	return s, nil
}

func WithServer(key string, opts ...ServerOpt) func(s *Store) {
	srv, err := NewServer(key, opts...)
	if err != nil {
		log.Fatal(StoreError{
			Err: err,
		})
	}
	return func(s *Store) {
		s.Server = srv
	}
}

func (s *Store) Serve() error {
	if s.Server == nil {
		return StoreError{
			Err: fmt.Errorf("store with key '%v' was not instantiated with a server", s.cfg.Key),
		}
	}
	caches, err := UseCache[Server](ServerStore, ServerCache)
	if err != nil {
		return err
	}
	if err = caches.Cache(s.Server, s.key); err != nil {
		return err
	}
	s.ListenAndServe()
	return nil
}

func (s *Store) Shutdown() {
	caches, err := UseCache[Server](ServerStore, ServerCache)
	if err != nil {
		log.Error(err)
	}
	item, ok := caches.Get(s.key)
	if !ok {
		err = StoreError{
			Err: fmt.Errorf("error retrieving store server with key: '%v'", s.key),
		}
		log.Error(err)
		return
	}
	item.Data.Shutdown()
}

func UseStore(key StoreKey) *Store {
	strMgr.mu.Lock()
	defer strMgr.mu.Unlock()
	s, ok := strMgr.stores[key]
	if !ok {
		log.Error(StoreError{
			Err: fmt.Errorf("no store with key: '%v", key),
		})
		return nil
	}
	return s
}

func NewCache[Cache any](s StoreKey, c CacheKey) (*cache[Cache], error) {
	store := UseStore(s)
	if store == nil {
		return nil, StoreError{
			Err: fmt.Errorf("no store with key '%v'", s),
		}
	}
	store.mu.Lock()
	defer store.mu.Unlock()

	_, ok := store.data[c]
	if ok {
		return nil, fmt.Errorf("key '%v' already exists", c)
	}

	cs := newCache[Cache]()
	store.data[c] = cs

	return cs, nil
}

func UseCache[Cache any](s StoreKey, c CacheKey) (*cache[Cache], error) {
	store := UseStore(s)
	if store == nil {
		return nil, StoreError{
			Err: fmt.Errorf("no store with key '%v'", s),
		}
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	data, ok := store.data[c]
	if !ok {
		return nil, StoreError{
			Err: fmt.Errorf("no cache with key '%v'", c),
		}
	}
	cache, ok := data.(*cache[Cache])
	if !ok {
		return nil, StoreError{
			Err: fmt.Errorf("invalid type for cache with key '%v'", c),
		}
	}
	return cache, nil
}
