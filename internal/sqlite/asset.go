package sqlite

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jjvanvark/lor-deluxe/internal/models"
)

func (d *db) AddAsset(file []byte) (*string, error) {

	var err error
	var result sql.Result
	var now time.Time
	var id int64
	var cursor *string

	now = time.Now()

	if result, err = d.d.Exec(`
		INSERT INTO asset
		(created_at, updated_at, cursor, file)
		VALUES (?, ?, ?, ?);
	`, now, now, uuid.New().String(), file); err != nil {
		return nil, err
	}

	if id, err = result.LastInsertId(); err != nil {
		return nil, err
	}

	if err = d.d.Get(cursor, `
		SELECT cursor FROM asset
		WHERE deleted_at IS NULL
		AND id = ?;
	`, id); err != nil {
		return nil, err
	}

	return cursor, nil

}

func (d *db) GetAsset(cursor string) (*models.Asset, error) {

	var err error
	var asset models.Asset

	if err = d.d.Get(&asset, `
		SELECT * FROM asset
		WHERE deleted_at IS NULL
		AND cursor = ?;
		`, cursor); err != nil {
		return nil, err
	}

	return &asset, nil

}

func (d *db) RemoveAsset(cursor string) error {

	var err error

	if _, err = d.d.Exec(`
		UPDATE asset
		SET deleted_at = ?
		WHERE deleted_at IS NULL
		AND cursor = ?;
	`, time.Now(), cursor); err != nil {
		return err
	}

	return nil

}
