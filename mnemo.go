// Package mnemo provides a robust way to generate and manage caches for any data type.
package mnemo

type (
	// Mnemo is the main struct for the Mnemo package.
	Mnemo struct {
		Server *Server
		logger Logger
	}
	Opt[T any] func(t *T)
)

// New returns a new Mnemo instance.
func New(opts ...Opt[Mnemo]) *Mnemo {
	m := &Mnemo{
		logger: logger,
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

// WithServer sets the server for the Mnemo instance from the NewServer function.
func WithServer(key string, opts ...Opt[Server]) Opt[Mnemo] {
	srv, err := NewServer(key, opts...)
	if err != nil {
		NewError[Server](err.Error()).Log()
	}
	return func(m *Mnemo) {
		m.Server = srv
	}
}
