package database

import (
	"go_teknologi/app/models"
)

// Jalankan migrasi
func runMigration() {
	DB.AutoMigrate(
		// &models.User{},
		// &models.Role{},
		// &models.Permission{},
		// &models.UserRole{},
		// &models.RolePermission{},
		&models.Material{},
		&models.Plant{},
		&models.ProdDailies{},
		&models.Area{},
		&models.ItemEfisiensi{},
		&models.EfisiensiDailies{},
	)
}
