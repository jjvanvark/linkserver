package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jjvanvark/lor-deluxe/internal/models"
	"github.com/jjvanvark/lor-deluxe/internal/security"
)

func (r *Routes) HandleHttpLogin(rw http.ResponseWriter, req *http.Request) {

	var email string
	var password string

	email = chi.URLParam(req, "email")
	password = chi.URLParam(req, "password")

	r.handleLogin(rw, email, password)

}

func (r *Routes) handleLogin(rw http.ResponseWriter, email string, password string) {

	var err error
	var user *models.User
	var token string

	if user, err = r.db.GetUserByEmail(email); err != nil {
		fmt.Printf("Handlelogin GetUserByEmail error :: %v\n", err)
		http.Error(rw, "bad request", http.StatusBadRequest)
		return
	}

	if !security.CheckPasswordHash(password, user.Password) {
		fmt.Printf("Handlelogin Password incorrect :: %v\n", user)
		http.Error(rw, "bad request", http.StatusBadRequest)
		return
	}

	if token, err = security.GetToken(r.auth, user.ID); err != nil {
		fmt.Printf("Handlelogin Get Token :: %v\n", user)
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		MaxAge:   3600,
		HttpOnly: true,
		Path:     "/",
	})

}

func (r *Routes) HandleLogin(rw http.ResponseWriter, req *http.Request) {

	var err error
	var result struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err = json.NewDecoder(req.Body).Decode(&result); err != nil {
		fmt.Printf("handlelogin body decoder :: %v\n", err)
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}

	if err = validate.Struct(result); err != nil {
		fmt.Printf("handlelogin body struct validate error :: %v\n", err)
		http.Error(rw, "bad request", http.StatusBadRequest)
		return
	}

	r.handleLogin(rw, result.Email, result.Password)

}
