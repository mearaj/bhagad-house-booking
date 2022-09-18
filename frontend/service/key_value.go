package service

import "gorm.io/gorm"

type KeyValue struct {
	gorm.Model
	Key        string
	CustomerID string `gorm:"foreignKey"`
}
