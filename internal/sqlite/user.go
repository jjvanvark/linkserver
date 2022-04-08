package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jjvanvark/lor-deluxe/internal/models"
	"github.com/jmoiron/sqlx"
)

func (d *db) AddUser(
	email string,
	name string,
	password []byte,
	getKey func(id int64) (string, error),
) (*models.User, error) {

	var err error
	var tx *sqlx.Tx
	var now time.Time
	var ctx context.Context
	var cancel context.CancelFunc
	var result sql.Result
	var id int64
	var key string

	// transform
	email = strings.ToLower(email)

	// validate
	if err = validate.Var(email, "required,email"); err != nil {
		return nil, err
	}

	if err = validate.Var(name, "required"); err != nil {
		return nil, err
	}

	if len(password) == 0 {
		return nil, errors.New("invalid password blob")
	}

	// insert user
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	if tx, err = d.d.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: false}); err != nil {
		return nil, err
	}

	now = time.Now()

	if result, err = tx.Exec(`
		INSERT INTO user
		(created_at, updated_at, email, name, password)
		VALUES (?, ?, ?, ?, ?);
	`, now, now, email, name, password); err != nil {
		return nil, err
	}

	if id, err = result.LastInsertId(); err != nil {
		return nil, err
	}

	if key, err = getKey(id); err != nil {
		cancel()
		return nil, err
	}

	// TODO: validate key?

	if _, err = tx.Exec(`
		UPDATE user
		SET key = ?
		WHERE deleted_at IS NULL
		AND id = ?;
	`, key, id); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil

}

func (d *db) GetUser(id int64) (*models.User, error) {

	var err error
	var user models.User

	// query
	if err = d.d.Get(&user, `
		SELECT * FROM user
		WHERE deleted_at IS NULL
		AND id = ?;
	`, id); err != nil {
		return nil, err
	}

	return &user, nil

}

func (d *db) GetUserByEmail(email string) (*models.User, error) {

	var err error
	var user models.User

	// transform
	email = strings.ToLower(email)

	// validation
	if err = validate.Var(email, "required,email"); err != nil {
		return nil, err
	}

	// query
	if err = d.d.Get(&user, `
		SELECT * FROM user
		WHERE deleted_at IS NULL
		AND email = ?;
	`, email); err != nil {
		return nil, err
	}

	return &user, nil

}
