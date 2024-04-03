package http

import "net/http"

type router struct {
	*http.Server
}

func NewRouter(addr string, tableHandler *tableHandler, rowHandler *rowHandler) (*router, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tables", tableHandler.GetAllTables)
	mux.HandleFunc("GET /tables/{team}", tableHandler.GetTableByTeam)

	mux.HandleFunc("POST /error", rowHandler.InsertNewError)

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
