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

func (m *postgresDBRepo) GetNextFeedsToFetch(limit int) ([]models.Feed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, name, url, user_id, created_at, updated_at, last_fetched_at FROM feeds
				ORDER BY last_fetched_at ASC NULLS FIRST LIMIT $1`

	var feeds []models.Feed

	rows, err := m.DB.QueryContext(ctx, query, limit)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var feed models.Feed
		err := rows.Scan(
			&feed.ID,
			&feed.Name,
			&feed.URL,
			&feed.UserID,
			&feed.CreatedAt,
			&feed.UpdatedAt,
			&feed.LastFetchedAt,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		feeds = append(feeds, feed)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}

func (m *postgresDBRepo) MarkFeedsAsFetched(id int) (models.Feed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `UPDATE feeds SET last_fetched_at = $1, updated_at = $2 WHERE id = $3`

	_, err := m.DB.ExecContext(ctx, stmt,
		time.Now(),
		time.Now(),
		id,
	)
	if err != nil {
		log.Println(err)
		return models.Feed{}, err
	}

	stmt = `SELECT * FROM feeds WHERE id = $1`

	var feed models.Feed

	err = m.DB.QueryRowContext(ctx, stmt, id).Scan(
		&feed.ID,
		&feed.Name,
		&feed.URL,
		&feed.UserID,
		&feed.CreatedAt,
		&feed.UpdatedAt,
		&feed.LastFetchedAt,
	)
	if err != nil {
		log.Println(err)
		return feed, err
	}

	return feed, nil
}
