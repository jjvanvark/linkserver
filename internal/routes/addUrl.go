package routes

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jjvanvark/lor-deluxe/internal/pagereader"
	"github.com/jjvanvark/lor-deluxe/internal/security"
)

func (r *Routes) HandleAddUrl(rw http.ResponseWriter, req *http.Request) {

	var err error
	var urlString string
	var key string
	var userId int64
	var urlValue *url.URL

	key = req.URL.Query().Get("key")
	urlString = req.URL.Query().Get("url")

	if err = validate.Var(urlString, "required,url"); err != nil {
		fmt.Printf("Add url :: Validate url :: %v\n", err)
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}

	if urlValue, err = url.Parse(urlString); err != nil {
		fmt.Printf("Add url :: Parse url :: %v\n", err)
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}

	if err = validate.Var(key, "required"); err != nil {
		fmt.Printf("Add url :: Validate key :: %v\n", err)
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}

	if userId, err = security.ParseUserToken(key); err != nil {
		fmt.Printf("Add url :: Parse user token :: %v\n", err)
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}

	go pagereader.AddPage(r.db, urlValue, userId)

	rw.Write([]byte(fmt.Sprintf("User id: %v\n", userId)))

}
