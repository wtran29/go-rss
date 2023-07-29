package main

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wtran29/go-rss/internal/models"
)

func startScraping(repo *DBRepo, concurrency int, timeBetweenReq time.Duration) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenReq)
	ticker := time.NewTicker(timeBetweenReq)

	for ; ; <-ticker.C {
		feeds, err := repo.DB.GetNextFeedsToFetch(concurrency)
		if err != nil {
			log.Println("error fetching feeds:", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(repo, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(repo *DBRepo, wg *sync.WaitGroup, feed models.Feed) {
	defer wg.Done()

	_, err := repo.DB.MarkFeedsAsFetched(feed.ID)
	if err != nil {
		log.Println("error marking feed as fetched:", err)
		return
	}

	rssFeed, err := createFeedFromURL(feed.URL)
	if err != nil {
		log.Println("error fetching feed", err)
		return

	}

	for _, item := range rssFeed.Channel.Item {
		var publishedAt time.Time

		t, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Println("error parsing date:", err)
			return
		}
		publishedAt = t

		post := models.Post{
			ID:          uuid.New(),
			FeedID:      feed.ID,
			Title:       item.Title,
			Description: item.Description,
			URL:         item.Link,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			PublishedAt: publishedAt,
		}
		_, err = repo.DB.CreatePost(post)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("error creating post: %v", err)
		}

	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
