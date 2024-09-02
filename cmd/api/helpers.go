package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Create a new type that will hold malformed requests
// this will implement the Error Interface
type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func (a *api) decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			a.logger.Error("request body content type is not application/json")
			msg := "Content Type Header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	// Use http.MaxBytesReader to enforce a maximum read of 1MB from the
	// response body. A request body larger than that will now result in
	// Decode() returning a "http: request body too large" error.
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	//setup the json decoder
	//call DisallowUnknownFields() so there are not dissallowed fields passed in
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshallTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errors in the JSON and send an error message
		case errors.As(err, &syntaxError):
			a.logger.Error(err.Error())
			msg := fmt.Sprintf("Request body contains badly formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON
		case errors.Is(err, io.ErrUnexpectedEOF):
			a.logger.Error(err.Error())
			msg := fmt.Sprintln("Request body contains badly formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// check if there are invalid types provided
		case errors.As(err, &unmarshallTypeError):
			a.logger.Error(err.Error())
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshallTypeError.Field, unmarshallTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// Catch the error caused by extra unexpected fields in the request
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			a.logger.Error(err.Error())
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// An io.EOF error is returned by Decode() if the request body is empty
		case errors.Is(err, io.EOF):
			a.logger.Error("no input provided in the body", "err:", err.Error())
			msg := "Request Body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// Catch the error caused by the request body being too large.
		case err.Error() == "http: request body too large":
			a.logger.Error(err.Error())
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			a.logger.Error(err.Error())
			return err
		}

	}
	// Call decode again, using a pointer to an empty anonymous struct as
	// the destination. If the request body only contained a single JSON
	// object this will return an io.EOF error. So if we get anything else,
	// we know that there is additional data in the request body.
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		a.logger.Error("recevied additional data besides the JSON Object", "err:", err)
		msg := "Request body must only containg a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil

}
