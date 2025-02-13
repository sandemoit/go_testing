package repositories

import (
	"go_teknologi/app/database"

	"gorm.io/gorm"
)

// Scope untuk query dasar
func QueryScope(status int) *gorm.DB {
	return database.DB.
		Preload("Plant").
		Preload("Material").
		Select(`id, uuid, kode_material, wc_id, jml_prod, tgl_prod, status, created_at`).
		Group(`id, uuid, kode_material, wc_id, jml_prod, tgl_prod, status, created_at`).
		Where("status = ? AND deleted_at IS NULL", status).
		Order("created_at DESC")
}
