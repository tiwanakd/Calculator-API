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

	//UI
	dynamic := alice.New(a.sessionManager.LoadAndSave, noSurf, a.authenticate)

	router.Handle("GET /{$}", dynamic.ThenFunc(a.homeView))
	router.Handle("GET /calculation/{id}", dynamic.ThenFunc(a.calculationView))
	router.Handle("GET /user/signup", dynamic.ThenFunc(a.userSignup))
	router.Handle("POST /user/signup", dynamic.ThenFunc(a.userSignupPost))
	router.Handle("GET /user/login", dynamic.ThenFunc(a.userLogin))
	router.Handle("POST /user/login", dynamic.ThenFunc(a.userLoginPost))

	protected := dynamic.Append(a.requireAuthentication)

	router.Handle("GET /createcalculation", protected.ThenFunc(a.createCalulationView))
	router.Handle("POST /createcalculation", protected.ThenFunc(a.createCalulationPost))
	router.Handle("POST /user/logout", protected.ThenFunc(a.userLogoutPost))

	standard := alice.New(a.recoverPanic, a.logRequest, commonHeaders)
	return standard.Then(router)
}
