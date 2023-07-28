package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/wtran29/go-rss/internal/models"
)

func (m *postgresDBRepo) CreateFeedFollow(ff models.FeedFollow) (models.FeedFollow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO feed_follows (created_at, updated_at, user_id, feed_id)
				VALUES ($1, $2, $3, $4)
				RETURNING id`

	var ffID int
	err := m.DB.QueryRowContext(ctx, stmt,
		time.Now(),
		time.Now(),
		ff.UserID,
		ff.FeedID,
	).Scan(&ffID)
	if err != nil {
		log.Println(err)
		return ff, err
	}

	log.Println("id:", ffID)

	query := `SELECT id, user_id, feed_id, created_at, updated_at FROM feed_follows WHERE id = $1`

	var newFeedFollows models.FeedFollow

	row := m.DB.QueryRowContext(ctx, query, ffID)

	err = row.Scan(
		&newFeedFollows.ID,
		&newFeedFollows.UserID,
		&newFeedFollows.FeedID,
		&newFeedFollows.CreatedAt,
		&newFeedFollows.UpdatedAt,
	)
	if err != nil {
		log.Println(err)
		return newFeedFollows, err
	}
	return newFeedFollows, nil
}

func (m *postgresDBRepo) GetFeedFollowsByUserID(uid int) ([]models.FeedFollow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, feed_id, created_at, updated_at FROM feed_follows WHERE user_id = $1`

	var feedFollows []models.FeedFollow

	rows, err := m.DB.QueryContext(ctx, query, uid)
	if err != nil {
		return feedFollows, err
	}

	defer rows.Close()

	for rows.Next() {
		var f models.FeedFollow
		err := rows.Scan(
			&f.ID,
			&f.UserID,
			&f.FeedID,
			&f.CreatedAt,
			&f.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
			return feedFollows, err
		}
		feedFollows = append(feedFollows, f)
	}
	return feedFollows, nil
}

func (m *postgresDBRepo) DeleteFeedFollow(id, uid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM feed_follows WHERE id = $1 AND user_id = $2`

	_, err := m.DB.ExecContext(ctx, stmt, id, uid)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
