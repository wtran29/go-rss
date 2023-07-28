package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/wtran29/go-rss/internal/models"
)

func (m *postgresDBRepo) CreateFeed(f models.Feed) (models.Feed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO feeds (created_at, updated_at, name, url, user_id)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id`

	var feedID int
	err := m.DB.QueryRowContext(ctx, stmt,
		time.Now(),
		time.Now(),
		f.Name,
		f.URL,
		f.UserID,
	).Scan(&feedID)
	if err != nil {
		log.Println(err)
	}

	query := `SELECT id, name, url, user_id, created_at, updated_at FROM feeds WHERE id = $1`

	var newFeed models.Feed

	row := m.DB.QueryRowContext(ctx, query, feedID)

	err = row.Scan(
		&newFeed.ID,
		&newFeed.Name,
		&newFeed.URL,
		&newFeed.UserID,
		&newFeed.CreatedAt,
		&newFeed.UpdatedAt,
	)
	if err != nil {
		log.Println(err)
		return newFeed, err
	}
	return newFeed, nil
}

func (m *postgresDBRepo) GetFeeds() ([]models.Feed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, name, url, user_id, created_at, updated_at FROM feeds`

	var feeds []models.Feed

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return feeds, err
	}

	defer rows.Close()

	for rows.Next() {
		var f models.Feed
		err := rows.Scan(
			&f.ID,
			&f.Name,
			&f.URL,
			&f.UserID,
			&f.CreatedAt,
			&f.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
			return feeds, err
		}
		feeds = append(feeds, f)
	}
	return feeds, nil
}
