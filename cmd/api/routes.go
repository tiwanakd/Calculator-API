package main

import "net/http"

func (a *api) routes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /add", a.add)
	router.HandleFunc("POST /subtract", a.subtract)
	router.HandleFunc("POST /multiply", a.multiply)
	router.HandleFunc("POST /divide", a.divide)
	router.HandleFunc("GET /calculations", a.allCalculations)
	return router
}
