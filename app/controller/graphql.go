package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/librarios/go-librarios/app/model"
	"github.com/librarios/go-librarios/app/service"
)

var bookType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Book",
		Fields: graphql.Fields{
			"id":            &graphql.Field{Type: graphql.Int},
			"isbn13":        &graphql.Field{Type: graphql.String},
			"isbn10":        &graphql.Field{Type: graphql.String},
			"title":         &graphql.Field{Type: graphql.String},
			"originalIsbn":  &graphql.Field{Type: graphql.String},
			"originalTitle": &graphql.Field{Type: graphql.String},
			"contents":      &graphql.Field{Type: graphql.String},
			"url":           &graphql.Field{Type: graphql.String},
			"pubDate":       &graphql.Field{Type: graphql.String},
			"authors":       &graphql.Field{Type: graphql.String},
			"translators":   &graphql.Field{Type: graphql.String},
			"publisher":     &graphql.Field{Type: graphql.String},
			"price":         &graphql.Field{Type: graphql.Float},
			"currency":      &graphql.Field{Type: graphql.String},
		},
	},
)

var ownedBookType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "OwnedBook",
		Fields: graphql.Fields{
			"id":           &graphql.Field{Type: graphql.Int},
			"isbn":         &graphql.Field{Type: graphql.String},
			"owner":        &graphql.Field{Type: graphql.String},
			"acquiredAt":   &graphql.Field{Type: graphql.String},
			"scanneddAt":   &graphql.Field{Type: graphql.String},
			"paidPrice":    &graphql.Field{Type: graphql.Float},
			"actualPages":  &graphql.Field{Type: graphql.Int},
			"hasPaperBook": &graphql.Field{Type: graphql.Boolean},
		},
	},
)

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"interestBookList": &graphql.Field{
			Type:        graphql.NewList(graphql.String),
			Description: "get interestBook ISBN list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// TODO: not implemented
				return []string{}, nil
			},
		},
		"ownedBookList": &graphql.Field{
			Type:        graphql.NewList(ownedBookType),
			Description: "get ownedBook list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// TODO: not implemented
				return []model.OwnedBook{}, nil
			},
		},
	},
})

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"addOwnedBook": &graphql.Field{
			Type:        ownedBookType,
			Description: "Add owned book",
			Args: graphql.FieldConfigArgument{
				"isbn": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return service.BookService.AddOwnedBook(service.AddOwnedBook{
					Isbn: params.Args["isbn"].(string),
				})
			},
		},
		"updateBook": &graphql.Field{
			Type:        bookType,
			Description: "update book",
			Args: graphql.FieldConfigArgument{
				"id":            &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"isbn13":        &graphql.ArgumentConfig{Type: graphql.String},
				"isbn10":        &graphql.ArgumentConfig{Type: graphql.String},
				"title":         &graphql.ArgumentConfig{Type: graphql.String},
				"originalIsbn":  &graphql.ArgumentConfig{Type: graphql.String},
				"originalTitle": &graphql.ArgumentConfig{Type: graphql.String},
				"contents":      &graphql.ArgumentConfig{Type: graphql.String},
				"url":           &graphql.ArgumentConfig{Type: graphql.String},
				"pubDate":       &graphql.ArgumentConfig{Type: graphql.String},
				"authors":       &graphql.ArgumentConfig{Type: graphql.String},
				"translators":   &graphql.ArgumentConfig{Type: graphql.String},
				"publisher":     &graphql.ArgumentConfig{Type: graphql.String},
				"price":         &graphql.ArgumentConfig{Type: graphql.Float},
				"currency":      &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(uint)
				update := params.Args
				return service.BookService.UpdateBook(id, update)
			},
		},
		"updateOwnedBook": &graphql.Field{
			Type:        ownedBookType,
			Description: "update owned book",
			Args: graphql.FieldConfigArgument{
				"id":           &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"isbn":         &graphql.ArgumentConfig{Type: graphql.String},
				"acquiredAt":   &graphql.ArgumentConfig{Type: graphql.String},
				"scanneddAt":   &graphql.ArgumentConfig{Type: graphql.String},
				"paidPrice":    &graphql.ArgumentConfig{Type: graphql.Float},
				"actualPages":  &graphql.ArgumentConfig{Type: graphql.Int},
				"hasPaperBook": &graphql.ArgumentConfig{Type: graphql.Boolean},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(uint)
				update := params.Args
				return service.BookService.UpdateOwnedBook(id, update)
			},
		},
	},
})

var subscriptionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootSubscription",
	Fields: graphql.Fields{
	},
})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:        queryType,
		Mutation:     mutationType,
		Subscription: subscriptionType,
	},
)

// GraphqlHandler returns handler function to handle GraphQL query.
func GraphqlHandler() gin.HandlerFunc {
	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
