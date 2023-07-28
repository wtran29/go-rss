package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/wtran29/go-rss/internal/models"
)

func (repo *DBRepo) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user models.User) {
	type parameters struct {
		Name string `name`
		URL  string `url`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorJSON(w, fmt.Sprintf("Error parsing json: %v", err), 400)
		return
	}
	fmt.Println(params.Name)

	// apiKey := migrations.GenerateAPIKey()
	feed := models.Feed{
		Name:      params.Name,
		URL:       params.URL,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newFeed, err := repo.DB.CreateFeed(feed)
	if err != nil {
		log.Println(err)
		errorJSON(w, fmt.Sprintf("Could not create feed: %v", err), 400)
		return
	}

	fmt.Println(newFeed)

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Created feed %s", newFeed.Name),
		Data:    newFeed,
	}

	writeJson(w, 201, payload)
}

func (repo *DBRepo) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {

	feeds, err := repo.DB.GetFeeds()
	if err != nil {
		errorJSON(w, fmt.Sprintf("Could not retrieve feeds: %v", err))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Data:    feeds,
		Message: "Retrieved all feeds",
	}

	writeJson(w, 200, payload)
}
