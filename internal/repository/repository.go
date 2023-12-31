package repository

import "github.com/wtran29/go-rss/internal/models"

type DatabaseRepo interface {
	CreateUser(u models.User) (models.User, error)
	GetUserByAPIKey(apiKey string) (models.User, error)
	CreateFeed(f models.Feed) (models.Feed, error)
	GetFeeds() ([]models.Feed, error)
	CreateFeedFollow(ff models.FeedFollow) (models.FeedFollow, error)
	GetFeedFollowsByUserID(uid int) ([]models.FeedFollow, error)
	DeleteFeedFollow(id, uid int) error
	GetNextFeedsToFetch(limit int) ([]models.Feed, error)
	MarkFeedsAsFetched(id int) (models.Feed, error)
	CreatePost(p models.Post) (models.Post, error)
	GetPostsByUser(uid, lim int) ([]models.Post, error)
}
