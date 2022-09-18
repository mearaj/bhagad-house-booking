package service

import (
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	Name    string
	Email   string
	Contact string
	Address string
}
