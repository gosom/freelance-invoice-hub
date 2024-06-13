package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Router struct {
	e *echo.Echo
}

func NewRouter() *Router {
	ans := Router{
		e: echo.New(),
	}

	ans.e.Debug = true

	baseHandler := baseHandler{
		e: ans.e,
	}

	handlers := []handler{
		&baseHandler,
	}

	for _, h := range handlers {
		h.RegisterRoutes()
	}

	return &ans
}

func (r *Router) Handler() http.Handler {
	return r.e
}

type handler interface {
	RegisterRoutes()
}

type baseHandler struct {
	e *echo.Echo
}

func (b *baseHandler) RegisterRoutes() {
	b.e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "HELLO WORLD")
	})
}
