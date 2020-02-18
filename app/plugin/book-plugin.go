package plugin

import (
	"gopkg.in/guregu/null.v3"
)

type Book struct {
	ISBN13        string    `json:"isbn13,omitempty"`
	ISBN10        string    `json:"isbn10,omitempty"`
	Title         string    `json:"title,omitempty"`
	OriginalISBN  string    `json:"originalISBN,omitempty"`
	OriginalTitle string    `json:"originalTitle,omitempty"`
	Contents      string    `json:"contents,omitempty"`
	Url           string    `json:"url,omitempty"`
	PubDate       null.Time `json:"pubDate,omitempty"`
	Authors       []string  `json:"authors,omitempty"`
	Translators   []string  `json:"translators,omitempty"`
	Publisher     string    `json:"publisher,omitempty"`
	Price         float64   `json:"price,omitempty"`
	Currency      string    `json:"currency,omitempty"`
}

type BookPlugin interface {
	FindByISBN(string) ([]*Book, error)
	FindByPerson(string) ([]*Book, error)
	FindByPublisher(string) ([]*Book, error)
	FindByTitle(string) ([]*Book, error)
}
