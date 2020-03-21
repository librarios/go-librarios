package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/librarios/go-librarios/app/config"
	"github.com/librarios/go-librarios/app/model"
	"github.com/librarios/go-librarios/app/plugin"
	"github.com/librarios/go-librarios/app/util"
	"gopkg.in/guregu/null.v3"
)

type AddBookCommand struct {
	ISBN         string
	Owner        string
	AcquiredAt   string
	ScannedAt    string
	PaidPrice    float64
	ActualPages  int64
	HasPaperBook bool
	Detail       *plugin.Book
}

type UpdateBookCommand struct {
	ISBN          string
	Owner         null.String
	AcquiredAt    null.String
	ScannedAt     null.String
	PaidPrice     null.Float
	ActualPages   null.Int
	HasPaperBook  null.Bool
	Title         null.String
	OriginalISBN  null.String
	OriginalTitle null.String
	Contents      null.String
	Url           null.String
	PubDate       null.Time
	Authors       null.String
	Translators   null.String
	Publisher     null.String
	Price         null.Float
	Currency      null.String
}

type IBookService interface {
	Search(isbn string,
		publisher string,
		person string,
		title string,
	) ([]*plugin.Book, error)

	AddBook(book AddBookCommand) (*model.Book, error)
	UpdateBook(isbn string, update gin.H) (*model.Book, error)
	UpdateOwnedBook(isbn string, update gin.H) (*model.OwnedBook, error)
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

// AddBook adds/updates book and ownedBook.
func (s *BookService) AddBook(cmd AddBookCommand) (*model.Book, error) {
	book := new(model.Book)

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// check existingBook existence
		if existingBook, err := model.FindBookByISBN(cmd.ISBN); err != nil {
			return err
		} else if existingBook == nil {
			// search pBook info if not pBook info is not supplied to 'detail' parameter.
			pBook := cmd.Detail
			if pBook == nil {
				books, err := s.Search(cmd.ISBN, "", "", "")
				if err != nil {
					return err
				}

				if len(books) > 0 {
					pBook = books[0]
				}
			}

			// insert pBook
			if pBook != nil {
				book.ISBN13 = pBook.ISBN13
				book.ISBN10 = null.StringFrom(pBook.ISBN10)
				book.Title = pBook.Title
				book.OriginalISBN = util.NullString(pBook.OriginalISBN)
				book.OriginalTitle = util.NullString(pBook.OriginalTitle)
				book.Contents = util.NullString(pBook.Contents)
				book.Url = util.NullString(pBook.Url)
				book.PubDate = pBook.PubDate
				book.Authors = util.NullStringJoin(pBook.Authors, ",")
				book.Translators = util.NullStringJoin(pBook.Translators, ",")
				book.Publisher = util.NullString(pBook.Publisher)
				book.Price = util.NullFloat(pBook.Price)
				book.Currency = util.NullString(pBook.Currency)
			} else {
				if len(cmd.ISBN) == 13 {
					book.ISBN13 = cmd.ISBN
				} else {
					return fmt.Errorf("cannot add unknown ISBN10: %s", cmd.ISBN)
				}
			}
			if err := model.Save(tx, book, true); err != nil {
				return err
			}
		} else {
			// set existing book as result
			book = existingBook
		}

		// add/update ownedBook
		ownedBook, err := model.FindOwnedBookByISBN(cmd.ISBN)
		if err != nil {
			return err
		}

		insert := false
		if ownedBook == nil {
			ownedBook = &model.OwnedBook{ISBN: book.ISBN13}
			insert = true
		}
		ownedBook.Owner = util.NullString(cmd.Owner)
		ownedBook.AcquiredAt = util.NullTimeFromString(cmd.AcquiredAt)
		ownedBook.ScannedAt = util.NullTimeFromString(cmd.ScannedAt)
		ownedBook.PaidPrice = util.NullFloat(cmd.PaidPrice)
		ownedBook.ActualPages = util.NullInt(cmd.ActualPages)
		ownedBook.HasPaperBook = cmd.HasPaperBook

		return model.Save(tx, ownedBook, insert)
	})

	return book, err
}

func (s *BookService) UpdateBook(isbn string, update gin.H) (*model.Book, error) {
	var err error
	book, err := model.FindBookByISBN(isbn)

	if err != nil {
		return nil, err
	}

	if err := config.DB.Model(book).Updates(update).Error; err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookService) UpdateOwnedBook(isbn string, update gin.H) (*model.OwnedBook, error) {
	var err error
	book, err := model.FindOwnedBookByISBN(isbn)

	if err != nil {
		return nil, err
	}

	if err := config.DB.Model(book).Updates(update).Error; err != nil {
		return nil, err
	}

	return book, nil
}
