package http

import (
	"errors"
	"fmt"
	"microservices-go/pkg/request"
	"microservices-go/services/order/internal/repository"
	"microservices-go/services/order/internal/service"
	"net/http"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: service}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var dto service.CreateOrderDTO

	if err := request.ReadJSON(w, r, &dto); err != nil {
		request.BadRequestResponse(w, r, err)
		return
	}

	input := service.CreateOrderDTO{
		BookID: dto.BookID,
		Email:  dto.Email,
	}

	err := h.orderService.Create(r.Context(), input)
	if err != nil {
		request.ServerErrorResponse(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (h *OrderHandler) ShowOrder(w http.ResponseWriter, r *http.Request) {
	//email, err := request.ReadEmailParam(r)
	//if err != nil {
	//	request.ServerErrorResponse(w, r, err)
	//	return
	//}

	var dto service.CreateOrderDTO
	fmt.Println("ShowOrder")
	if err := request.ReadJSON(w, r, &dto); err != nil {
		request.BadRequestResponse(w, r, err)
		fmt.Println(dto)
		return
	}
	fmt.Println(dto)
	orders, err := h.orderService.Show(r.Context(), dto.Email)
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
	err = request.WriteJSON(w, http.StatusOK, map[string]any{"orders": orders}, nil)
	if err != nil {
		request.ServerErrorResponse(w, r, err)
		return
	}
}
