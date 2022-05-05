package repositories

import "github.com/Dsmit05/metida/internal/models"

type StorageI interface {
	Close()
	CreateUser(UserName, Password, Email, Role string) error
	ReadUser(Email string) (models.User, error)
	UpdatePassword(Password, Email string) error
	DeleteUser(Email string) error
	CreateSession(Email, RefreshToken string) error
	UpdateSession(Email, RefreshToken string) error
	CreateUserAndSession(UserName, Password, Email, Role, RefreshToken string) error
	GetEmailRoleFromSession(RefreshToken string) (Email, Role string, err error)
}
