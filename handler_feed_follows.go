package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wtran29/go-rss/internal/models"
)

func (repo *DBRepo) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user models.User) {
	type parameters struct {
		FeedID int `feed_id`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorJSON(w, fmt.Sprintf("Error parsing json: %v", err), 400)
		return
	}
	fmt.Println(params.FeedID)

	// apiKey := migrations.GenerateAPIKey()
	ff := models.FeedFollow{
		UserID:    user.ID,
		FeedID:    params.FeedID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	log.Println(ff.UserID)

	newFeedFollow, err := repo.DB.CreateFeedFollow(ff)
	if err != nil {
		log.Println(err)
		errorJSON(w, fmt.Sprintf("Could not create feed follow: %v", err), 400)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Created feed follow %v", newFeedFollow.ID),
		Data:    newFeedFollow,
	}

	writeJson(w, http.StatusCreated, payload)
}

func (repo *DBRepo) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user models.User) {

	feedFollows, err := repo.DB.GetFeedFollowsByUserID(user.ID)
	if err != nil {
		errorJSON(w, fmt.Sprintf("Could not retrieve feed follows: %v", err))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Data:    feedFollows,
		Message: "Retrieved all feed follows",
	}

	writeJson(w, http.StatusAccepted, payload)
}

func (repo *DBRepo) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user models.User) {
	id, err := strconv.Atoi(chi.URLParam(r, "feedFollowID"))
	if err != nil {
		errorJSON(w, fmt.Sprintf("Could not parse feed follow id: %v", err))
	}
	err = repo.DB.DeleteFeedFollow(id, user.ID)
	if err != nil {
		errorJSON(w, fmt.Sprintf("Could not delete feed follow: %v", err))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Unfollowed feed %v", id),
	}

	writeJson(w, http.StatusAccepted, payload)
}
