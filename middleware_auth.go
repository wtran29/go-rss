package main

import (
	"fmt"
	"net/http"

	"github.com/wtran29/go-rss/internal/auth"
	"github.com/wtran29/go-rss/internal/models"
)

type authHandler func(http.ResponseWriter, *http.Request, models.User)

func (repo *DBRepo) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			errorJSON(w, fmt.Sprintf("Auth error: %v", err), 403)
			return
		}

		user, err := repo.DB.GetUserByAPIKey(apiKey)
		if err != nil {
			errorJSON(w, fmt.Sprintf("Could not get user: %v", err), 400)
			return
		}
		handler(w, r, user)
	}
}
