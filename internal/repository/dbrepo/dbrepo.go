package dbrepo

import (
	"database/sql"

	"github.com/wtran29/go-rss/internal/repository"
)

type postgresDBRepo struct {
	DB *sql.DB
}

func NewPostgresRepo(Conn *sql.DB) repository.DatabaseRepo {
	return &postgresDBRepo{
		DB: Conn,
	}
}
