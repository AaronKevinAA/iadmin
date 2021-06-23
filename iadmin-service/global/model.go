package global

import (
	"gorm.io/gorm"
)

type GvaModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt int64 `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}