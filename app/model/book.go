package model

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/guregu/null.v3"
)

type Book struct {
	gorm.Model
	ISBN          string      `gorm:"size:13;unique"`
	Title         string      `gorm:"size:255"`
	OriginalISBN  null.String `gorm:"size:13"`
	OriginalTitle null.String `gorm:"size:255"`
	Contents      null.String `gorm:"size:8192"`
	Url           null.String `gorm:"size:1024"`
	PubDate       null.Time
	Authors       null.String `gorm:"size:255"`
	Translators   null.String `gorm:"size:255"`
	Publisher     null.String `gorm:"size:255"`
	Price         null.Float
	Currency      null.String
}

type OwnedBook struct {
	gorm.Model
	ISBN        string      `gorm:"size:13;unique"`
	Owner       null.String `gorm:"size:255"`
	AcquiredAt  null.Time
	ScannedAt   null.Time
	PaidPrice   null.Float
	ActualPages null.Int
}
