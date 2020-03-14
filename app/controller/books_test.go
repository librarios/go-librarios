package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/librarios/go-librarios/app/model"
	"github.com/librarios/go-librarios/app/service"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func deleteAllBooks() {
	_ = service.DeleteAllBooks()
}
func deleteAllOwnedBooks() {
	_ = service.DeleteAllOwnedBooks()
}

func TestAddBookIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	deleteAllBooks()
	deleteAllOwnedBooks()

	Convey("addBook", t, func() {
		body := service.AddBookCommand{
			ISBN:  "9788960778320",
			Owner: "foo",
		}
		var result struct {
			Data *model.Book
		}
		rw := testServer.Post("/books", body, &result)

		book := result.Data
		So(rw.Code, ShouldEqual, http.StatusCreated)
		So(book.ISBN13, ShouldEqual, body.ISBN)
		So(book.Title, ShouldContainSubstring, "The Go Programming Language")
	})
}

func TestUpdateBookSpec(t *testing.T) {
	deleteAllBooks()

	isbn := "9788960778320"
	_, _ = service.InsertBook(&model.Book{ISBN13: isbn})

	Convey("updateBook", t, func() {
		body := gin.H{
			"title":   "foo",
			"authors": "tom,james",
			"dummy":   "foo,bar",
		}
		var result struct {
			Data *model.Book
		}
		rw := testServer.Patch(fmt.Sprintf("/books/%s", isbn), body, &result)

		book := result.Data
		So(rw.Code, ShouldEqual, http.StatusOK)
		So(book.Title, ShouldEqual, body["title"])
		So(book.Authors.String, ShouldEqual, body["authors"])
	})
}

func TestUpdateOwnedBookSpec(t *testing.T) {
	deleteAllOwnedBooks()

	isbn := "9788960778320"
	_, _ = service.InsertOwnedBook(&model.OwnedBook{
		ISBN: isbn,
	})

	Convey("updateOwnedBook", t, func() {
		body := gin.H{
			"owner":        "foo",
			"paidPrice":    25000,
			"hasPaperBook": true,
		}
		var result struct {
			Data *model.OwnedBook
		}
		rw := testServer.Patch(fmt.Sprintf("/books/%s/owned", isbn), body, &result)

		book := result.Data
		So(rw.Code, ShouldEqual, http.StatusOK)
		So(book.Owner.String, ShouldEqual, body["owner"])
		So(book.PaidPrice.Float64, ShouldEqual, body["paidPrice"])
		So(book.HasPaperBook, ShouldEqual, body["hasPaperBook"])
	})
}
