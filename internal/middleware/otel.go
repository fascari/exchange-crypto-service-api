package middleware

import (
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func SetupOTEL(r *mux.Router, serviceName string) {
	r.Use(otelmux.Middleware(serviceName))
}
