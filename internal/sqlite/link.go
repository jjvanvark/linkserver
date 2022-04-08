package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jjvanvark/lor-deluxe/internal/models"
	"github.com/jmoiron/sqlx"
)

func (d *db) AddLink(
	userId int64,
	url string,
	title string,
	byline string,
	content string,
	textContent string,
	length int,
	excerpt string,
	siteName string,
	getImageId func() (*string, error),
	getFaviconId func() (*string, error),
) (*models.Link, error) {

	var err error
	var tx *sqlx.Tx
	var ctx context.Context
	var cancel context.CancelFunc
	var link *models.Link
	var imageId *string
	var faviconId *string

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	if tx, err = d.d.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: false}); err != nil {
		return nil, err
	}

	if link, err = d.getLinkByUrl(url, userId); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		if imageId, err = getImageId(); err != nil {
			fmt.Printf("sqlite :: link :: getImageId :: %v\n", err)
		}

		if faviconId, err = getFaviconId(); err != nil {
			fmt.Printf("sqlite :: link :: getFaviconId :: %v\n", err)
		}

		if err = addNewLink(
			tx,
			userId,
			url,
			title,
			byline,
			content,
			textContent,
			length,
			excerpt,
			siteName,
			imageId,
			faviconId,
		); err != nil {
			return nil, err
		}
	} else {
		// if link.Article != nil && link.Article.Title == title && link.Article.Content == content {
		// 	if _, err = tx.Exec(`
		// 		UPDATE link
		// 		SET updated_at = ?
		// 		WHERE deleted_at IS NULL
		// 		AND id = ?
		// 		AND user_id = ?;
		// 	`, time.Now(), link.ID, userId); err != nil {
		// 		return nil, err
		// 	}
		// } else {
		if err = addArticleOnly(
			tx,
			link.ID,
			title,
			byline,
			content,
			textContent,
			length,
			excerpt,
			siteName,
			getImageId,
			getFaviconId,
		); err != nil {
			return nil, err
		}
		// }

	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil

}

func (d *db) UpdateLink(cursor string, userId int64, field string, value interface{}) (*models.Link, error) {

	var err error
	var ok bool
	var timeValue time.Time

	if err = validate.Var(field, "required,oneof=updated_at"); err != nil {
		return nil, err
	}

	switch field {
	case "updated_at":
		if timeValue, ok = value.(time.Time); !ok {
			return nil, errors.New("timevalue typecasting error")
		}

		if _, err = d.d.Exec(`
			UPDATE link
			SET updated_at = ?
			WHERE deleted_at IS NULL
			AND cursor = ?
			AND user_id = ?;
		`, timeValue, cursor, userId); err != nil {
			return nil, err
		}
	}

	return d.GetLink(cursor, userId)

}

func addNewLink(
	tx *sqlx.Tx,
	userId int64,
	url string,
	title string,
	byline string,
	content string,
	textContent string,
	length int,
	excerpt string,
	siteName string,
	image *string,
	favicon *string,
) error {

	var err error
	var now time.Time
	var result sql.Result
	var lastId int64

	now = time.Now()
	if result, err = tx.Exec(`
	INSERT INTO link
		(created_at,
		updated_at,
		cursor,
		user_id,
		url,
		host
		)
		VALUES (?, ?, ?, ?, ?, ?);
	`, now, now, uuid.New().String(), userId, url, "host"); err != nil {
		return err
	}

	if lastId, err = result.LastInsertId(); err != nil {
		return err
	}

	now = time.Now()
	if _, err = tx.Exec(`
	INSERT INTO article
		(created_at,
		updated_at,
		link_id,
		title,
		by_line,
		content,
		text_content,
		length,
		excerpt,
		site_name,
		image,
		favicon
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`, now, now, lastId, title, byline, content, textContent, length, excerpt, siteName, image, favicon); err != nil {
		return err
	}

	return nil

}

func addArticleOnly(
	tx *sqlx.Tx,
	linkId int64,
	title string,
	byline string,
	content string,
	textContent string,
	length int,
	excerpt string,
	siteName string,
	getImageId func() (*string, error),
	getFaviconId func() (*string, error),
) error {

	var err error
	var now time.Time
	var result sql.Result
	var imageId *string
	var faviconId *string

	if imageId, err = getImageId(); err != nil {
		fmt.Printf("sqlite :: link :: addArticle only :: getImageId :: %v\n", err)
	}

	if faviconId, err = getFaviconId(); err != nil {
		fmt.Printf("sqlite :: link :: addArticle only :: getFaviconId :: %v\n", err)
	}

	now = time.Now()
	if result, err = tx.Exec(`
	INSERT INTO article
		(created_at,
		updated_at,
		link_id,
		title,
		by_line,
		content,
		text_content,
		length,
		excerpt,
		site_name,
		image,
		favicon
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`, now, now, linkId, title, byline, content, textContent, length, excerpt, siteName, imageId, faviconId); err != nil {
		return err
	}

	if _, err = result.LastInsertId(); err != nil {
		return err
	}

	return nil

}

func (d *db) getLinkByUrl(url string, userId int64) (*models.Link, error) {

	var err error
	var cursor string

	if err = d.d.Get(&cursor, `
		SELECT cursor
		FROM link
		WHERE deleted_at IS NULL
		AND url = ?
		AND user_id = ?;
	`, url, userId); err != nil {
		return nil, err
	}

	return d.GetLink(cursor, userId)
}

func (d *db) GetLink(cursor string, userId int64) (*models.Link, error) {

	var err error
	var article models.Article
	var link models.Link
	var archive []int64

	if err = d.d.Get(&link, `
		SELECT id, created_at, updated_at, deleted_at, cursor, url, host FROM link
		WHERE deleted_at IS NULL
		AND cursor = ?
		AND user_id = ?;
	`, cursor, userId); err != nil {
		return nil, err
	}

	if article, archive, err = d.getArticleByLink(link.ID); err != nil {
		return nil, err
	}

	link.Article = article
	link.Archive = archive

	return &link, nil

}

func (d *db) getArticleByLink(id int64) (models.Article, []int64, error) {

	var err error
	var article models.Article
	var archive []int64

	if err = d.d.Get(&article, `
		SELECT id, created_at, updated_at, deleted_at, title, by_line, content, text_content, length, excerpt, site_name, image, favicon FROM article
		WHERE deleted_at IS NULL
		AND link_id = ?
		ORDER BY created_at DESC;
	`, id); err != nil {
		return article, nil, err
	}

	if err = d.d.Select(&archive, `
		SELECT id FROM article
		WHERE deleted_at IS NULL
		AND link_id = ?
		ORDER BY created_at DESC;
	`, id); err != nil {
		return article, nil, err
	}

	return article, archive, nil

}

func (d *db) TotalLinks(id int64) (int, error) {

	var err error
	var count int

	if err = d.d.Get(&count, `
		SELECT count(id)
		FROM link
		WHERE deleted_at IS NULL
		AND user_id = ?;
	`, id); err != nil {
		return 0, err
	}

	return count, nil

}

func (d *db) GetLinks(userId int64, amount int, isDesc bool, cursor *string) ([]models.Link, bool, bool, error) {

	var err error
	var timeStamp time.Time
	var result []models.Link
	var direction string
	var sign string
	var link models.Link
	var index int
	var count int = amount + 1
	var hasNext bool
	var hasPrevious bool

	if cursor == nil {

		if isDesc {
			direction = "DESC"
			sign = "<="
		} else {
			direction = "ASC"
			sign = ">="
		}

		if err = d.d.Get(&timeStamp, fmt.Sprintf(`
			SELECT updated_at
			FROM link
			WHERE deleted_at IS NULL
			AND user_id = ?
			ORDER BY updated_at %v
			LIMIT 1;
		`, direction), userId); err != nil {
			return nil, hasPrevious, hasNext, err
		}

		if err = d.d.Select(&result, fmt.Sprintf(`
		SELECT *
		FROM link
		WHERE deleted_at IS NULL
		AND user_id = ?
		AND updated_at %v ?
		ORDER BY updated_at %v
		LIMIT ?;
	`, sign, direction), userId, timeStamp, count); err != nil {
			return nil, hasPrevious, hasNext, err
		}

	} else {

		if isDesc {
			direction = "DESC"
			sign = "<"
		} else {
			direction = "ASC"
			sign = ">"
		}

		if err = d.d.Get(&timeStamp, `
			SELECT updated_at
			FROM link
			WHERE deleted_at IS NULL
			AND user_id = ?
			AND cursor = ?;
		`, userId, cursor); err != nil {
			return nil, hasPrevious, hasNext, err
		}

		hasPrevious = true

		if err = d.d.Select(&result, fmt.Sprintf(`
			SELECT *
			FROM link
			WHERE deleted_at IS NULL
			AND user_id = ?
			AND updated_at %v ?
			ORDER BY updated_at %v
			LIMIT ?;
		`, sign, direction), userId, timeStamp, count); err != nil {
			return nil, hasPrevious, hasNext, err
		}

	}

	if len(result) == count {
		hasNext = true
		result = result[:len(result)-1]
	}

	var archive []int64

	for index, link = range result {

		var articles []models.Article

		if err = d.d.Select(&articles, `
			SELECT *
			FROM article
			WHERE deleted_at IS NULL
			AND link_id = ?
			ORDER BY updated_at DESC;
		`, link.ID); err != nil {
			return nil, hasPrevious, hasNext, err
		}

		if len(articles) < 1 {
			return nil, hasPrevious, hasNext, errors.New("no article found by link")
		}

		result[index].Article = articles[0]

		archive = make([]int64, len(articles))
		for index, article := range articles {
			archive[index] = article.ID
		}

		result[index].Archive = archive
	}

	return result, hasPrevious, hasNext, nil

}
