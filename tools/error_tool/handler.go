package error_tool

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/raythx98/url-shortener/dto"
	"net/http"
)

func Handle(w http.ResponseWriter, err error) {
	fmt.Errorf("error: %v", err)

	var appError *AppError
	if errors.Is(err, appError) {
		HandleAppError(w, appError)
		return
	}

	//var invalidValidationErr *validator.InvalidValidationError
	var invalidValidationErr *validator.InvalidValidationError
	if errors.As(err, &invalidValidationErr) {
		HandleInvalidValidationError(w, invalidValidationErr)
		return
	}

	//var validationErr
	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		HandleValidationError(w, validationErr)
		return
	}

	HandleInternalServerError(w, err)
}

func HandleAppError(w http.ResponseWriter, appError *AppError) {
	marshal, err := json.Marshal(dto.ErrorResponse{
		Message: appError.message,
		Code:    appError.code,
		Data:    appError.err.Error(),
	})
	if err != nil {
		HandleInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(marshal)
}

func HandleInvalidValidationError(w http.ResponseWriter, validationErr *validator.InvalidValidationError) {
	marshal, err := json.Marshal(dto.ErrorResponse{
		Message: "Validation Error",
		Code:    422,
		Data:    validationErr.Error(),
	})
	if err != nil {
		HandleInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusUnprocessableEntity)
	_, _ = w.Write(marshal)
}

func HandleValidationError(w http.ResponseWriter, validationErr validator.ValidationErrors) {
	marshal, err := json.Marshal(dto.ErrorResponse{
		Message: "Validation Error",
		Code:    422,
		Data:    validationErr.Error(),
	})
	if err != nil {
		HandleInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusUnprocessableEntity)
	_, _ = w.Write(marshal)
}

func HandleInternalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	marshal, err := json.Marshal(dto.ErrorResponse{
		Message: "Internal Server Error",
		Code:    500,
		Data:    err.Error(),
	})
	if err != nil {
		_, _ = w.Write([]byte("Internal Server Error"))
	}

	_, _ = w.Write(marshal)
}
