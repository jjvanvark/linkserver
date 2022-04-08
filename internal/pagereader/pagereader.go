package pagereader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/jjvanvark/lor-deluxe/internal/models"
)

func AddPage(db models.DataInterface, url *url.URL, userId int64) {

	var err error
	var article readability.Article
	var file []byte

	if article, err = readability.FromURL(url.String(), 30*time.Second); err != nil {
		fmt.Printf("Readability FromUrl error :: %v\n", err)
		return
	}

	if _, err = db.AddLink(userId,
		url.String(),
		article.Title,
		article.Byline,
		article.Content,
		article.TextContent,
		article.Length,
		article.Excerpt,
		article.SiteName,
		func() (*string, error) {
			if file, err = downloadFile(article.Image); err != nil {
				return nil, err
			}

			return db.AddAsset(file)
		},
		func() (*string, error) {
			if file, err = downloadFile(article.Favicon); err != nil {
				return nil, err
			}

			return db.AddAsset(file)
		},
	); err != nil {
		fmt.Printf("Readability Addlink error :: %v\n", err)
		return
	}

}

func downloadFile(url string) ([]byte, error) {

	var err error
	var response *http.Response
	var result []byte

	if response, err = http.Get(url); err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if result, err = io.ReadAll(response.Body); err != nil {
		return nil, err
	}

	return result, nil

}
