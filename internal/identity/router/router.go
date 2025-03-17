package router

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/maketaio/apiserver/internal/types"
	"github.com/maketaio/apiserver/pkg/api"
)

type Backend interface {
	SignUp(ctx context.Context, input *types.SignUpInput) (*api.User, error)
}

type Router struct {
	backend Backend
}

func New(backend Backend) *Router {
	return &Router{backend: backend}
}

func (r *Router) Attach(g *echo.Group) {
	g.POST("/signup", r.signUp)
}

func (r *Router) signUp(c echo.Context) error {
	type payload struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
	}

	var p payload
	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	input := &types.SignUpInput{
		Email:     p.Email,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Password:  p.Password,
	}

	if _, err := r.backend.SignUp(c.Request().Context(), input); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, nil)
}
