package database

import (
	"go_crud_postgres/internal/config"
	"log"

	"go_crud_postgres/internal/models"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// func Connect(cfg *config.Config) (*pgxpool.Pool, error) {
// 	db, err := pgxpool.New(context.Background(), cfg.GetDBURL())
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to database: %w", err)
// 	}

// 	err = db.Ping(context.Background())
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to ping database: %w", err)
// 	}

//		log.Println("Connected to PostgreSQL")
//		return db, nil
//	}

func Connect(cfg *config.Config) error {
	// Connect using GORM
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	// AutoMigrate - creates/updates tables
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	DB = db
	log.Println("✅ Connected to PostgreSQL with GORM")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
