package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/wtran29/go-rss/internal/driver"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

}

type ApiConfig struct {
	DB *driver.DB
}

func run() error {
	log := log.New(os.Stdout, "GO-RSS : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	var cfg struct {
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:8080"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
	}

	if err := conf.Parse(os.Args[1:], "GO-RSS", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("GO-RSS", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
	}

	log.Printf("main : Started : Application initializing")
	defer log.Println("main : Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main : Config :\n%v\n", out)

	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("port is not found in env")
	}

	// dbType := os.Getenv("DATABASE_TYPE")
	dbHost := os.Getenv("DATABASE_HOST")
	dbUser := os.Getenv("DATABASE_USER")
	dbPort := os.Getenv("DATABASE_PORT")
	dbPw := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")
	dbSSL := os.Getenv("DATABASE_SSL_MODE")

	dsnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
		dbHost, dbPort, dbUser, dbPw, dbName, dbSSL,
	)
	if dsnString == "" {
		log.Fatal("invalid DSN string")
	}

	db, err := driver.ConnectPostgres(dsnString)
	if err != nil {
		log.Fatal("Cannot connect to database!", err)
	}

	apiCfg := ApiConfig{
		DB: db,
	}

	app = &apiCfg

	repo := NewPostgresqlHandlers(db, app)
	NewHandlers(repo, app)

	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1 := chi.NewRouter()
	v1.Get("/healthz", handlerReadiness)
	v1.Get("/err", handlerErr)
	v1.Post("/users", repo.handlerCreateUser)
	v1.Get("/users", repo.handlerGetUser)

	mux.Mount("/v1", v1)

	srv := &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf(":%s", port),
	}

	log.Printf("Starting server on port: %v", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Port:", port)

	return nil
}
