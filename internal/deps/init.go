package deps

import (
	"exchange-crypto-service-api/internal/infra"
)

type (
	Dependencies struct {
		Repositories Repositories
		UseCases     UseCases
	}
)

func New(app infra.App) Dependencies {
	repos := initRepos(app.DB)
	return Dependencies{
		Repositories: repos,
		UseCases:     initUseCases(repos),
	}
}
