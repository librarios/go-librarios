package main

type Book struct {
	ISBN          string
	Title         string
	Contents      string
	Url           string
	PubDate       string
	Authors       []string
	Publisher     string
	Translators   []string
	Price         float64
	PriceCurrency string
	IsSale        bool
}