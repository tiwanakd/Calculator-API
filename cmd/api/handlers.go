package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Numbers struct {
	A int `json:"a"`
	B int `json:"b"`
}

func (a *api) add(w http.ResponseWriter, r *http.Request) {
	var nums Numbers
	err := a.decodeJSONBody(w, r, &nums)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			a.genericServerError(w, r, err)
		}
		return
	}

	sum := nums.A + nums.B
	err = a.calculations.Insert("Addition", nums.A, nums.B, float64(sum))
	if err != nil {
		a.genericServerError(w, r, err)
		return
	}

	jsonResponse := fmt.Sprintf("{\"result\":%d}\n", sum)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResponse))
}

func (a *api) subtract(w http.ResponseWriter, r *http.Request) {
	var nums Numbers
	err := a.decodeJSONBody(w, r, &nums)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			a.genericServerError(w, r, err)
		}
		return
	}

	subtract := nums.A - nums.B
	err = a.calculations.Insert("Subtraction", nums.A, nums.B, float64(subtract))
	if err != nil {
		a.genericServerError(w, r, err)
		return
	}
	jsonResponse := fmt.Sprintf("{\"result\":%d}\n", subtract)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResponse))
}

func (a *api) multiply(w http.ResponseWriter, r *http.Request) {
	var nums Numbers
	err := a.decodeJSONBody(w, r, &nums)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			a.genericServerError(w, r, err)
		}
		return
	}

	multiply := nums.A * nums.B
	err = a.calculations.Insert("Multiplication", nums.A, nums.B, float64(multiply))
	if err != nil {
		a.genericServerError(w, r, err)
		return
	}
	jsonResponse := fmt.Sprintf("{\"result\":%d}\n", multiply)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResponse))
}

func (a *api) divide(w http.ResponseWriter, r *http.Request) {
	var nums Numbers
	err := a.decodeJSONBody(w, r, &nums)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			a.genericServerError(w, r, err)
		}
		return
	}

	if nums.B == 0 {
		a.logger.Error("Divsion by Zero")
		http.Error(w, "Cannot Divide by Zero", http.StatusBadRequest)
		return
	}

	divide := float64(nums.A) / float64(nums.B)
	err = a.calculations.Insert("Division", nums.A, nums.B, float64(divide))
	if err != nil {
		a.genericServerError(w, r, err)
		return
	}
	jsonResponse := fmt.Sprintf("{\"result\":%0.2f}\n", divide)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResponse))
}

func (a *api) allCalculations(w http.ResponseWriter, r *http.Request) {
	calculations, err := a.calculations.GetAll()
	if err != nil {
		a.genericServerError(w, r, err)
		return
	}
	jsonData, err := json.MarshalIndent(calculations, "", "\t")
	if err != nil {
		a.genericServerError(w, r, err)
		return
	}

	w.Write(jsonData)
}

func (a *api) getCalculations(w http.ResponseWriter, r *http.Request) {
	operation := r.URL.Query().Get("operation")
	calculations, err := a.calculations.GetCalculations(operation)
	if err != nil {
		a.genericServerError(w, r, err)
		return
	}

	jsonData, err := json.MarshalIndent(calculations, "", "\t")
	if err != nil {
		a.genericServerError(w, r, err)
		return
	}

	w.Write(jsonData)
}

func (a *api) homeView(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	calculations, err := a.calculations.GetLatestCalculations()
	if err != nil {
		a.genericServerError(w, r, err)
		return
	}

	data := templateData{
		Calculations: calculations,
	}

	a.render(w, r, http.StatusOK, "home.tmpl.html", data)
}
