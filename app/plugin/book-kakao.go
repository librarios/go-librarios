package plugin

import (
	"encoding/json"
	"fmt"
	"github.com/librarios/go-librarios/app/util"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// https://developers.kakao.com/docs/restapi/search#%EC%B1%85-%EA%B2%80%EC%83%89

type Kakao struct {
	apiKey string
}

func (k *Kakao) Type() string {
	return TypeBook
}

func (k *Kakao) Name() string {
	return "kakao"
}

func (k *Kakao) SetProperties(p map[string]interface{}) {
	apiKey, exists := p["apiKey"]
	if exists {
		k.apiKey = apiKey.(string)
	}
}

func (k *Kakao) searhBook(target string, query string) ([]*Book, error) {
	u := fmt.Sprintf("https://dapi.kakao.com/v3/search/book?target=%s&query=%s", target, query)
	log.Printf("[GET] %s", u)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("KakaoAK %s", k.apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var kakaoResp KakaoResponse
	if err := json.Unmarshal(bytes, &kakaoResp); err != nil {
		return nil, err
	}

	books := make([]*Book, 0)

	for _, doc := range kakaoResp.Documents {
		books = append(books, k.toBook(doc))
	}

	return books, nil
}

func (k *Kakao) toBook(doc KakaoDocument) *Book {
	book := &Book{}

	for _, token := range strings.Split(doc.Isbn, " ") {
		switch len(token) {
		case 10:
			book.ISBN10 = token
		case 13:
			book.ISBN13 = token
		}
	}
	book.Title = doc.Title
	book.Contents = doc.Contents
	book.Url = doc.Url
	book.PubDate = util.NullTimeFromString(doc.Datetime)
	book.Authors = doc.Authors
	book.Publisher = doc.Publisher
	book.Translators = doc.Translators
	book.Price = float64(doc.Price)
	book.Currency = "KRW"

	return book
}

func (k *Kakao) FindByISBN(isbn string) ([]*Book, error) {
	return k.searhBook("isbn", isbn)
}

func (k *Kakao) FindByTitle(title string) ([]*Book, error) {
	return k.searhBook("title", title)
}

func (k *Kakao) FindByPublisher(publisher string) ([]*Book, error) {
	return k.searhBook("publisher", publisher)
}

func (k *Kakao) FindByPerson(person string) ([]*Book, error) {
	return k.searhBook("person", person)
}

// Kakao API definitions

type KakaoMeta struct {
	IsEnd         bool `json:"is_end"`
	PageableCount int  `json:"pageable_count"`
	TotalCount    int  `json:"total_count"`
}

type KakaoDocument struct {
	Title       string   `json:"title"`
	Contents    string   `json:"contents"`
	Url         string   `json:"url"`
	Isbn        string   `json:"isbn"`
	Datetime    string   `json:"datetime"`
	Authors     []string `json:"authors"`
	Publisher   string   `json:"publisher"`
	Translators []string `json:"translators"`
	Price       int      `json:"price"`
	SalePrice   int      `json:"sale_price"`
	Thumbnail   string   `json:"thumbnail"`
	Status      string   `json:"status"`
}

type KakaoResponse struct {
	Meta      KakaoMeta       `json:"meta"`
	Documents []KakaoDocument `json:"documents"`
}
