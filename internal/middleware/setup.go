package middleware

import (
	"github.com/gorilla/mux"
)

func Setup(router *mux.Router, serviceName string) {
	SetupOTEL(router, serviceName)
	SetupRateLimiter(router)
	SetupJWTAuth(router)
}
