package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tiwanakd/Calculator-API/internal/models"
	"github.com/tiwanakd/Calculator-API/internal/validator"
)

type Numbers struct {
	A int `json:"a"`
	B int `json:"b"`
}

func (a *api) calculate(operation string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var result float64

		switch operation {
		case "Addition":
			result = float64(nums.A) + float64(nums.B)
		case "Subtraction":
			result = float64(nums.A) - float64(nums.B)
		case "Multiplication":
			result = float64(nums.A) * float64(nums.B)
		case "Division":
			if nums.B == 0 {
				a.logger.Error("Divsion by Zero")
				http.Error(w, "Cannot Divide by Zero", http.StatusBadRequest)
				return
			}
			result = float64(nums.A) / float64(nums.B)
		}
		_, err = a.calculations.Insert(operation, nums.A, nums.B, result)
		if err != nil {
			a.genericServerError(w, r, err)
			return
		}
		jsonResponse := fmt.Sprintf("{\"result\":%0.2f}\n", result)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jsonResponse))
	}
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
	calculations, err := a.calculations.GetLatestCalculations()
	if err != nil {
		a.genericServerError(w, r, err)
		return
	}

	data := a.newTemplateData(r)
	data.Calculations = calculations

	a.render(w, r, http.StatusOK, "home.tmpl.html", data)
}

func (a *api) calculationView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		a.genericServerError(w, r, err)
		return
	}

	calculation, err := a.calculations.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.NotFound(w, r)
		} else {
			a.genericServerError(w, r, err)
		}
		return
	}

	data := a.newTemplateData(r)
	data.Calculation = calculation

	a.render(w, r, http.StatusOK, "calculation.tmpl.html", data)
}

func (a *api) createCalulationView(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	data.Form = resultForm{}
	a.render(w, r, http.StatusOK, "createCalculation.tmpl.html", data)
}

type resultForm struct {
	Id     int
	A      int
	B      int
	Result float64
	validator.Validator
}

func (a *api) createCalulationPost(w http.ResponseWriter, r *http.Request) {
	var form resultForm
	err := r.ParseForm()
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(r.PostForm.Get("numberA")), "numberA", "This field cannot be blank")
	form.CheckField(validator.NotBlank(r.PostForm.Get("numberB")), "numberB", "This field cannot be blank")

	operation := r.Form.Get("submitbtn")
	numberA, err := strconv.Atoi(r.Form.Get("numberA"))
	if err != nil {
		form.AddFieldError("numberA", "Invalid Number Provided")
	}
	numberB, err := strconv.Atoi(r.Form.Get("numberB"))
	if err != nil {
		form.AddFieldError("numberB", "Invalid Number Provided")
	}

	var result float64

	switch operation {
	case "Addition":
		result = float64(numberA) + float64(numberB)
	case "Subtraction":
		result = float64(numberA) - float64(numberB)
	case "Multiplication":
		result = float64(numberA) * float64(numberB)
	case "Division":
		if numberB == 0 {
			form.AddFieldError("numberB", "Cannot divide by Zero")
		}
		result = float64(numberA) / float64(numberB)
	}

	form.A = numberA
	form.B = numberB

	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form
		a.render(w, r, http.StatusUnprocessableEntity, "createCalculation.tmpl.html", data)
		return
	}

	id, err := a.calculations.Insert(operation, numberA, numberB, result)
	if err != nil {
		a.genericServerError(w, r, err)
		return
	}

	form.Id = id
	form.Result = result

	data := a.newTemplateData(r)
	data.Form = form

	//a.sessionManager.Put(r.Context(), "resultflash", result)

	a.render(w, r, http.StatusOK, "createCalculation.tmpl.html", data)
}

type userSignupForm struct {
	Name     string
	Email    string
	Password string
	validator.Validator
}

func (a *api) userSignup(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	data.Form = userSignupForm{}
	a.render(w, r, http.StatusOK, "signup.tmpl.html", data)
}

func (a *api) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := r.ParseForm()
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	form.Name = r.PostForm.Get("name")
	form.Email = r.PostForm.Get("email")
	form.Password = r.PostForm.Get("password")

	form.CheckField(validator.NotBlank(form.Name), "name", "This field Cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field Cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field Cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be al least 8 characters long")

	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form
		a.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	err = a.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address alredy in use")

			data := a.newTemplateData(r)
			data.Form = form
			a.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			a.genericServerError(w, r, err)
		}

		return
	}

	a.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email    string
	Password string
	validator.Validator
}

func (a *api) userLogin(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	data.Form = userLoginForm{}
	a.render(w, r, http.StatusOK, "login.tmpl.html", data)
}

func (a *api) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := r.ParseForm()
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	form.Email = r.PostForm.Get("email")
	form.Password = r.PostForm.Get("password")

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be empty")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be empty")

	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form
		a.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	id, err := a.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentails) {
			form.AddNonFieldError("Incorrect Username or Password")

			data := a.newTemplateData(r)
			data.Form = form
			a.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)

		} else {
			a.genericServerError(w, r, err)
		}

		return
	}

	err = a.sessionManager.RenewToken(r.Context())
	if err != nil {
		a.genericServerError(w, r, err)
		return
	}

	a.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/createcalculation", http.StatusSeeOther)

}

func (a *api) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := a.sessionManager.RenewToken(r.Context())
	if err != nil {
		a.genericServerError(w, r, err)
		return
	}

	a.sessionManager.Remove(r.Context(), "authenticatedUserID")
	a.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
