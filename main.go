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
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

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
