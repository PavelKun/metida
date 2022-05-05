package models

import (
	"time"
)

// ToDO: тут все должно быть DTO?
type Blog struct {
	ID          int32
	Name        string
	Description string
	Img         []byte
}

type Content struct {
	ID          int32
	UserEmail   string
	Name        string
	Description string
}

type Session struct {
	ID           int32
	UserEmail    string
	RefreshToken string
	AccessToken  string
	UserAgent    string
	IP           string
	ExpiresIn    int64
	CreatedAt    time.Time
}

type User struct {
	ID        int32
	Name      string
	Password  string
	Email     string
	Role      string
	IsDeleted bool
}
