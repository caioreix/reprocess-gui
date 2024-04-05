package http

import "net/http"

type Router struct {
	TableHandler    *tableHandler
	RowHandler      *rowHandler
	ConsumerHandler *consumerHandler
}

type router struct {
	*http.Server
}

func (r *Router) NewRouter(addr string) (*router, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tables", r.TableHandler.GetAllTables)
	mux.HandleFunc("GET /tables/{team}", r.TableHandler.GetTableByTeam)

	mux.HandleFunc("POST /error", r.RowHandler.InsertNewError)

	mux.HandleFunc("POST /consumer", r.ConsumerHandler.InsertNewConsumer)

	return &router{
		&http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}, nil
}

func (r *router) Serve() error {
	return r.ListenAndServe()
}
