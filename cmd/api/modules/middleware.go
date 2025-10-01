package modules

import (
	"exchange-crypto-service-api/internal/middleware"

	"github.com/gorilla/mux"
)

func setupMiddlewares(router *mux.Router, serviceName string) {
	middleware.SetupOTEL(router, serviceName)
	middleware.SetupRateLimiter(router)
}
