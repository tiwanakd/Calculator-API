package main

import "net/http"

func (a *api) routes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /add", a.add)
	router.HandleFunc("POST /subtract", a.subtract)
	router.HandleFunc("POST /multiply", a.multiply)
	router.HandleFunc("POST /divide", a.divide)
	router.HandleFunc("GET /allcalculations", a.allCalculations)
	router.HandleFunc("GET /getcalculations", a.getCalculations)
	router.HandleFunc("GET /{$}", a.homeView)

	return router
}
