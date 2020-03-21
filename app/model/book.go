package model

import (
	"github.com/jinzhu/gorm"
)

// FindBookByISBN finds a book with matching ISBN.
// returns (nil, nil) if not found.
func FindBookByISBN(isbn string) (*Book, error) {
	book := new(Book)
	book.ISBN13 = isbn

	if err := DB(nil).Where(book).First(book).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return book, nil
}

// FindOwnedBookByISBN finds an owned book with matching ISBN.
// returns (nil, nil) if not found.
func FindOwnedBookByISBN(isbn string) (*OwnedBook, error) {
	book := new(OwnedBook)
	book.ISBN = isbn

	if err := DB(nil).Where(book).First(book).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return book, nil
}

// DeleteAllBooks deletes all books.
func DeleteAllBooks(tx *gorm.DB) error {
	return DB(tx).Unscoped().Delete(Book{}).Error
}

// DeleteAllOwnBooks deletes all owned books.
func DeleteAllOwnedBooks(tx *gorm.DB, ) error {
	return DB(tx).Unscoped().Delete(OwnedBook{}).Error
}
