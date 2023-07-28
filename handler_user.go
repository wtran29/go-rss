package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/wtran29/go-rss/internal/driver"
	"github.com/wtran29/go-rss/internal/models"
	"github.com/wtran29/go-rss/internal/repository"
	"github.com/wtran29/go-rss/internal/repository/dbrepo"
)

var Repo *DBRepo
var app *ApiConfig

type DBRepo struct {
	App *ApiConfig
	DB  repository.DatabaseRepo
}

// NewHandlers creates the handlers
func NewHandlers(repo *DBRepo, a *ApiConfig) {
	Repo = repo
	app = a
}

// NewPostgresqlHandlers creates db repo for postgres
func NewPostgresqlHandlers(db *driver.DB, a *ApiConfig) *DBRepo {
	return &DBRepo{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL),
	}
}

func (repo *DBRepo) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `name`
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
	user := models.User{
		Name:      params.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	log.Println(user)

	newUser, err := repo.DB.CreateUser(user)
	if err != nil {
		log.Println(err)
		errorJSON(w, fmt.Sprintf("Could not create user: %v", err), 400)
		return
	}

	fmt.Println(newUser)

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Created user %s", newUser.Name),
		Data:    newUser,
	}

	writeJson(w, 201, payload)
}

func (repo *DBRepo) handlerGetUser(w http.ResponseWriter, r *http.Request, user models.User) {

	payload := jsonResponse{
		Error:   false,
		Data:    user,
		Message: fmt.Sprintf("Retrieved user %s", user.Name),
	}

	writeJson(w, 200, payload)
}
