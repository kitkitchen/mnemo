package store

import (
	"fmt"
	"sync"
)

type Pool struct {
	mu    sync.Mutex
	conns map[interface{}]*Connection
}

func NewPool() *Pool {
	return &Pool{
		conns: make(map[interface{}]*Connection),
	}
}

func (p *Pool) Connections() map[interface{}]*Connection {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.conns
}

func (p *Pool) AddConnection(c *Connection) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, ok := p.conns[c.Key]
	if ok {
		return fmt.Errorf("connection with key %v already exists", c.Key)
	}
	p.conns[c.Key] = c
	c.Pool = p
	return nil
}

func (p *Pool) removeConnection(c *Connection) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.conns, c)
}
