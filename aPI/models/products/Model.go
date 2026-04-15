package products

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	CreatedAt time.Time      `gorm:json:"created_at"`
	UpdatedAt time.Time      `gorm:json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
