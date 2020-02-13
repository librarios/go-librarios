package service

import (
	"encoding/json"
	"fmt"
	"github.com/librarios/go-librarios/app/model"
	"github.com/librarios/go-librarios/app/util"
	"gopkg.in/guregu/null.v3"
	"io/ioutil"
	"net/http"
	"strings"
)

// https://developers.kakao.com/docs/restapi/search#%EC%B1%85-%EA%B2%80%EC%83%89

var kakaoDef = PluginDef{
	Type: PluginTypeBook,
	Name: "kakao",
	NewFunc: func() Plugin {
		return &Kakao{}
	},
}

type Kakao struct {
	apiKey string
}

func (k *Kakao) Type() string {
	return PluginTypeBook
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

func (k *Kakao) searhBook(target string, query string) ([]*model.Book, error) {
	u := fmt.Sprintf("https://dapi.kakao.com/v3/search/book?target=%s&query=%s", target, query)
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

	books := make([]*model.Book, 0)

	for _, doc := range kakaoResp.Documents {
		book := model.Book{
			ISBN:        doc.Isbn,
			Title:       doc.Title,
			Contents:    null.StringFrom(doc.Contents),
			Url:         null.StringFrom(doc.Url),
			PubDate:     util.NullTimeFromString(doc.Datetime),
			Authors:     null.StringFrom(strings.Join(doc.Authors, ",")),
			Publisher:   null.StringFrom(doc.Publisher),
			Translators: null.StringFrom(strings.Join(doc.Translators, ",")),
			Price:       null.FloatFrom(float64(doc.Price)),
			Currency:    null.StringFrom("KRW"),
		}
		books = append(books, &book)
	}

	return books, nil
}

func (k *Kakao) FindByISBN(isbn string) ([]*model.Book, error) {
	return k.searhBook("isbn", isbn)
}

func (k *Kakao) FindByTitle(title string) ([]*model.Book, error) {
	return k.searhBook("title", title)
}

func (k *Kakao) FindByPublisher(publisher string) ([]*model.Book, error) {
	return k.searhBook("publisher", publisher)
}

func (k *Kakao) FindByPerson(person string) ([]*model.Book, error) {
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
