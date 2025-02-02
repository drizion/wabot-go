package db

import (
	"log"

	"github.com/drizion/wabot-go/database/models"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB
	err   error
	Query *models.Query
)

// DbConnection create database connection
func DbConnection(masterDSN string) error {
	var db = DB

	logDb := viper.GetBool("DB_DEBUG")

	loglevel := logger.Silent
	if logDb {
		loglevel = logger.Info
	}

	db, err = gorm.Open(postgres.New(postgres.Config{
		DriverName: "postgres",
		DSN:        masterDSN,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(loglevel),
	})

	if err != nil {
		log.Fatalf("Db connection error")
		return err
	}
	DB = db

	// Migrate the schema
	// Query = models.Use(db)

	return nil
}

// GetDB connection
func GetDB() *gorm.DB {
	return DB
}
