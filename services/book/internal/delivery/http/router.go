package http

import (
	"microservices-go/services/book/internal/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type router struct {
	book BookHandler
}

func NewRouter(bookService service.BookService) *router {
	return &router{book: *NewHandler(bookService)}
}

func (r *router) GetRoutes() http.Handler {

	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/v1/books", r.book.CreateBookHandler)
	router.HandlerFunc(http.MethodGet, "/v1/books/:id", r.book.ShowBookHandler)
	router.HandlerFunc(http.MethodGet, "/v1/books", r.book.ListBooksHandler)
	
	return router
}
