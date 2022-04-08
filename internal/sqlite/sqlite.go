package sqlite

import (
	"github.com/go-playground/validator/v10"
	"github.com/jjvanvark/lor-deluxe/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type db struct {
	d *sqlx.DB
}

var validate *validator.Validate

func Init(filename string) (models.DataInterface, error) {
	var err error
	var database *sqlx.DB

	validate = validator.New()

	// connect the db
	if database, err = sqlx.Connect("sqlite3", filename); err != nil {
		return nil, err
	}

	// check db
	if err = database.Ping(); err != nil {
		defer database.Close()
		return nil, err
	}

	// run schema
	if _, err = database.Exec(dbSchema); err != nil {
		defer database.Close()
		return nil, err
	}

	return &db{d: database}, nil
}

func (d *db) Close() error {

	var err error

	if d.d == nil {
		return nil
	}

	if err = d.d.Close(); err != nil {
		return err
	}

	return nil
}
