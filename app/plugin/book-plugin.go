package plugin

import (
	"github.com/shopspring/decimal"
)

type Book struct {
	Isbn13        string          `json:"isbn13,omitempty"`
	Isbn10        string          `json:"isbn10,omitempty"`
	Title         string          `json:"title,omitempty"`
	OriginalIsbn  string          `json:"originalIsbn,omitempty"`
	OriginalTitle string          `json:"originalTitle,omitempty"`
	Contents      string          `json:"contents,omitempty"`
	Url           string          `json:"url,omitempty"`
	PubDate       string          `json:"pubDate,omitempty"`
	Authors       []string        `json:"authors,omitempty"`
	Translators   []string        `json:"translators,omitempty"`
	Publisher     string          `json:"publisher,omitempty"`
	Price         decimal.Decimal `json:"price,omitempty"`
	Currency      string          `json:"currency,omitempty"`
	Thumbnail     string          `json:"thumbnail,omitempty"`
}

type BookPlugin interface {
	FindByIsbn(string) ([]*Book, error)
	FindByPerson(string) ([]*Book, error)
	FindByPublisher(string) ([]*Book, error)
	FindByTitle(string) ([]*Book, error)
}
