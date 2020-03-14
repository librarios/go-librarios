package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
		ownedBook, err := FindOwnedBookByISBN(cmd.ISBN)
		if err != nil {
			return err
		}

		insert := false
		if ownedBook == nil {
			ownedBook = &model.OwnedBook{
				ISBN: bookModel.ISBN13,
			}
			insert = true
		}
		ownedBook.Owner = util.NullString(cmd.Owner)
		ownedBook.AcquiredAt = util.NullTimeFromString(cmd.AcquiredAt)
		ownedBook.ScannedAt = util.NullTimeFromString(cmd.ScannedAt)
		ownedBook.PaidPrice = util.NullFloat(cmd.PaidPrice)
		ownedBook.ActualPages = util.NullInt(cmd.ActualPages)
		ownedBook.HasPaperBook = cmd.HasPaperBook

		if insert {
			tx.Create(ownedBook)
		} else {
			tx.Save(ownedBook)
		}
		return nil
	})

	return &bookModel, err
}

func (s *BookService) UpdateBook(isbn string, update gin.H) (*model.Book, error) {
	var err error
	book, err := FindBookByISBN(isbn)

	if err != nil {
		return nil, err
	}

	if err := dbConn.Model(book).Updates(update).Error; err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookService) UpdateOwnedBook(isbn string, update gin.H) (*model.OwnedBook, error) {
	var err error
	book, err := FindOwnedBookByISBN(isbn)

	if err != nil {
		return nil, err
	}

	if err := dbConn.Model(book).Updates(update).Error; err != nil {
		return nil, err
	}

	return book, nil
}

func FindBookByISBN(isbn string) (*model.Book, error) {
	book := &model.Book{}
	book.ISBN13 = isbn

	if err := dbConn.Where(book).First(book).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return book, nil
}

func FindOwnedBookByISBN(isbn string) (*model.OwnedBook, error) {
	book := &model.OwnedBook{}
	book.ISBN = isbn

	if err := dbConn.Where(book).First(book).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return book, nil
}

func InsertBook(book *model.Book) (*model.Book, error) {
	if err := dbConn.Create(book).Error; err != nil {
		return nil, err
	}
	return book, nil
}

func InsertOwnedBook(book *model.OwnedBook) (*model.OwnedBook, error) {
	if err := dbConn.Create(book).Error; err != nil {
		return nil, err
	}
	return book, nil
}

func DeleteAllBooks() error {
	return dbConn.Unscoped().Delete(model.Book{}).Error
}

func DeleteAllOwnedBooks() error {
	return dbConn.Unscoped().Delete(model.OwnedBook{}).Error
}
