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
	_ = model.DeleteAllBooks(nil)
}
func deleteAllOwnedBooks() {
	_ = model.DeleteAllOwnedBooks(nil)
}

func TestAddOwnedBookIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	deleteAllBooks()
	deleteAllOwnedBooks()

	Convey("addBook", t, func() {
		body := service.AddOwnedBook{
			Isbn: "9788960778320",
		}
		var result struct {
			Data *model.OwnedBook
		}
		rw := testServer.Post("/books/own", body, &result)

		book := result.Data
		So(rw.Code, ShouldEqual, http.StatusCreated)
		So(book.Isbn, ShouldEqual, body.Isbn)
	})
}

func TestUpdateBookSpec(t *testing.T) {
	deleteAllBooks()

	isbn := "9788960778320"
	_ = model.Save(nil, &model.Book{Isbn13: isbn}, true)

	Convey("updateBook", t, func() {
		body := gin.H{
			"title":   "foo",
			"authors": "tom,james",
			"dummy":   "foo,bar",
			"price": 100,
		}
		var result struct {
			Data *model.Book
		}
		rw := testServer.Patch(fmt.Sprintf("/books/book/%s", isbn), body, &result)

		book := result.Data
		So(rw.Code, ShouldEqual, http.StatusOK)
		So(book.Title, ShouldEqual, body["title"])
		So(book.Authors.String, ShouldEqual, body["authors"])
		So(book.Price.Decimal.IntPart(), ShouldEqual, body["price"])
	})
}

func TestUpdateOwnedBookSpec(t *testing.T) {
	deleteAllOwnedBooks()

	isbn := "9788960778320"
	_ = model.Save(nil, &model.OwnedBook{Isbn: isbn}, true)

	Convey("updateOwnedBook", t, func() {
		body := gin.H{
			"owner":        "foo",
			"paidPrice":    25000,
			"hasPaperBook": true,
		}
		var result struct {
			Data *model.OwnedBook
		}
		rw := testServer.Patch(fmt.Sprintf("/books/own/%s", isbn), body, &result)

		book := result.Data
		So(rw.Code, ShouldEqual, http.StatusOK)
		So(book.Owner.String, ShouldEqual, body["owner"])
		So(book.PaidPrice.Decimal.IntPart(), ShouldEqual, body["paidPrice"])
		So(book.HasPaperBook, ShouldEqual, body["hasPaperBook"])
	})
}
