package model

import (
	"github.com/jinzhu/gorm"
)

// FindBookByIsbn finds a book with matching Isbn.
// returns (nil, nil) if not found.
func FindBookByIsbn(isbn string) (*Book, error) {
	book := new(Book)
	book.Isbn13 = isbn

	if err := DB(nil).Where(book).First(book).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return book, nil
}

// FindOwnedBookByIsbn finds an owned book with matching Isbn.
// returns (nil, nil) if not found.
func FindOwnedBookByIsbn(isbn string) (*OwnedBook, error) {
	book := new(OwnedBook)
	book.Isbn = isbn

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
func DeleteAllOwnedBooks(tx *gorm.DB) error {
	return DB(tx).Unscoped().Delete(OwnedBook{}).Error
}
