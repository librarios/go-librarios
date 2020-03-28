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
	book := model.Book{Isbn13: isbn}
	_ = model.Save(nil, &book, true)

	Convey("updateBook", t, func() {
		body := gin.H{
			"title":    "foo",
			"authors":  "tom,james",
			"dummy":    "foo,bar",
			"price":    100,
			"currency": "USD",
		}
		var result struct {
			Data *model.Book
		}
		rw := testServer.Patch(fmt.Sprintf("/books/book/%d", book.ID), body, &result)

		res := result.Data
		So(rw.Code, ShouldEqual, http.StatusOK)
		So(res.Title, ShouldEqual, body["title"])
		So(res.Authors.String, ShouldEqual, body["authors"])
		So(res.Price.Decimal.IntPart(), ShouldEqual, body["price"])
		So(res.Currency.String, ShouldEqual, body["currency"])
	})
}

func TestUpdateOwnedBookSpec(t *testing.T) {
	deleteAllOwnedBooks()

	isbn := "9788960778320"
	ownedBook := model.OwnedBook{Isbn: isbn}
	_ = model.Save(nil, &ownedBook, true)

	Convey("updateOwnedBook", t, func() {
		body := gin.H{
			"owner":        "foo",
			"acquiredAt":   "2020-01-02",
			"scannedAt":    "2020-01-03",
			"paidPrice":    25000,
			"hasPaperBook": true,
		}
		var result struct {
			Data *model.OwnedBook
		}
		rw := testServer.Patch(fmt.Sprintf("/books/own/%d", ownedBook.ID), body, &result)

		res := result.Data
		So(rw.Code, ShouldEqual, http.StatusOK)
		So(res.Owner.String, ShouldEqual, body["owner"])
		So(res.AcquiredAt.String, ShouldEqual, body["acquiredAt"])
		So(res.ScannedAt.String, ShouldEqual, body["scannedAt"])
		So(res.PaidPrice.Decimal.IntPart(), ShouldEqual, body["paidPrice"])
		So(res.HasPaperBook, ShouldEqual, body["hasPaperBook"])
	})
}
