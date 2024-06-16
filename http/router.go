package http

import (
	"crypto/subtle"
	"invoicehub"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct {
	e *echo.Echo
}

func NewRouter(
	companies invoicehub.CompanyRepository,
	invoices invoicehub.InvoiceService,
) *Router {
	ans := Router{
		e: echo.New(),
	}

	ans.e.Debug = true

	username := os.Getenv("FH_USERNAME")
	password := os.Getenv("FH_PASSWORD")

	if username == "" || password == "" {
		panic("username or password not set")
	}

	// middlewares
	ans.e.Use(middleware.BasicAuth(func(u, p string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(u), []byte(username)) == 1 &&
			subtle.ConstantTimeCompare([]byte(p), []byte(password)) == 1 {
			return true, nil
		}
		return false, nil
	}))

	baseHandler := baseHandler{
		e:         ans.e,
		companies: companies,
		invoices:  invoices,
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
	e         *echo.Echo
	companies invoicehub.CompanyRepository
	invoices  invoicehub.InvoiceService
}

func (b *baseHandler) RegisterRoutes() {
	b.e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "THANK YOU")
	})
}
