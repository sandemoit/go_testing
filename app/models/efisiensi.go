package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type ItemEfisiensi struct {
	ID            uint           `gorm:"primaryKey"`
	IDArea        uint           `gorm:"not null"`
	EfisiensiDesc string         `gorm:"type:string;not null"`
	CreatedAt     *time.Time     `json:"created_at,omitempty"`
	UpdatedAt     *time.Time     `json:"updated_at,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	Area             Area               `gorm:"foreignKey:IDArea"`
	EfisiensiDailies []EfisiensiDailies `gorm:"foreignKey:IDItemEfisiensi"`
}

type Area struct {
	ID        uint           `gorm:"primaryKey"`
	Area      string         `gorm:"type:string;not null"`
	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type EfisiensiDailies struct {
	ID              uint           `gorm:"primaryKey"`
	UUID            string         `gorm:"size:36;not null"`
	IDItemEfisiensi uint           `gorm:"not null"`
	Aktual          float64        `gorm:"type:float;not null"`
	TglEfisiensi    time.Time      `gorm:"type:date;not null"`
	CreatedBy       string         `gorm:"type:string;not null" json:"created_by,omitempty"`
	CreatedAt       *time.Time     `json:"created_at,omitempty"`
	UpdatedAt       *time.Time     `json:"updated_at,omitempty"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	ItemEfisiensi ItemEfisiensi `gorm:"foreignKey:IDItemEfisiensi"`
}

type EfisiensiRow struct {
	Periode string             `json:"periode"`
	UUID    string             `json:"UUID"`
	Values  map[string]float64 `json:"-"` // Menyimpan nilai sebelum diubah ke JSON
}

func GetAllEfisiensi(db *gorm.DB, IDArea uint, page, limit int) ([]map[string]interface{}, int64, error) {
	var items []ItemEfisiensi
	var response []EfisiensiRow
	itemDescSet := make(map[string]bool)
	var totalItems int64

	// Hitung total items sebelum pagination diterapkan
	err := db.Model(&ItemEfisiensi{}).Where("id_area = ?", IDArea).Count(&totalItems).Error
	if err != nil {
		return nil, 0, err
	}

	// Terapkan pagination
	offset := (page - 1) * limit
	err = db.Preload("EfisiensiDailies", func(db *gorm.DB) *gorm.DB {
		return db.Select("id_item_efisiensi, aktual, uuid, tgl_efisiensi")
	}).
		Where("id_area = ?", IDArea).
		Offset(offset).Limit(limit).
		Find(&items).Error

	if err != nil {
		return nil, 0, err
	}

	// Kumpulan periode unik
	periodeMap := make(map[string]map[string]float64)
	uuidMap := make(map[string]string)

	// Loop semua data untuk membentuk struktur JSON dinamis
	for _, item := range items {
		for _, daily := range item.EfisiensiDailies {
			tanggal := daily.TglEfisiensi.Format("2006-01-02")

			// Simpan item_desc unik agar nanti jadi header di JSON
			itemDescSet[strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(item.EfisiensiDesc, " ", "_"), "/", "_"))] = true

			// Jika belum ada data untuk tanggal ini, buat dulu
			if _, exists := periodeMap[tanggal]; !exists {
				periodeMap[tanggal] = make(map[string]float64)
				uuidMap[tanggal] = daily.UUID
			}

			// Masukkan nilai ke dalam map
			periodeMap[tanggal][strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(item.EfisiensiDesc, " ", "_"), "/", "_"))] = daily.Aktual
		}
	}

	// Ubah map menjadi slice dengan format yang sesuai
	for tanggal, values := range periodeMap {
		row := EfisiensiRow{
			Periode: tanggal,
			UUID:    uuidMap[tanggal],
			Values:  values,
		}
		response = append(response, row)
	}

	// Pastikan setiap row memiliki semua item_desc yang tersedia
	var finalResponse []map[string]interface{}
	for _, row := range response {
		rowData := map[string]interface{}{
			"periode": row.Periode,
			"uuid":    row.UUID,
		}

		// Tambahkan semua item_desc dengan nilai default 0 jika tidak ada
		for desc := range itemDescSet {
			if val, exists := row.Values[desc]; exists {
				rowData[desc] = val
			} else {
				rowData[desc] = 0
			}
		}

		finalResponse = append(finalResponse, rowData)
	}

	return finalResponse, totalItems, nil
}
