package main

import "net/http"

func (a *api) routes() *http.ServeMux {
	router := http.NewServeMux()

	//static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	//JSON API
	router.HandleFunc("POST /add", a.add)
	router.HandleFunc("POST /subtract", a.subtract)
	router.HandleFunc("POST /multiply", a.multiply)
	router.HandleFunc("POST /divide", a.divide)
	router.HandleFunc("GET /allcalculations", a.allCalculations)
	router.HandleFunc("GET /getcalculations", a.getCalculations)

	//UI
	router.HandleFunc("GET /{$}", a.homeView)
	router.HandleFunc("GET /calculation/{id}", a.calculationView)

	return router
}
