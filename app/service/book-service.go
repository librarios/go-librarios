package service

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/librarios/go-librarios/app/model"
	"github.com/librarios/go-librarios/app/plugin"
	"github.com/librarios/go-librarios/app/util"
	"gopkg.in/guregu/null.v3"
)

type AddBookCommand struct {
	ISBN         string  `json:"isbn"`
	Owner        string  `json:"owner"`
	AcquiredAt   string  `json:"acquiredAt"`
	ScannedAt    string  `json:"scannedAt"`
	PaidPrice    float64 `json:"paidPrice"`
	ActualPages  int64   `json:"actualPages"`
	HasPaperBook bool    `json:"hasPaperBook"`
}

type IBookService interface {
	Search(isbn string,
		publisher string,
		person string,
		title string,
	) ([]*plugin.Book, error)

	AddBook(book AddBookCommand) (*model.Book, error)
}

type BookService struct {
	bookPlugins []plugin.BookPlugin
}

func NewBookService() IBookService {
	bookPlugins := make([]plugin.BookPlugin, 0)
	for _, p := range plugin.GetPluginsByType(plugin.TypeBook) {
		bookPlugin := p.(plugin.BookPlugin)
		bookPlugins = append(bookPlugins, bookPlugin)
	}

	return &BookService{
		bookPlugins: bookPlugins,
	}
}

// Search book information using bookPlugin
func (s *BookService) Search(isbn string,
	publisher string,
	person string,
	title string,
) ([]*plugin.Book, error) {
	var fn func(string) ([]*plugin.Book, error) = nil

	for _, p := range plugin.GetPluginsByType(plugin.TypeBook) {
		bookPlugin := p.(plugin.BookPlugin)
		var query string

		if isbn != "" {
			fn = bookPlugin.FindByISBN
			query = isbn
		} else if publisher != "" {
			fn = bookPlugin.FindByPublisher
			query = publisher
		} else if person != "" {
			fn = bookPlugin.FindByPerson
			query = person
		} else if title != "" {
			fn = bookPlugin.FindByTitle
			query = title
		} else {
			return nil, errors.New("search parameter is not set")
		}

		books, err := fn(query)
		if err != nil {
			return nil, err
		}

		// try next plugin if not found
		if len(books) == 0 {
			continue
		}

		return books, nil
	}
	return make([]*plugin.Book, 0), nil
}

func (s *BookService) AddBook(cmd AddBookCommand) (*model.Book, error) {
	bookModel := model.Book{}

	err := dbConn.Transaction(func(tx *gorm.DB) error {
		// check book existence
		if err := tx.Where(&model.Book{ISBN13: cmd.ISBN}).First(&bookModel).Error; err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				return err
			}

			// search book info
			books, err := s.Search(cmd.ISBN, "", "", "")
			if err != nil {
				return err
			}

			// insert book
			if len(books) > 0 {
				book := books[0]
				bookModel.ISBN13 = book.ISBN13
				bookModel.ISBN10 = null.StringFrom(book.ISBN10)
				bookModel.Title = book.Title
				bookModel.OriginalISBN = util.NullString(book.OriginalISBN)
				bookModel.OriginalTitle = util.NullString(book.OriginalTitle)
				bookModel.Contents = util.NullString(book.Contents)
				bookModel.Url = util.NullString(book.Url)
				bookModel.PubDate = book.PubDate
				bookModel.Authors = util.NullStringJoin(book.Authors, ",")
				bookModel.Translators = util.NullStringJoin(book.Translators, ",")
				bookModel.Publisher = util.NullString(book.Publisher)
				bookModel.Price = util.NullFloat(book.Price)
				bookModel.Currency = util.NullString(book.Currency)
			} else {
				if len(cmd.ISBN) == 13 {
					bookModel.ISBN13 = cmd.ISBN
				} else {
					return fmt.Errorf("cannot add unknown ISBN10: %s", cmd.ISBN)
				}
			}
			tx.Create(&bookModel)
		}

		// add ownedBook
		ownedBook := model.OwnedBook{}
		insert := false
		if err := tx.Where(&model.OwnedBook{ISBN: cmd.ISBN}).First(&ownedBook).Error; err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				return err
			}
			ownedBook.ISBN = bookModel.ISBN13
			insert = true
		}
		ownedBook.Owner = util.NullString(cmd.Owner)
		ownedBook.AcquiredAt = util.NullTimeFromString(cmd.AcquiredAt)
		ownedBook.ScannedAt = util.NullTimeFromString(cmd.ScannedAt)
		ownedBook.PaidPrice = util.NullFloat(cmd.PaidPrice)
		ownedBook.ActualPages = util.NullInt(cmd.ActualPages)
		ownedBook.HasPaperBook = cmd.HasPaperBook

		if insert {
			tx.Create(&ownedBook)
		} else {
			tx.Save(&ownedBook)
		}
		return nil
	})

	return &bookModel, err
}
