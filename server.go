package mnemo

import (
	"fmt"
	"log"
	"sync"

	"github.com/labstack/echo/v4"
)

var srvMgr = &srvManager{
	servers: make(map[string]*Server),
}

type (
	srvManager struct {
		mu      sync.Mutex
		servers map[string]*Server
	}
	Server struct {
		app             *echo.Echo
		cfg             serverConfig
		msgs            chan []byte
		onNewConnection func(c *Conn)
		ConnectionPool  *Pool
	}
	serverConfig struct {
		Port string
		Path string
		Key  string
	}
)

func NewServer(cfg serverConfig) (*Server, error) {
	if cfg.Path == "" {
		return nil, fmt.Errorf("server path cannot root")
	}
	srvMgr.mu.Lock()
	for port, sv := range srvMgr.servers {
		if port == cfg.Port {
			return nil, fmt.Errorf("server with port %s already exists", cfg.Port)
		}
		if sv.cfg.Path == cfg.Path {
			return nil, fmt.Errorf("server with path %s already exists", cfg.Path)
		}
		if sv.cfg.Key == cfg.Key {
			return nil, fmt.Errorf("server with key %s already exists", cfg.Key)
		}
	}

	server := &Server{
		app:            echo.New(),
		cfg:            cfg,
		msgs:           make(chan []byte, 16),
		ConnectionPool: NewPool(),
	}

	srvMgr.servers[server.cfg.Port] = server
	srvMgr.mu.Unlock()

	ws := server.app.Group(cfg.Path + "/ws")
	ws.GET("/subscribe", server.handleSubscribe)

	return server, nil
}

func NewServerConfig(port, path, key string) serverConfig {
	//TODO: Use functional config
	return serverConfig{
		port, path, key,
	}
}

func (s *Server) Config() serverConfig {
	return s.cfg
}

func (s *Server) listenAndServe() {
	go s.app.Start(s.cfg.Port)
}

func (s *Server) Shutdown() error {
	srvMgr.mu.Lock()
	delete(srvMgr.servers, s.cfg.Port)
	srvMgr.mu.Unlock()
	err := s.app.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) handleSubscribe(ctx echo.Context) error {
	conn, err := NewConnection(ctx)
	if err != nil {
		return err
	}
	if err := s.ConnectionPool.AddConnection(conn); err != nil {
		return err
	}
	// Trigger user defined call back on new connection
	if s.onNewConnection != nil {
		go s.onNewConnection(conn)
	}
	conn.Listen()
	return conn.Close()
}

func (s *Server) SetOnNewConnection(callback func(c *Conn)) {
	s.onNewConnection = callback
}

func (s *Server) Publish(msg interface{}) {
	for _, conn := range s.ConnectionPool.Connections() {
		select {
		case conn.Messages <- msg:
		default:
			log.Println("closing connection: ", conn.Key)
			conn.Close()
		}
	}
}
