package deps

import (
	"exchange-crypto-service-api/internal/infra"
)

type (
	Dependencies struct {
		Repositories Repositories
	}
)

func New(app infra.App) Dependencies {
	return Dependencies{
		Repositories: initRepos(app.DB),
	}
}
