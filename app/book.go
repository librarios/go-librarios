package app

type KakaoBookResponse struct {
	Isbn          string
	Title         string
	Contents      string
	Url           string
	PubDate       string
	Authors       []string
	Publisher     string
	Translators   []string
	Price         float64
	PriceCurrency string
	Pages         int32
	IsSale        bool
}
