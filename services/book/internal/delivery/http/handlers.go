package http

import (
	"errors"
	"fmt"
	"microservices-go/pkg/request"
	"microservices-go/pkg/validator"
	"microservices-go/services/book/internal/repository"
	"microservices-go/services/book/internal/service"
	"net/http"
)

type BookHandler struct {
	bookService service.BookService
}

func NewHandler(service service.BookService) *BookHandler {
	return &BookHandler{bookService: service}
}

func (h *BookHandler) CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var input service.CreateBookDTO

	err := request.ReadJSON(w, r, &input)
	if err != nil {
		request.BadRequestResponse(w, r, err)
		return
	}

	book := service.CreateBookDTO{
		Title:  input.Title,
		Year:   input.Year,
		Genres: input.Genres,
	}

	err = h.bookService.CreateBook(r.Context(), book)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrFailedValidation):
			request.BadRequestResponse(w, r, err)
			return
		case errors.Is(err, service.ErrDuplicate):
			request.RecordDuplicationResponse(w, r)
			return
		default:
			request.ServerErrorResponse(w, r, err)
			return
		}
	}

	err = request.WriteJSON(w, http.StatusCreated, map[string]any{"book": book}, nil)
	if err != nil {
		request.ServerErrorResponse(w, r, err)
		return
	}

}

func (h *BookHandler) ShowBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := request.ReadIDParam(r)
	if err != nil {
		request.NotFoundResponse(w, r)
		return
	}

	book, err := h.bookService.GetBookByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			request.NotFoundResponse(w, r)
			return
		default:
			request.ServerErrorResponse(w, r, err)
			return
		}
	}
	err = request.WriteJSON(w, http.StatusOK, map[string]any{"book": book}, nil)
	if err != nil {
		request.ServerErrorResponse(w, r, err)
		return
	}
}

func (h *BookHandler) ListBooksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		repository.Filters
	}

	fmt.Println("ListBooksHandler")
	v := validator.New()
	qs := r.URL.Query()
	input.Title = request.ReadString(qs, "title", "")
	input.Genres = request.ReadCSV(qs, "genres", []string{})
	input.Filters.Page = request.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = request.ReadInt(qs, "page_size", 20, v)
	input.Filters.Sort = request.ReadString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "year", "-id", "-title", "-year"}

	books, err := h.bookService.GetBooks(r.Context(), input.Title, input.Genres, input.Filters)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			request.NotFoundResponse(w, r)
			return
		default:
			request.ServerErrorResponse(w, r, err)
			return
		}
	}
	err = request.WriteJSON(w, http.StatusOK, map[string]any{"books": books}, nil)
	if err != nil {
		request.ServerErrorResponse(w, r, err)
		return
	}
}
