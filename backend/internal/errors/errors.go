package errors

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Message string            `json:"error"`
	Fields  map[string]string `json:"fields,omitempty"`
	Code    int               `json:"-"`
}

func (e *APIError) Error() string {
	return e.Message
}

func BadRequest(message string) *APIError {
	return &APIError{Message: message, Code: http.StatusBadRequest}
}

func BadRequestValidation(fields map[string]string) *APIError {
	return &APIError{Message: "validation failed", Fields: fields, Code: http.StatusBadRequest}
}

func Unauthorized() *APIError {
	return &APIError{Message: "unauthorized", Code: http.StatusUnauthorized}
}

func Forbidden() *APIError {
	return &APIError{Message: "forbidden", Code: http.StatusForbidden}
}

func NotFound() *APIError {
	return &APIError{Message: "not found", Code: http.StatusNotFound}
}

func Conflict(message string) *APIError {
	return &APIError{Message: message, Code: http.StatusConflict}
}

func InternalServerError() *APIError {
	return &APIError{Message: "internal server error", Code: http.StatusInternalServerError}
}

func WriteError(w http.ResponseWriter, err *APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Message})
}

func WriteErrorWithFields(w http.ResponseWriter, err *APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(err)
}
