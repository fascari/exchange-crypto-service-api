package createtransaction

import (
	"errors"
	"net/http"
	"strconv"

	"exchange-crypto-service-api/internal/app/transaction/domain"
	"exchange-crypto-service-api/internal/app/transaction/usecase/createtransaction"
	httpjson "exchange-crypto-service-api/pkg/http"

	"github.com/gorilla/mux"
)

const Path = "/transactions/{transactionType}/accounts/{accountID}"

type Handler struct {
	useCase createtransaction.UseCase
}

func NewHandler(useCase createtransaction.UseCase) Handler {
	return Handler{useCase: useCase}
}

func RegisterEndpoint(r *mux.Router, h Handler) {
	r.HandleFunc(Path, h.Handle).Methods(http.MethodPost)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	accountID, err := strconv.ParseUint(vars["accountID"], 10, 32)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, errors.New("invalid account ID"))
		return
	}

	var payload InputPayload
	if err := httpjson.ReadJSON(r, &payload); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, err)
		return
	}

	transactionType, err := domain.ParseTransactionType(vars["transactionType"])
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, err)
		return
	}

	transactionID, err := h.useCase.Execute(r.Context(), transactionType, uint(accountID), payload.Amount, payload.IdempotencyKey)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := OutputPayload{
		TransactionID: transactionID,
	}

	httpjson.WriteJSON(w, http.StatusCreated, response)
}
