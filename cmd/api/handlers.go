package main

import (
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
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	sum := nums.A + nums.B
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
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	subtract := nums.A - nums.B
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
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	multiply := nums.A * nums.B
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
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if nums.B == 0 {
		a.logger.Error("Cannot Divide by Zero")
		http.Error(w, "Cannot Divide by Zero", http.StatusBadRequest)
		return
	}

	divide := float64(nums.A) / float64(nums.B)
	jsonResponse := fmt.Sprintf("{\"result\":%0.2f}\n", divide)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResponse))
}
