package main

import (
	"fmt"
	"net/http"

	"github.com/wtran29/go-rss/internal/models"
)

func (repo *DBRepo) handlerCreatePost(w http.ResponseWriter, r *http.Request, user models.User) {

}

func (repo *DBRepo) handlerGetPostsByUser(w http.ResponseWriter, r *http.Request, user models.User) {

	posts, err := repo.DB.GetPostsByUser(user.ID, 10)
	if err != nil {
		errorJSON(w, fmt.Sprintf("Could not fetch posts by user %v", user.ID), http.StatusInternalServerError)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Data:    posts,
		Message: fmt.Sprintf("Retrieved posts from user %v", user.ID),
	}
	writeJson(w, http.StatusOK, payload)
}
