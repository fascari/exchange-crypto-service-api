package bootstrap

import (
	"exchange-crypto-service-api/internal/config"
	"exchange-crypto-service-api/internal/database"

	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Database = fx.Module("database",
	fx.Provide(NewDatabase),
)

func NewDatabase(cfg config.Database) (*gorm.DB, error) {
	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		return nil, err
	}
	return db, nil
}
