package model

import (
	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v3"
	"time"
)

type BaseModel struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `sql:"index" json:"deleted_at,omitempty"`
}

type Book struct {
	BaseModel
	ISBN13        string              `gorm:"size:13;unique" json:"isbn13,omitempty"`
	ISBN10        null.String         `gorm:"size:10" json:"isbn10,omitempty"`
	Title         string              `gorm:"size:255" json:"title,omitempty"`
	OriginalISBN  null.String         `gorm:"size:13" json:"originalISBN,omitempty"`
	OriginalTitle null.String         `gorm:"size:255" json:"originalTitle,omitempty"`
	Contents      null.String         `gorm:"size:8192" json:"contents,omitempty"`
	Url           null.String         `gorm:"size:1024" json:"url,omitempty"`
	PubDate       null.Time           `json:"pubDate,omitempty"`
	Authors       null.String         `gorm:"size:255" json:"authors,omitempty"`
	Translators   null.String         `gorm:"size:255" json:"translators,omitempty"`
	Publisher     null.String         `gorm:"size:255" json:"publisher,omitempty"`
	Price         decimal.NullDecimal `gorm:"type:decimal(20,2)" json:"price,omitempty"`
	Currency      null.String         `json:"currency,omitempty"`
	Thumbnail     null.String         `json:"thumbnail,omitempty"`
}

func (b *Book) TableName() string {
	return "books"
}

type InterestBooks struct {
	BaseModel
	LoginID string `gorm:"type:varchar(40);unique_index:interest_books_uq;not null"`
	Isbn    string `gorm:"type:varchar(13);unique_index:interest_books_uq;not null"`
}

func (c *InterestBooks) TableName() string { return "interest_books" }

type OwnedBook struct {
	BaseModel
	ISBN         string              `gorm:"size:13;unique" json:"ISBN13"`
	Owner        null.String         `gorm:"size:255" json:"owner"`
	AcquiredAt   null.Time           `json:"acquiredAt"`
	ScannedAt    null.Time           `json:"scannedAt"`
	PaidPrice    decimal.NullDecimal `gorm:"type:decimal(20,2)"`
	ActualPages  null.Int            `json:"actualPages"`
	HasPaperBook bool                `json:"hasPaperBook"`
}

func (b *OwnedBook) TableName() string {
	return "owned_books"
}
