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
	"log"
)

type AddOwnedBook struct {
	Isbn     string
	Searched *model.Book
}

// convertBook converts plugin.Book to model.Book
func convertBook(pBook *plugin.Book) *model.Book {
	if pBook == nil {
		return nil
	}

	return &model.Book{
		Isbn13:        pBook.Isbn13,
		Isbn10:        null.StringFrom(pBook.Isbn10),
		Title:         pBook.Title,
		OriginalIsbn:  util.NullString(pBook.OriginalIsbn),
		OriginalTitle: util.NullString(pBook.OriginalTitle),
		Contents:      util.NullString(pBook.Contents),
		Url:           util.NullString(pBook.Url),
		PubDate:       util.NullString(pBook.PubDate),
		Authors:       util.NullStringJoin(pBook.Authors, ","),
		Translators:   util.NullStringJoin(pBook.Translators, ","),
		Publisher:     util.NullString(pBook.Publisher),
		Price:         util.NullDecimal(pBook.Price),
		Currency:      util.NullString(pBook.Currency),
		Thumbnail:     util.NullString(pBook.Thumbnail),
	}
}

type IBookService interface {
	Search(isbn string,
		publisher string,
		person string,
		title string,
	) ([]*model.Book, error)

	AddOwnedBook(book AddOwnedBook) (*model.OwnedBook, error)
	UpdateOwnedBook(id uint, update gin.H) (*model.OwnedBook, error)
	UpdateBook(id uint, update gin.H) (*model.Book, error)
}

type BookServiceImpl struct {
	bookPlugins []plugin.BookPlugin
}

func NewBookService() IBookService {
	bookPlugins := make([]plugin.BookPlugin, 0)
	for _, p := range plugin.GetPluginsByType(plugin.TypeBook) {
		bookPlugin := p.(plugin.BookPlugin)
		bookPlugins = append(bookPlugins, bookPlugin)
	}

	return &BookServiceImpl{
		bookPlugins: bookPlugins,
	}
}

// Search book information using bookPlugin
func (s *BookServiceImpl) Search(isbn string,
	publisher string,
	person string,
	title string,
) ([]*model.Book, error) {
	var searchFn func(string) ([]*plugin.Book, error) = nil

	books := make([]*model.Book, 0)
	for _, p := range plugin.GetPluginsByType(plugin.TypeBook) {
		bookPlugin := p.(plugin.BookPlugin)
		var query string

		if isbn != "" {
			searchFn = bookPlugin.FindByIsbn
			query = isbn
		} else if publisher != "" {
			searchFn = bookPlugin.FindByPublisher
			query = publisher
		} else if person != "" {
			searchFn = bookPlugin.FindByPerson
			query = person
		} else if title != "" {
			searchFn = bookPlugin.FindByTitle
			query = title
		} else {
			return nil, errors.New("search parameter is not set")
		}

		pBooks, err := searchFn(query)
		if err != nil {
			return nil, err
		}

		// try next plugin if not found
		if len(pBooks) == 0 {
			continue
		}

		for _, pBook := range pBooks {
			books = append(books, convertBook(pBook))
		}
	}
	return books, nil
}

func (s *BookServiceImpl) searchByIsbn(isbn string) (*model.Book, error) {
	if books, err := s.Search(isbn, "", "", ""); err != nil {
		return nil, err
	} else {
		if len(books) > 0 {
			return books[0], nil
		} else {
			return nil, nil
		}
	}
}

// AddOwnedBook adds/updates book and ownedBook.
func (s *BookServiceImpl) AddOwnedBook(body AddOwnedBook) (*model.OwnedBook, error) {
	var ownedBook *model.OwnedBook
	var err error

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		var book *model.Book

		// check book existence
		if book, err = model.FindBookByIsbn(body.Isbn); err != nil {
			return err
		}
		// insert book
		if book == nil {
			searched := body.Searched
			if searched == nil {
				if searched, err = s.searchByIsbn(body.Isbn); err != nil {
					return err
				} else if searched == nil {
					return fmt.Errorf("book not found. Isbn: %s", body.Isbn)
				}
			}

			// insert searched book
			if book, err = s.addBook(tx, searched); err != nil {
				return err
			}
		}

		// add ownedBook
		if b, err := model.FindOwnedBookByIsbn(body.Isbn); err != nil {
			return err
		} else {
			if b == nil {
				ownedBook = &model.OwnedBook{Isbn: book.Isbn13}
				return model.Save(tx, ownedBook, true)
			} else {
				return nil
			}
		}
	})

	return ownedBook, err
}

// addBook inserts model.Book to DB.
func (s *BookServiceImpl) addBook(tx *gorm.DB, book *model.Book) (*model.Book, error) {
	if err := model.Save(tx, book, true); err != nil {
		return nil, err
	} else {
		return book, nil
	}
}

// UpdateBook updates Book.
func (s *BookServiceImpl) UpdateBook(id uint, update gin.H) (*model.Book, error) {
	if book, err := model.FindBook(id); err != nil {
		return nil, err
	} else if book == nil {
		return nil, fmt.Errorf("book not found. id=%d", id)
	} else {
		if err := config.DB.Model(book).Updates(update).Error; err != nil {
			return nil, err
		}
		return book, nil
	}
}

// UpdateOwnedBook updates OwnedBook.
func (s *BookServiceImpl) UpdateOwnedBook(id uint, update gin.H) (*model.OwnedBook, error) {
	if book, err := model.FindOwnedBook(id); err != nil {
		return nil, err
	} else if book == nil {
		return nil, fmt.Errorf("owned book not found. id=%d", id)
	} else {
		log.Printf("update: %#v", update)
		if err := config.DB.Model(book).Updates(update).Error; err != nil {
			return nil, err
		} else {
			return book, nil
		}
	}
}
