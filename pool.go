package mnemo

import (
	"fmt"
	"sync"
)

type Pool struct {
	mu    sync.Mutex
	conns map[interface{}]*Conn
}

func NewPool() *Pool {
	return &Pool{
		conns: make(map[interface{}]*Conn),
	}
}

func (p *Pool) Connections() map[interface{}]*Conn {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.conns
}

func (p *Pool) AddConnection(c *Conn) error {
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

func (p *Pool) removeConnection(c *Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.conns, c)
}
