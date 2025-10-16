package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"exchange-crypto-service-api/pkg/apperror"

	"github.com/rs/zerolog/log"
)

// WriteJSON handles JSON encoding and writing to response writer
func WriteJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Error().Err(err).Msg("failed to encode JSON response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// ReadJSON handles JSON decoding from request body
func ReadJSON(r *http.Request, payload any) error {
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		log.Error().Err(err).Msg("failed to decode JSON request")
		return err
	}
	return nil
}

// WriteError writes an error response as JSON
func WriteError(w http.ResponseWriter, statusCode int, err error) {
	statusText := strings.ToUpper(strings.ReplaceAll(http.StatusText(statusCode), " ", "_"))
	appErr := apperror.New(statusText, "%s", err.Error())
	WriteJSON(w, statusCode, appErr)
}
