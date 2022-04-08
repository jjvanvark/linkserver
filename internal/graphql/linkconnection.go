package graphql

import (
	"github.com/jjvanvark/lor-deluxe/internal/models"
)

type linkConnectionResolver struct {
	total       int
	links       []models.Link
	hasPrevious bool
	hasNext     bool
}

func (r *linkConnectionResolver) Total() (int32, error) {
	return int32(r.total), nil
}

func (r *linkConnectionResolver) Nodes() ([]*linkResolver, error) {

	var result []*linkResolver = make([]*linkResolver, len(r.links))

	for index, link := range r.links {
		result[index] = &linkResolver{link}
	}

	return result, nil

}

func (r *linkConnectionResolver) Edges() ([]*linkEdgeResolver, error) {

	var result []*linkEdgeResolver = make([]*linkEdgeResolver, len(r.links))

	for index, link := range r.links {
		result[index] = &linkEdgeResolver{link}
	}

	return result, nil

}

func (r *linkConnectionResolver) PageInfo() (*pageInfoResolver, error) {
	var startCursor *string
	var endCursor *string

	if len(r.links) > 0 {
		startCursor = &r.links[0].Cursor
		endCursor = &r.links[len(r.links)-1].Cursor
	}

	return &pageInfoResolver{r.hasPrevious, r.hasNext, startCursor, endCursor}, nil
}
