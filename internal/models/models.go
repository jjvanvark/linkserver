package models

import "time"

type Model struct {
	ID        int64      `json:"id" validate:"-" db:"id"`
	CreatedAt time.Time  `json:"created_at" validate:"-" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" validate:"-" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" validate:"-" db:"deleted_at"`
}

type User struct {
	Model
	Email    string `json:"email" validate:"required,email" db:"email"`
	Name     string `json:"name" validate:"required" db:"name"`
	Password []byte `json:"-" validate:"-" db:"password"`
	Key      string `json:"key" validate:"required" db:"key"`
}

type Asset struct {
	Model
	Cursor string `json:"cursor" validate:"required,uuid" db:"cursor"`
	File   []byte `json:"-" validate:"-" db:"file"`
}

type Link struct {
	Model
	UserId  int64   `json:"user_id" validate:"required" db:"user_id"`
	Cursor  string  `json:"cursor" validate:"required,uuid" db:"cursor"`
	Url     string  `json:"url" validate:"requred,url" db:"url"`
	Host    string  `json:"host" validate:"requred" db:"host"`
	Article Article `json:"article"`
	Archive []int64 `json:"articles"`
}

type Article struct {
	Model
	Title       string  `json:"title" db:"title"`
	Byline      string  `json:"by_line" db:"by_line"`
	Content     string  `json:"content" db:"content"`
	TextContent string  `json:"text_content" db:"text_content"`
	Length      int     `json:"length" db:"length"`
	Excerpt     string  `json:"excerpt" db:"excerpt"`
	SiteName    string  `json:"site_name" db:"site_name"`
	Image       *string `json:"image" db:"image"`
	Favicon     *string `json:"favicon" db:"favicon"`
	LinkId      *int64  `json:"-" db:"link_id"`
}
