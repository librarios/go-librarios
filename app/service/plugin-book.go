package service

import (
	"github.com/librarios/go-librarios/app/model"
)

type BookPlugin interface {
	FindByISBN(string) ([]*model.Book, error)
	FindByPerson(string) ([]*model.Book, error)
	FindByPublisher(string) ([]*model.Book, error)
	FindByTitle(string) ([]*model.Book, error)
}
