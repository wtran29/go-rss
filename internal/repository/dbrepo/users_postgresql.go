package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/wtran29/go-rss/internal/models"
)

func (m *postgresDBRepo) CreateUser(u models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO users (created_at, updated_at, name, api_key)
				VALUES ($1, $2, $3, encode(sha256(random()::text::bytea), 'hex'))
				RETURNING id`

	var newId int
	err := m.DB.QueryRowContext(ctx, stmt,
		time.Now(),
		time.Now(),
		u.Name,
	).Scan(&newId)
	if err != nil {
		log.Println(err)
	}

	log.Println(newId)

	// userID, err := res.LastInsertId()
	// if err != nil {
	// 	log.Println(err)
	// }

	query := `SELECT id, name, created_at, updated_at, api_key FROM users WHERE id = $1`

	var newUser models.User

	row := m.DB.QueryRowContext(ctx, query, newId)

	err = row.Scan(
		&newUser.ID,
		&newUser.Name,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
		&newUser.APIKey,
	)
	if err != nil {
		log.Println(err)
		return newUser, err
	}
	return newUser, nil
}

func (m *postgresDBRepo) GetUserByAPIKey(apiKey string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, name, created_at, updated_at FROM users WHERE api_key = $1`

	var u models.User

	row := m.DB.QueryRowContext(ctx, query, apiKey)

	err := row.Scan(
		&u.ID,
		&u.Name,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		log.Println(err)
		return u, err
	}
	return u, nil
}
