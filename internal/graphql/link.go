package graphql

import (
	"fmt"

	"github.com/jjvanvark/lor-deluxe/internal/models"
)

type linkResolver struct {
	link models.Link
}

func (r *linkResolver) Cursor() (string, error) {
	return r.link.Cursor, nil
}

func (r *linkResolver) Title() (string, error) {
	return r.link.Article.Title, nil
}

func (r *linkResolver) Url() (string, error) {
	return r.link.Url, nil
}

func (r *linkResolver) Byline() (string, error) {
	return r.link.Article.Byline, nil
}

func (r *linkResolver) Content() (string, error) {
	return r.link.Article.Content, nil
}

func (r *linkResolver) TextContent() (string, error) {
	return r.link.Article.TextContent, nil
}

func (r *linkResolver) Excerpt() (string, error) {
	return r.link.Article.Excerpt, nil
}

func (r *linkResolver) SiteName() (string, error) {
	return r.link.Article.SiteName, nil
}

func (r *linkResolver) Image() (*string, error) {
	return r.link.Article.Image, nil
}

func (r *linkResolver) Favicon() (*string, error) {
	return r.link.Article.Favicon, nil
}

func (r *linkResolver) Archive() ([]int32, error) {

	var index int
	var archive []int32 = make([]int32, len(r.link.Archive))
	var value int64

	for index, value = range r.link.Archive {
		archive[index] = int32(value)
	}

	fmt.Printf("archive :: %v :: %v\n", archive, r.link.Archive)

	return archive, nil
}

type linkEdgeResolver struct {
	link models.Link
}

func (r *linkEdgeResolver) Cursor() (string, error) {
	return r.link.Cursor, nil
}

func (r *linkEdgeResolver) Node() (*linkResolver, error) {
	return &linkResolver{r.link}, nil
}
