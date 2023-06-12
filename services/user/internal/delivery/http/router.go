package http

import (
	"microservices-go/services/user/internal/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type router struct {
	user UserHandler
}

func NewRouter(userService service.UserService) *router {
	return &router{user: *NewHandler(userService)}
}

func (r *router) GetRoutes() http.Handler {

	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/v1/user/signup", r.user.RegisterUser)
	router.HandlerFunc(http.MethodPost, "/v1/user/signin", r.user.LoginUser)

	return router
}
