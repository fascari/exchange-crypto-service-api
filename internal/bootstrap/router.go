package bootstrap

import (
	"exchange-crypto-service-api/internal/middleware"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

var Router = fx.Module("router",
	fx.Provide(NewRouters),
)

type RouterOut struct {
	fx.Out

	MainRouter *mux.Router `name:"main"`
	APIRouter  *mux.Router `name:"api"`
}

func NewRouters() RouterOut {
	mainRouter := mux.NewRouter()
	apiRouter := mainRouter.PathPrefix("/api/v1").Subrouter()
	middleware.Setup(apiRouter, serviceName)

	return RouterOut{
		MainRouter: mainRouter,
		APIRouter:  apiRouter,
	}
}
