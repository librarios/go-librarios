package model

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type Book struct {
	gorm.Model
	ISBN          string         `gorm:"size:13;unique"`
	Title         string         `gorm:"size:255"`
	OriginalISBN  sql.NullString `gorm:"size:13"`
	OriginalTitle sql.NullString `gorm:"size:255"`
	Contents      sql.NullString `gorm:"size:8192"`
	Url           sql.NullString `gorm:"size:1024"`
	PubDate       sql.NullTime
	Authors       sql.NullString `gorm:"size:255"`
	Translators   sql.NullString `gorm:"size:255"`
	Publisher     sql.NullString `gorm:"size:255"`
	Price         sql.NullFloat64
	Currency      sql.NullString
}

type OwnedBook struct {
	gorm.Model
	ISBN         string         `gorm:"size:13;unique"`
	Owner        sql.NullString `gorm:"size:255"`
	AcquiredAt   sql.NullTime
	ScannedAt    sql.NullTime
	PaidPrice    sql.NullFloat64
	ScannedPages sql.NullInt32
}
