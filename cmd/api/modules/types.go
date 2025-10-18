package modules

import (
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

// RouterParams is used to inject routers into handlers
type RouterParams struct {
	fx.In

	Router    *mux.Router `name:"main"`
	APIRouter *mux.Router `name:"api"`
}
