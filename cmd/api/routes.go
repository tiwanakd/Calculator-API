package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (a *api) routes() http.Handler {
	router := http.NewServeMux()

	//static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	//JSON API
	router.HandleFunc("POST /add", a.calculate("Addition"))
	router.HandleFunc("POST /subtract", a.calculate("Subtraction"))
	router.HandleFunc("POST /multiply", a.calculate("Multiplication"))
	router.HandleFunc("POST /divide", a.calculate("Division"))
	router.HandleFunc("GET /allcalculations", a.allCalculations)
	router.HandleFunc("GET /getcalculations", a.getCalculations)

	dynamic := alice.New(a.sessionManager.LoadAndSave)

	//UI
	router.Handle("GET /{$}", dynamic.ThenFunc(a.homeView))
	router.Handle("GET /calculation/{id}", dynamic.ThenFunc(a.calculationView))
	router.Handle("GET /createCalculation", dynamic.ThenFunc(a.createCalulationView))
	router.Handle("POST /createCalculation", dynamic.ThenFunc(a.createCalulationPost))

	standard := alice.New(a.recoverPanic, a.logRequest, commonHeaders)
	return standard.Then(router)
}
