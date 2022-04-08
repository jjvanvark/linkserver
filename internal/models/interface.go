package models

import (
	"fmt"
	"time"
)

type DataInterface interface {
	ddUser(
		email string,
		name string,
		password []byte,
		getKey func(id int64) (string, error),
	) (*User, error)
	GetUser(id int64) (*User, error)
	GetUserByEmail(email string) (*User, error)
	Close() error
	AddAsset(file []byte) (*string, error)
	RemoveAsset(cursor string) error
	GetAsset(cursor string) (*Asset, error)
	AddLink(
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
	) (*Link, error)
	GetLink(cursor string, userId int64) (*Link, error)
	UpdateLink(cursor string, userId int64, field string, value interface{}) (*Link, error)
	TotalLinks(userId int64) (int, error)
	GetLinks(userId int64, amount int, isDesc bool, cursor *string) ([]Link, bool, bool, error)
}

func hello() {
	fmt.Println("now", time.Now())
}
