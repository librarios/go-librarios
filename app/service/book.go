package service

import (
	"errors"
	"github.com/librarios/go-librarios/app/model"
)

type IBookService interface {
	Search(isbn string,
		publisher string,
		person string,
		title string,
	) ([]*model.Book, error)
}

type BookService struct {
	bookPlugins []BookPlugin
}

func NewBookService() IBookService {
	bookPlugins := make([]BookPlugin, 0)
	for _, plugin := range pluginManager.GetPluginsByType(PluginTypeBook) {
		bookPlugin := plugin.(BookPlugin)
		bookPlugins = append(bookPlugins, bookPlugin)
	}

	return &BookService{
		bookPlugins: bookPlugins,
	}
}

func (s *BookService) Search(isbn string,
	publisher string,
	person string,
	title string,
) ([]*model.Book, error) {
	var fn func(string) ([]*model.Book, error) = nil

	for _, plugin := range pluginManager.GetPluginsByType(PluginTypeBook) {
		bookPlugin := plugin.(BookPlugin)
		if isbn != "" {
			fn = bookPlugin.FindByISBN
		} else if publisher != "" {
			fn = bookPlugin.FindByPublisher
		} else if person != "" {
			fn = bookPlugin.FindByPerson
		} else if title != "" {
			fn = bookPlugin.FindByTitle
		} else {
			return nil, errors.New("search parameter is not set")
		}

		books, err := fn(title)
		if err != nil {
			return nil, err
		}

		// try next plugin if not found
		if len(books) == 0 {
			continue
		}

		return books, nil
	}
	return make([]*model.Book, 0), nil
}
