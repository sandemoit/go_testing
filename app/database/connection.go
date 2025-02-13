package database

import (
	"fmt"
	"go_teknologi/utils"
	"os"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.LogError(err)
	}

	DB = db

	sqlDB, err := db.DB()
	if err != nil {
		utils.LogError(err)
	}

	if err := sqlDB.Ping(); err != nil {
		utils.LogError(err)
	}

	// Jalankan migrasi
	runMigration()

	utils.LogInfo("Connected to SQL Server!")
}
