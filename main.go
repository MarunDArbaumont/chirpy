package main 

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main () {
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	serv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /healthz", handlerReadiness)

	mux.HandleFunc("GET /metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /reset", cfg.handlerResetNumber)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(serv.ListenAndServe())
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}