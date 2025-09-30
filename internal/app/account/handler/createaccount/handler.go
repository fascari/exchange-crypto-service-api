package createaccount

import (
	"encoding/json"
	"net/http"

	"exchange-crypto-service-api/internal/app/account"
	"exchange-crypto-service-api/internal/app/account/usecase/createaccount"
	"exchange-crypto-service-api/pkg/apperror"

	"github.com/gorilla/mux"
)

const Path = "/accounts"

type (
	Handler struct {
		useCase createaccount.UseCase
	}
)

func NewHandler(useCase createaccount.UseCase) Handler {
	return Handler{
		useCase: useCase,
	}
}

func RegisterEndpoint(r *mux.Router, h Handler) {
	r.HandleFunc(Path, h.Handle).Methods(http.MethodPost)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var payload InputPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if err := payload.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdAccount, err := h.useCase.Execute(r.Context(), payload.ToDomain())
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(ToOutputPayload(createdAccount)); err != nil {
		handleError(w, err)
	}
}

func handleError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	if apperror.As(err, account.ErrCodeInvalidMinimumAge) {
		statusCode = http.StatusBadRequest
	}
	http.Error(w, err.Error(), statusCode)
}
