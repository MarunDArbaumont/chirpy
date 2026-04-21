package main 

import _ "github.com/lib/pq"

import (
	"log"
	"net/http"
	"sync/atomic"
	"os"
	"database/sql"

	"github.com/joho/godotenv"
	"github.com/MarunDArbaumont/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	database *database.Queries
	platform string
}

func main () {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error while opening bd: %v", err)
	}


	dbQueries := database.New(db)
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		database: dbQueries,
		platform: platform,
	}

	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	serv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/chirps", cfg.handlerChirps)
	mux.HandleFunc("GET /api/chirps", cfg.handlerAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerSingleChirp)
	mux.HandleFunc("POST /api/users", cfg.handlerUsers)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerResetUsers)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(serv.ListenAndServe())
}
