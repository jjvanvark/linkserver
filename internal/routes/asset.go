package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jjvanvark/lor-deluxe/internal/models"
)

func (r *Routes) HandleAsset(rw http.ResponseWriter, req *http.Request) {

	var err error
	var cursor string
	var asset *models.Asset

	cursor = chi.URLParam(req, "cursor")

	if err = validate.Var(cursor, "required,uuid"); err != nil {
		fmt.Printf("Routes :: HandleAsset :: validate cursor %v\n", err)
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}

	if asset, err = r.db.GetAsset(cursor); err != nil {
		fmt.Printf("Routes :: HandleAsset :: Get asset %v\n", err)
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}

	rw.Write(asset.File)

}
