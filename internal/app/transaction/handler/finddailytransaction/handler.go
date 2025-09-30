package finddailytransaction

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"exchange-crypto-service-api/internal/app/transaction/usecase/finddailytransaction"

	"github.com/gorilla/mux"
)

const Path = "/transactions/daily"

type Handler struct {
	useCase finddailytransaction.UseCase
}

func NewHandler(useCase finddailytransaction.UseCase) Handler {
	return Handler{
		useCase: useCase,
	}
}

func RegisterEndpoint(r *mux.Router, h Handler) {
	r.HandleFunc(Path, h.Handle).Methods(http.MethodGet)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if err := validateRequiredParams(startDateStr, endDateStr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startDate, endDate, err := parseDateParams(startDateStr, endDateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dailyTrans, err := h.useCase.Execute(r.Context(), startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, ToOutputPayloads(dailyTrans))
}

func validateRequiredParams(startDate, endDate string) error {
	if startDate == "" || endDate == "" {
		return errors.New("start_date and end_date query parameters are required")
	}
	return nil
}

func parseDateParams(startDateStr, endDateStr string) (startDate, endDate time.Time, err error) {
	startDate, err = parseDate(startDateStr, "start_date")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endDate, err = parseDate(endDateStr, "end_date")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return startDate, endDate, nil
}

func parseDate(dateStr, paramName string) (time.Time, error) {
	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid %s format, use YYYY-MM-DD", paramName)
	}
	return date, nil
}

func writeJSONResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
