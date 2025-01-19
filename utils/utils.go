package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func WriteValidationError(w http.ResponseWriter, err error) {
	validationErrors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		validationErrors[err.Field()] = getValidationMessage(err)
	}

	errorResponse := map[string]interface{}{
		"error": "validation error",
		"details": validationErrors,
	}
	WriteJSON(w, http.StatusBadRequest, errorResponse)
}

func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("'%s' is required", fe.Field())
	case "email":
		return fmt.Sprintf("'%s' is invalid", fe.Field())
	case "min":
		return fmt.Sprintf("'%s' is less than the required minimum length", fe.Field())
	case "max":
		return fmt.Sprintf("'%s' is greater than the required maximum length", fe.Field())
	}
	return fmt.Sprintf("%v", fe.Error())
}