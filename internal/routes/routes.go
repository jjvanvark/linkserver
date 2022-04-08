package routes

import (
	"github.com/go-chi/jwtauth"
	"github.com/go-playground/validator/v10"
	"github.com/jjvanvark/lor-deluxe/internal/models"
)

var validate *validator.Validate

type Routes struct {
	db   models.DataInterface
	auth *jwtauth.JWTAuth
}

func InitRoutes(db models.DataInterface, auth *jwtauth.JWTAuth) *Routes {
	validate = validator.New()

	return &Routes{db, auth}
}
