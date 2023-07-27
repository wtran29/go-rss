package repository

import "github.com/wtran29/go-rss/internal/models"

type DatabaseRepo interface {
	CreateUser(u models.User) (models.User, error)
	GetUserByAPIKey(apiKey string) (models.User, error)
}
