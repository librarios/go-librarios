package model

import (
	"github.com/jinzhu/gorm"
)

// FindBook finds a book by id
// returns (nil, nil) if not found.
func FindBook(id uint) (*Book, error) {
	book := new(Book)
	if err := DB(nil).First(book, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return book, nil
}

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

// FindOwnedBook finds owned book by id
// returns (nil, nil) if not found.
func FindOwnedBook(id uint) (*OwnedBook, error) {
	ownedBook := new(OwnedBook)
	if err := DB(nil).First(ownedBook, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return ownedBook, nil
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
