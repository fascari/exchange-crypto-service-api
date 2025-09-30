package finduserbalance

import (
	"encoding/json"
	"net/http"
	"strconv"

	"exchange-crypto-service-api/internal/app/user/usecase/finduserbalance"

	"github.com/gorilla/mux"
)

const Path = "/users/{id}/balance"

type Handler struct {
	useCase finduserbalance.UseCase
}

func NewHandler(useCase finduserbalance.UseCase) Handler {
	return Handler{
		useCase: useCase,
	}
}

func RegisterEndpoint(r *mux.Router, h Handler) {
	r.HandleFunc(Path, h.Handle).Methods(http.MethodGet)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userIDStr := mux.Vars(r)["id"]

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	userBalance, err := h.useCase.Execute(r.Context(), uint(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(userBalance); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
