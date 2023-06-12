package http

import (
	"errors"
	"microservices-go/pkg/request"
	"microservices-go/services/user/internal/service"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
}

func NewHandler(service service.UserService) *UserHandler {
	return &UserHandler{userService: service}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var dto service.UserSignUpDTO

	if err := request.ReadJSON(w, r, &dto); err != nil {
		request.BadRequestResponse(w, r, err)
		return
	}

	input := service.UserSignUpDTO{
		Name:         dto.Name,
		Email:        dto.Email,
		HashPassword: dto.HashPassword,
	}

	err := h.userService.SignUp(r.Context(), input)
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
	w.WriteHeader(http.StatusOK)

}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {

	var input service.UserSignInDTO

	if err := request.ReadJSON(w, r, &input); err != nil {
		request.BadRequestResponse(w, r, err)
		return
	}

	//input := service.UserSignInDTO{
	//	Email:        dto.Email,
	//	HashPassword: dto.HashPassword,
	//}

	token, err := h.userService.SignIn(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrWrongCredentials):
			request.NotFoundResponse(w, r)
			return
		default:
			request.ServerErrorResponse(w, r, err)
			return
		}
	}
	request.WriteJSON(w, http.StatusOK, map[string]any{"token": token.PlainText}, nil)
}
