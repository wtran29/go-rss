package dbrepo

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/wtran29/go-rss/internal/models"
)

func (m *postgresDBRepo) CreatePost(p models.Post) (models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO posts (id, feed_id, title, description, url, published_at, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
				RETURNING id, feed_id, title, description, url, published_at, created_at, updated_at`

	// if description is blank, set to null in db
	desc := sql.NullString{}
	if p.Description != "" {
		desc.String = p.Description
		desc.Valid = true
	}

	var post models.Post
	err := m.DB.QueryRowContext(ctx, stmt,
		p.ID,
		p.FeedID,
		p.Title,
		desc,
		p.URL,
		p.PublishedAt,
		p.CreatedAt,
		p.UpdatedAt,
	).Scan(
		&post.ID,
		&post.FeedID,
		&post.Title,
		&post.Description,
		&post.URL,
		&post.PublishedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {

		return post, err
	}
	return post, nil
}

func (m *postgresDBRepo) GetPostsByUser(uid, lim int) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT p.id, p.title, p.description, p.published_at, p.url, p.feed_id, p.created_at, p.updated_at
				FROM posts p
				LEFT JOIN feed_follows ff ON (p.feed_id = ff.feed_id)
				WHERE ff.user_id = $1
				ORDER BY p.published_at DESC
				LIMIT $2`

	rows, err := m.DB.QueryContext(ctx, query, uid, lim)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Description,
			&post.PublishedAt,
			&post.URL,
			&post.FeedID,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
