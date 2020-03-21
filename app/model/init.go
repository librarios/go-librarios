package model

import (
	"github.com/jinzhu/gorm"
	"github.com/librarios/go-librarios/app/config"
	"log"
)

// AutoMigrate run auto migration for given models, will only add missing fields, won't delete/change current data
func AutoMigrate() {
	config.DB.AutoMigrate(
		new(Book),
		new(OwnedBook),
	)
	log.Println("DB auto migration finished.")
}

// DB returns 'tx' argument if not nil.
// returns default DB connection if 'tx' is nil.
func DB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return config.DB
}

// Save inserts/updates an entity.
func Save(tx *gorm.DB, entity interface{}, insert bool) error {
	if insert {
		return DB(tx).Save(entity).Error
	} else {
		return DB(tx).Update(entity).Error
	}
}
