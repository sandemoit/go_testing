package models

import (
	"time"
)

type ProdDailies struct {
	ID           uint    `gorm:"primaryKey"`
	UUID         string  `gorm:"size:36;not null;index"`
	KodeMaterial string  `gorm:"size:15;not null;index"`
	WCID         string  `gorm:"size:15;not null;index"`
	JmlProd      float64 `gorm:"type:decimal"`
	TglProd      string  `gorm:"type:date;index"`
	Status       uint    `gorm:"type:int;index"`
	CreatedBy    string  `gorm:"type:string" json:"created_by,omitempty"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"default:null"`

	// Definisi relasi
	Plant    Plant    `json:"Plant" gorm:"foreignKey:WCID;references:Plant"`
	Material Material `json:"Material" gorm:"foreignKey:KodeMaterial;references:KodeMat"`
}

type ProdDailiesResponse struct {
	ID           uint       `json:"ID"`
	UUID         string     `json:"UUID"`
	KodeMaterial string     `json:"KodeMaterial"`
	WCID         string     `json:"WCID"`
	JmlProd      int        `json:"JmlProd"`
	TglProd      time.Time  `json:"TglProd"`
	Status       int        `json:"Status"`
	CreatedAt    time.Time  `json:"CreatedAt"`
	UpdatedAt    time.Time  `json:"UpdatedAt"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"default:null"`
	Plant        Plant      `json:"Plant"`
	Material     Material   `json:"Material"`
}

type Material struct {
	ID        uint   `gorm:"primaryKey"`
	KodeMat   string `gorm:"size:15;not null;unique"`
	Material  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"default:null"`
}

type Plant struct {
	ID        uint   `gorm:"primaryKey"`
	Plant     string `gorm:"size:15;not null;unique"`
	PlantDesc string `gorm:"type:string;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"default:null"`
}

type UpdateProdDailiesRequest struct {
	Method string        `json:"method"`
	Data   []ProdDailies `json:"data"`
}
