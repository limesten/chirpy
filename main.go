package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type apiConfig struct {
	fileserverHits int
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func main() {
	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	router := chi.NewRouter()

	fileServer := apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))
	router.Mount("/", fileServer)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", readinessHandler)
	apiRouter.Post("/validate_chirp", chirpValidationHandler)
	router.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.metricsHandler)
	router.Mount("/admin", adminRouter)

	corsMux := middlewareCors(router)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: corsMux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
