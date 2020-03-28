package controller

import (
	"github.com/librarios/go-librarios/app/model"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestBookSearchIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	Convey("search by isbn", t, func() {
		var result struct {
			Data []*model.Book
		}
		rw := testServer.Get("/book/search?isbn=9788960778320", &result)

		book := result.Data[0]
		So(rw.Code, ShouldEqual, http.StatusOK)
		So(book.Title, ShouldContainSubstring, "The Go Programming Language")
	})
}
