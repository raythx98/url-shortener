package error_helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/raythx98/url-shortener/dto"

	"github.com/go-playground/validator/v10"
)

func Handle(w http.ResponseWriter, err error) {
	fmt.Errorf("error: %v", err)

	var appError *AppError
	if errors.Is(err, appError) {
		HandleAppError(w, appError)
		return
	}

	var invalidValidationErr *validator.InvalidValidationError
	if errors.As(err, &invalidValidationErr) {
		HandleInvalidValidationError(w, invalidValidationErr)
		return
	}

	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		HandleValidationError(w, validationErr)
		return
	}

	HandleInternalServerError(w, err)
}

func HandleAppError(w http.ResponseWriter, appError *AppError) {
	marshal, err := json.Marshal(&dto.ErrorResponse{
		Message: appError.Message(),
		Code:    appError.Code(),
		Data:    appError.Error(),
	})
	if err != nil {
		HandleInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(marshal)
}

func HandleInvalidValidationError(w http.ResponseWriter, validationErr *validator.InvalidValidationError) {
	marshal, err := json.Marshal(dto.NewValidationError(validationErr))
	if err != nil {
		HandleInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusUnprocessableEntity)
	_, _ = w.Write(marshal)
}

func HandleValidationError(w http.ResponseWriter, validationErr validator.ValidationErrors) {
	marshal, err := json.Marshal(dto.NewValidationError(validationErr))
	if err != nil {
		HandleInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusUnprocessableEntity)
	_, _ = w.Write(marshal)
}

func HandleInternalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	marshal, err := json.Marshal(dto.NewInternalServerError(err))
	if err != nil {
		_, _ = w.Write([]byte("Internal Server Error"))
	}

	_, _ = w.Write(marshal)
}
