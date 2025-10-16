package tokengen

import (
	"errors"
	"net/http"
	"strconv"

	"exchange-crypto-service-api/internal/app/auth/usecase/generatetoken"
	httpjson "exchange-crypto-service-api/pkg/http"

	"github.com/gorilla/mux"
)

const Path = "/tokens/generate/{user_id}"

type Handler struct {
	useCase generatetoken.UseCase
}

func NewHandler(useCase generatetoken.UseCase) Handler {
	return Handler{useCase: useCase}
}

func RegisterEndpoint(r *mux.Router, h Handler) {
	r.HandleFunc(Path, h.Handle).Methods(http.MethodPost)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]

	if userIDStr == "" {
		httpjson.WriteError(w, http.StatusBadRequest, errors.New("user_id path parameter is required"))
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, errors.New("user_id must be a valid positive integer"))
		return
	}

	token, err := h.useCase.Execute(r.Context(), uint(userID))
	if err != nil {
		httpjson.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, toOutputPayload(token))
}
