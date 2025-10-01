package createuser

import (
	"net/http"

	"exchange-crypto-service-api/internal/app/user/usecase/createuser"
	httpjson "exchange-crypto-service-api/pkg/http"

	"github.com/gorilla/mux"
)

const Path = "/users"

type Handler struct {
	useCase createuser.UseCase
}

func NewHandler(useCase createuser.UseCase) Handler {
	return Handler{
		useCase: useCase,
	}
}

func RegisterEndpoint(r *mux.Router, h Handler) {
	r.HandleFunc(Path, h.Handle).Methods(http.MethodPost)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var payload InputPayload
	if err := httpjson.ReadJSON(r, &payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if err := payload.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := payload.ToDomain()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.useCase.Execute(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, ToOutputPayload(createdUser))
}
