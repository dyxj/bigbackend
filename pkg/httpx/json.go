package httpx

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/dyxj/bigbackend/pkg/errorx"
)

const contentTypeJSON = "application/json"
const headerKeyContentType = "Content-Type"

const internalServerErrorDefaultMessage = "internal server error"

func JsonResponse(statusCode int, resp any, w http.ResponseWriter) {
	// Why not directly in the response writer?
	// In the event encoding fails, we would be able to
	// change the status code accordingly.
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(true)

	if err := encoder.Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(headerKeyContentType, contentTypeJSON)
	w.WriteHeader(statusCode)

	// Errors are likely due to client disconnect.
	_, _ = w.Write(buf.Bytes())
}

func BadRequestResponse(message string, details map[string]string, w http.ResponseWriter) {
	JsonResponse(
		http.StatusBadRequest,
		ErrorResponse{
			Message: message,
			Details: details,
		},
		w)
}

func ConflictResponse(message string, details map[string]string, w http.ResponseWriter) {
	JsonResponse(
		http.StatusConflict,
		ErrorResponse{
			Message: message,
			Details: details,
		},
		w)
}

func ValidationFailedResponse(validationFailure *errorx.ValidationError, w http.ResponseWriter) {
	JsonResponse(
		http.StatusBadRequest,
		ErrorResponse{
			Message: "validation failed",
			Details: validationFailure.Properties,
		},
		w)
}

func InternalServerErrorResponse(message string, w http.ResponseWriter) {
	if message == "" {
		message = internalServerErrorDefaultMessage
	}
	JsonResponse(
		http.StatusInternalServerError,
		ErrorResponse{Message: message},
		w)
}
