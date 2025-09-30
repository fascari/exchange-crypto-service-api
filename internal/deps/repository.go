package deps

import (
	accountrepo "exchange-crypto-service-api/internal/app/account/repository"
	exchangerepo "exchange-crypto-service-api/internal/app/exchange/repository"
	userrepo "exchange-crypto-service-api/internal/app/user/repository"

	"gorm.io/gorm"
)

type Repositories struct {
	User     userrepo.Repository
	Exchange exchangerepo.Repository
	Account  accountrepo.Repository
}

func initRepos(db *gorm.DB) Repositories {
	return Repositories{
		User:     userrepo.New(db),
		Exchange: exchangerepo.New(db),
		Account:  accountrepo.New(db),
	}
}
