package http

import (
	"net/http"
	"time"
)

// Router provides a routing mechanism for handling HTTP requests.
type Router struct {
	TableHandler    *tableHandler
	RowHandler      *rowHandler
	ConsumerHandler *consumerHandler
}

type router struct {
	*http.Server
}

// NewRouter creates a new HTTP router with the specified address.
// It sets up handlers for various HTTP endpoints and returns a router.
func (r *Router) NewRouter(addr string, readHeaderTimeout time.Duration) (*router, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tables", r.TableHandler.GetAllTables)
	mux.HandleFunc("GET /tables/{team}", r.TableHandler.GetTableByTeam)

	mux.HandleFunc("POST /error", r.RowHandler.InsertNewError)

	mux.HandleFunc("GET /consumers", r.ConsumerHandler.GetPagedConsumers)
	mux.HandleFunc("POST /consumer", r.ConsumerHandler.InsertNewConsumer)

	return &router{
		&http.Server{
			ReadHeaderTimeout: readHeaderTimeout,
			Addr:              addr,
			Handler:           mux,
		},
	}, nil
}

// Serve starts serving HTTP requests.
func (r *router) Serve() error {
	return r.ListenAndServe()
}
