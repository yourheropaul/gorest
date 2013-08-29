package gorest

import (
	"log"
	"net/http"
)

var error_handlers = make(map[int]ErrorHandler, 0)

//Signiture of functions to be used as Authorizers
type ErrorHandler func(w http.ResponseWriter, r *http.Request)

//Registers an Authorizer for the specified realm.
func RegisterErrorHandler(code int, handler ErrorHandler) {

	if _, found := error_handlers[code]; !found {
		error_handlers[code] = handler
	}
}

//Returns the registred Authorizer for the specified realm.
func GetErrorHandler(code int) ErrorHandler {

	eh, exists := error_handlers[code]

	if !exists {
		return DefaulErrorHandler
	}

	return eh
}

func DefaulErrorHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Default error message")
	w.Header().Add("error-description", "Resource not found.")
	w.WriteHeader(http.StatusNotFound)
}
