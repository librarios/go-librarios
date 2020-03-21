package model

import (
	"gopkg.in/guregu/null.v3"
	"time"
)

type BaseModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
}

type Book struct {
	BaseModel
	ISBN13        string      `gorm:"size:13;unique" json:"isbn13,omitempty"`
	ISBN10        null.String `gorm:"size:10" json:"isbn10,omitempty"`
	Title         string      `gorm:"size:255" json:"title,omitempty"`
	OriginalISBN  null.String `gorm:"size:13" json:"originalISBN,omitempty"`
	OriginalTitle null.String `gorm:"size:255" json:"originalTitle,omitempty"`
	Contents      null.String `gorm:"size:8192" json:"contents,omitempty"`
	Url           null.String `gorm:"size:1024" json:"url,omitempty"`
	PubDate       null.Time   `json:"pubDate,omitempty"`
	Authors       null.String `gorm:"size:255" json:"authors,omitempty"`
	Translators   null.String `gorm:"size:255" json:"translators,omitempty"`
	Publisher     null.String `gorm:"size:255" json:"publisher,omitempty"`
	Price         null.Float  `json:"price,omitempty"`
	Currency      null.String `json:"currency,omitempty"`
}

func (b *Book) TableName() string {
	return "books"
}

type OwnedBook struct {
	BaseModel
	ISBN         string      `gorm:"size:13;unique" json:"ISBN13"`
	Owner        null.String `gorm:"size:255" json:"owner"`
	AcquiredAt   null.Time   `json:"acquiredAt"`
	ScannedAt    null.Time   `json:"scannedAt"`
	PaidPrice    null.Float  `json:"paidPrice"`
	ActualPages  null.Int    `json:"actualPages"`
	HasPaperBook bool        `json:"hasPaperBook"`
}

func (b *OwnedBook) TableName() string {
	return "owned_books"
}
