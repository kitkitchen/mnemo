// Package store blah blah blah
package store

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/labstack/gommon/log"
)

const ServerStore StoreKey = "servers"

var ServerCache = CacheKey("servers_cache")

var StoreManager = storeManager{
	stores: make(map[StoreKey]*Store),
}

func init() {
	// TODO: put this initializing and server starting as an option
	serverStore, _ := NewStore(ServerStore)
	cache, err := NewCache[Server](serverStore, ServerCache)
	if err != nil {
		log.Fatal(err)
	}
	cache.SetReducer(func(state Server) (mutation any) {
		return state.Config()
	})

	port := ":" + os.Getenv("STORE_PORT")
	if port == ":" {
		port = ":8080"
	}
	serverStore.Serve(port, "/store")
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

func UseStore(key StoreKey) *Store {
	StoreManager.mu.Lock()
	defer StoreManager.mu.Unlock()
	store, ok := StoreManager.stores[key]
	if !ok {
		return nil
	}
	return store
}

func NewStore(key StoreKey) (*Store, error) {
	StoreManager.mu.Lock()
	defer StoreManager.mu.Unlock()
	s := &Store{
		key:      key,
		data:     make(map[CacheKey]any),
		Commands: NewCommands(),
	}
	if _, ok := StoreManager.stores[key]; ok {
		return nil, fmt.Errorf("store with key '%v' already exists", key)
	}
	StoreManager.stores[key] = s
	return s, nil
}

func (s *Store) Serve(port string, path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return fmt.Errorf("empty store path for port: %s", port)
	}

	cfg := NewConfig(port, path, string(s.key))
	server, err := NewServer(cfg)
	if err != nil {
		return err
	}

	caches, err := UseCache[Server](ServerStore, ServerCache)
	if err != nil {
		log.Error(err)
	}
	if err = caches.Cache(server, s.key); err != nil {
		return err
	}
	s.listenAndServe()
	return nil
}

func (s *Store) Shutdown() {
	caches, err := UseCache[Server](ServerStore, ServerCache)
	if err != nil {
		log.Error(err)
	}
	item, ok := caches.Get(s.key)
	if !ok {
		log.Error("error retrieving store server with key: ", s.key)
	}
	item.Data.Shutdown()
}

func NewCache[Cache any](s *Store, key CacheKey) (*cache[Cache], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[key]
	if ok {
		return nil, fmt.Errorf("key '%v' already exists", key)
	}

	cs := newCache[Cache]()
	s.data[key] = cs

	return cs, nil
}

func UseCache[Cache any](s StoreKey, c CacheKey) (*cache[Cache], error) {
	store := UseStore(s)
	if store == nil {
		return nil, fmt.Errorf("no store with key '%v'", s)
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	data, ok := store.data[c]
	if !ok {
		return nil, fmt.Errorf("no cache with key '%v'", c)
	}
	cache, ok := data.(*cache[Cache])
	if !ok {
		return nil, fmt.Errorf("invalid type for cache with key %v", c)
	}
	return cache, nil
}
