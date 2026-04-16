package main 

import (
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

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerResetNumber)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(serv.ListenAndServe())
}
