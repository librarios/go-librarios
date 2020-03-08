package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/librarios/go-librarios/app/model"
	"github.com/librarios/go-librarios/app/plugin"
	"github.com/librarios/go-librarios/app/service"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestBookSearchSpec(t *testing.T) {
	if !isIntegrationTest() {
		return
	}

	url := "/book/search"

	sendGet := func() (service.IBookService, *gin.Engine) {
		s := service.NewBookService()
		r := newRouter()
		r.GET(url, searchBookHandlerFn(s))
		return s, r
	}

	Convey("search by isbn", t, func() {
		_, r := sendGet()
		req, rw := makeGet(url + "?isbn=9788960778320")
		r.ServeHTTP(rw, req)

		var result struct {
			Data []*plugin.Book
		}
		_ = parseJSON(rw, &result)
		So(rw.Code, ShouldEqual, http.StatusOK)
		So(len(result.Data), ShouldEqual, 1)
		book := result.Data[0]
		So(book.Title, ShouldContainSubstring, "The Go Programming Language")
	})
}

func TestAddBookSpec(t *testing.T) {
	if !isIntegrationTest() {
		return
	}

	url := "/book"

	sendPost := func() (service.IBookService, *gin.Engine) {
		s := service.NewBookService()
		r := newRouter()
		r.POST(url, addBookHandlerFn(s))
		return s, r
	}

	Convey("addBook", t, func() {
		_, r := sendPost()
		body := service.AddBookCommand{
			ISBN:  "9788960778320",
			Owner: "foo",
		}
		req, rw, err := makePostJSON(url, body)
		So(err, ShouldBeNil)

		r.ServeHTTP(rw, req)

		var result struct {
			Data *model.Book
		}
		_ = parseJSON(rw, &result)

		So(rw.Code, ShouldEqual, http.StatusCreated)
		So(result.Data, ShouldNotBeNil)
		book := result.Data
		So(book.ISBN13, ShouldEqual, body.ISBN)
		So(book.Title, ShouldContainSubstring, "The Go Programming Language")
	})
}
