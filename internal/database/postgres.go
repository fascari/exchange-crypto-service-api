package database

import (
	"fmt"
	"log"

	"exchange-crypto-service-api/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectPostgres(dbConfig config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbConfig.Conn.Host,
		dbConfig.Conn.Username,
		dbConfig.Conn.Password,
		dbConfig.Conn.DatabaseName,
		dbConfig.Conn.Port,
		dbConfig.Conn.SslMode)

	if dbConfig.Conn.Schema != "" {
		dsn += fmt.Sprintf(" search_path=%s", dbConfig.Conn.Schema)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(dbConfig.MaxLifetimeConnections)

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("successfully connected to database")
	return db, nil
}
