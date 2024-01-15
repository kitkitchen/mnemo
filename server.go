package mnemo

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var srvMgr = &srvManager{
	servers: make(map[int]*Server),
}

type (
	srvManager struct {
		mu sync.Mutex
		// Servers stored by port
		servers map[int]*Server
	}
	Server struct {
		*http.Server
		context.Context
		cfg             serverConfig
		msgs            chan []byte
		onNewConnection func(c *Conn)
		ConnectionPool  *Pool
	}
	ServerOpt    func(*Server)
	serverConfig struct {
		Port    int
		Pattern string
		Key     string
	}
)

func WithPort(port int) ServerOpt {
	return func(s *Server) {
		s.cfg.Port = port
		s.Server.Addr = ":" + strconv.Itoa(port)
	}
}

func WithPattern(pattern string) ServerOpt {
	return func(s *Server) {
		s.cfg.Pattern = pattern
	}
}

func NewServer(key string, opts ...ServerOpt) (*Server, error) {
	mux := http.NewServeMux()
	cfg := serverConfig{
		Key:  key,
		Port: srvMgr.AssignPort(),
	}
	srv := &Server{
		Server: &http.Server{
			Addr:           ":" + strconv.Itoa(cfg.Port),
			Handler:        mux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		cfg:            cfg,
		msgs:           make(chan []byte, 16),
		ConnectionPool: NewPool(),
	}

	for _, o := range opts {
		o(srv)
	}

	srvMgr.mu.Lock()
	defer srvMgr.mu.Unlock()

	for port, sv := range srvMgr.servers {
		if port == srv.cfg.Port {
			return nil, fmt.Errorf("server with port '%d' already exists", srv.cfg.Port)
		}
		if sv.cfg.Pattern == srv.cfg.Pattern {
			return nil, fmt.Errorf("server with pattern '%s' already exists", srv.cfg.Pattern)
		}
		if sv.cfg.Key == srv.cfg.Key {
			return nil, fmt.Errorf("server with key '%s' already exists", srv.cfg.Key)
		}
	}
	srvMgr.servers[srv.cfg.Port] = srv

	mux.HandleFunc(srv.cfg.Pattern+"/ws/subscribe", srv.HandleSubscribe)

	return srv, nil
}

func (s *Server) ListenAndServe() {
	go s.Server.ListenAndServe()
}

func (s *Server) Shutdown() error {
	srvMgr.mu.Lock()
	delete(srvMgr.servers, s.cfg.Port)
	srvMgr.mu.Unlock()
	err := s.Server.Shutdown(s.Context)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	conn, err := NewConnection(w, r)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := s.ConnectionPool.AddConnection(conn); err != nil {
		log.Fatal(err)
	}
	// Trigger user defined call back on new connection
	if s.onNewConnection != nil {
		go s.onNewConnection(conn)
	}
	conn.Listen()
}

func (s *Server) SetOnNewConnection(fn func(c *Conn)) {
	s.onNewConnection = fn
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

func (s *srvManager) AssignPort() int {
	next := 8080
	for port := range s.servers {
		if next < port {
			next = port + 1
		}
	}
	return next
}
