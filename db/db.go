package db

import "gorm.io/gorm"

type DB interface {
	GetInstance() *gorm.DB
}
