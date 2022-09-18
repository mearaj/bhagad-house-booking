package service

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	BookingID uint `gorm:"foreignKey"`
	Note      string
}
