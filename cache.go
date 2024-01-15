package mnemo

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

type (
	cache[T any] struct {
		mu        sync.Mutex
		createdAt time.Time
		raw       *raw[T]
		reducer   *reducer[T]
	}
	raw[T any] struct {
		caches  map[CacheKey]*Item[T]
		history map[time.Time]map[CacheKey]Item[T]
		feed    chan map[time.Time]map[CacheKey]Item[T]
	}
	reducer[T any] struct {
		reduce  *func([]reducerCache[T]) []reducerCache[any]
		history map[time.Time][]reducerCache[any]
		feed    chan reducerFeed[any]
	}
	Item[T any] struct {
		CreatedAt time.Time `json:"created_at"`
		Data      *T        `json:"data"`
	}
	cacheTimeoutConfig[T any] struct {
		data       *T
		key        any
		timeoutFun func(data *T)
		timeout    time.Duration
	}
	CacheKey                  any
	ReducerFunc[T any, U any] func(state T) (mutation U)
	ReducerConfig[T any]      struct {
		CreatedAt time.Time
		Data      []reducerCache[T]
	}
	reducerCache[T any] struct {
		Key       CacheKey  `json:"key"`
		CreatedAt time.Time `json:"created_at"`
		Data      T         `json:"data"`
	}
	reducerFeed[T any] struct {
		CreatedAt time.Time         `json:"created_at"`
		Cache     []reducerCache[T] `json:"cache"`
	}
	reducerHistory[T any] reducerFeed[T]
)

func newCache[T any]() (data *cache[T]) {
	c := &cache[T]{
		createdAt: time.Now(),
		raw: &raw[T]{
			caches:  make(map[CacheKey]*Item[T]),
			history: make(map[time.Time]map[CacheKey]Item[T]),
			feed:    make(chan map[time.Time]map[CacheKey]Item[T], 1024),
		},
		reducer: &reducer[T]{
			history: make(map[time.Time][]reducerCache[any]),
			feed:    make(chan reducerFeed[any], 1024),
		},
	}
	return c
}

func (c *cache[T]) monitorChanges(setup chan bool) {
	pRaw := c.copyRaw()
	setup <- true
	prev := c.reduce(pRaw)
	t := time.Now()
	c.cacheRaw(t, pRaw)
	c.cacheReduction(t, prev)
	for {
		raw := c.copyRaw()
		current := c.reduce(raw)
		if !reflect.DeepEqual(prev, current) {
			t := time.Now()
			c.cacheRaw(t, raw)
			c.cacheReduction(t, current)
			prev = current
		}
	}
}

func (c *cache[T]) copyRaw() map[CacheKey]Item[T] {
	c.mu.Lock()
	defer c.mu.Unlock()
	copy := make(map[CacheKey]Item[T])
	for k, v := range c.raw.caches {
		copy[k] = *v
	}
	return copy
}

func (c *cache[T]) cacheRaw(t time.Time, copy map[CacheKey]Item[T]) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.raw.history[t] = copy
	c.raw.feed <- c.raw.history
}

func newCacheReducer[T, U any](f ReducerFunc[T, U]) func([]reducerCache[T]) []reducerCache[U] {
	return func(rdt []reducerCache[T]) []reducerCache[U] {
		rdu := []reducerCache[U]{}
		for _, d := range rdt {
			rdu = append(rdu, reducerCache[U]{
				Key:       d.Key,
				CreatedAt: d.CreatedAt,
				Data:      f(d.Data),
			})
		}
		return rdu
	}
}
func (c *cache[T]) reduce(copy map[CacheKey]Item[T]) []reducerCache[any] {
	c.mu.Lock()
	defer c.mu.Unlock()

	data := []reducerCache[T]{}
	for key, item := range copy {
		data = append(data, reducerCache[T]{
			Key:       key,
			CreatedAt: item.CreatedAt,
			Data:      *item.Data,
		})
	}
	reduce := *c.reducer.reduce
	return reduce(data)
}

func (c *cache[T]) cacheReduction(t time.Time, r []reducerCache[any]) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.reducer.history[t] = r
	rf := reducerFeed[any]{}
	for _, k := range r {
		//TODO sort by createdAt
		rf.Cache = append(rf.Cache, k)
	}
	c.reducer.feed <- rf
}

func (c *cache[T]) SetReducer(rf ReducerFunc[T, any]) {
	cr := newCacheReducer(rf)
	c.mu.Lock()
	c.reducer.reduce = &cr
	c.mu.Unlock()

	setup := make(chan bool)
	defer close(setup)
	go c.monitorChanges(setup)
	<-setup
}

func (c *cache[T]) DefaultReducer(state T) (mutation any) {
	return state
}

func (c *cache[T]) RawFeed() chan map[time.Time]map[CacheKey]Item[T] {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.raw.feed
}

func (c *cache[T]) ReducerFeed() chan reducerFeed[any] {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.reducer.feed
}

func (c *cache[T]) RawHistory() map[time.Time]map[CacheKey]Item[T] {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.raw.history
}

func (c *cache[T]) ReducerHistory() []reducerHistory[any] {
	c.mu.Lock()
	defer c.mu.Unlock()
	rh := []reducerHistory[any]{}
	for time, cache := range c.reducer.history {
		rh = append(rh, reducerHistory[any]{
			CreatedAt: time,
			Cache:     cache,
		})
	}
	return rh
}

func (c *cache[T]) Get(key CacheKey) (Item[T], bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	data := c.raw.caches[key]
	if data == nil {
		return *new(Item[T]), false
	}
	return *data, true
}

// TODO: ensure this copy works
func (c *cache[T]) GetAll() map[CacheKey]Item[T] {
	c.mu.Lock()
	defer c.mu.Unlock()
	cache := make(map[CacheKey]Item[T])
	for key, item := range c.raw.caches {
		cache[key] = *item
	}
	return cache
}

func (c *cache[T]) Cache(data *T, key CacheKey) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.raw.caches[key] != nil {
		return fmt.Errorf("duplicate cache key: %v", key)
	}
	item := &Item[T]{
		Data:      data,
		CreatedAt: time.Now(),
	}
	c.raw.caches[key] = item
	return nil
}

func (c *cache[T]) CacheWithTimeout(cfg cacheTimeoutConfig[T]) error {
	c.Cache(cfg.data, cfg.key)
	if !(cfg.timeout > time.Second*0) {
		return fmt.Errorf(
			"cache not set for timeout: %v; timeout must be greater than 0", cfg.timeout,
		)
	}

	go func() {
		timer := time.NewTimer(cfg.timeout)
		<-timer.C
		cache, ok := c.Get(cfg.key)
		if !ok {
			logger.Fatalf("could not get cache with key %v", cfg.key)
		}
		err := c.Delete(cfg.key)
		if err != nil {
			logger.Fatal(err)
		}
		cfg.timeoutFun(cache.Data)
	}()
	return nil
}

func (c *cache[T]) Update(key CacheKey, update Item[T]) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.raw.caches[key]
	if !ok {
		return false
	}
	//TODO: ensure this is being updated in reducer
	c.raw.caches[key] = &update
	return true
}

func (c *cache[T]) Delete(key interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.raw.caches[key] == nil {
		return fmt.Errorf("no cache with key: %v", key)
	}
	delete(c.raw.caches, key)
	return nil
}

func NewCacheTimeoutConfig[T any](
	data *T,
	key interface{},
	timeoutFun func(data *T),
	timeout time.Duration,
) cacheTimeoutConfig[T] {
	return cacheTimeoutConfig[T]{
		data:       data,
		key:        key,
		timeoutFun: timeoutFun,
		timeout:    timeout,
	}
}
