package router

import (
	"net/http"
)

type Registerer interface {
	RegisterRoutes(mux *http.ServeMux)
}

func New(handlers ...Registerer) *http.ServeMux {
	mux := http.NewServeMux()
	for _, h := range handlers {
		h.RegisterRoutes(mux)
	}
	return mux
}
