package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type routeDef struct {
	Method  string
	Path    string
	Name    string
	Handler http.HandlerFunc
}

func makeRouter(routes []routeDef) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(route.Handler)
	}

	return router
}
