package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/jjvanvark/lor-deluxe/internal/models"
	"github.com/jjvanvark/lor-deluxe/internal/security"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

type resolver struct {
	db models.DataInterface
}

func InitGraphql(db models.DataInterface) (http.Handler, error) {

	var err error
	var schema *graphql.Schema

	if schema, err = graphql.ParseSchema(graphqlSchema, &resolver{db}); err != nil {
		return nil, err
	}

	return &relay.Handler{Schema: schema}, nil

}

// query

func (r *resolver) Me(ctx context.Context) (*userResolver, error) {

	var err error
	var id int64
	var user *models.User

	if id, err = security.GetIdFromContext(ctx); err != nil {
		fmt.Printf("Graphql :: Me :: %v\n", err)
		return nil, err
	}

	if user, err = r.db.GetUser(id); err != nil {
		fmt.Printf("Graphql :: Me :: %v\n", err)
		return nil, err
	}

	return &userResolver{user}, nil

}

func (r *resolver) Link(ctx context.Context, args struct {
	Cursor string
}) (*linkResolver, error) {

	var err error
	var id int64
	var link *models.Link

	if id, err = security.GetIdFromContext(ctx); err != nil {
		fmt.Printf("Graphql :: Link :: %v\n", err)
		return nil, err
	}

	if link, err = r.db.GetLink(args.Cursor, id); err != nil {
		fmt.Printf("Graphql :: Link :: GetLink %v\n", err)
		return nil, err
	}

	return &linkResolver{*link}, nil

}

func (r *resolver) Links(ctx context.Context, args struct {
	First *int32
	Last  *int32
	After *string
}) (*linkConnectionResolver, error) {

	var err error
	var amount int
	var order bool
	var id int64
	var total int
	var hasPrevious bool
	var hasNext bool

	if amount, order, err = parseConnectionParams(args.First, args.Last); err != nil {
		return nil, err
	}

	if id, err = security.GetIdFromContext(ctx); err != nil {
		fmt.Printf("Graphql :: Link :: %v\n", err)
		return nil, err
	}

	if total, err = r.db.TotalLinks(id); err != nil {
		return nil, err
	}

	var result []models.Link
	if result, hasPrevious, hasNext, err = r.db.GetLinks(id, amount, order, args.After); err != nil {
		if err == sql.ErrNoRows {
			result = make([]models.Link, 0)
		} else {
			return nil, err
		}
	}

	return &linkConnectionResolver{total, result, hasPrevious, hasNext}, nil
}

// helpers

func parseConnectionParams(first *int32, last *int32) (int, bool, error) {
	if first == nil && last == nil {
		return 0, true, errors.New("either first or last should have a value")
	}

	if first != nil {
		return int(*first), true, nil
	} else {
		return int(*last), false, nil
	}

}

// func idToInt(id graphql.ID) (int64, error) {

// 	var err error
// 	var result int64

// 	if result, err = strconv.ParseInt(string(id), 10, 64); err != nil {
// 		fmt.Printf("Graphql :: idToInt :: %v\n", err)
// 		return 0, err
// 	}

// 	return result, nil

// }

func intToId(id int64) graphql.ID {
	return graphql.ID(fmt.Sprint(id))
}
