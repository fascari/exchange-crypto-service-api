package deps

import (
	createuserrepo "exchange-crypto-service-api/internal/app/user/repository/createuser"

	"gorm.io/gorm"
)

type Repositories struct {
	User createuserrepo.Repository
}

func initRepos(db *gorm.DB) Repositories {
	return Repositories{
		User: createuserrepo.New(db),
	}
}
