package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/emilmalmsten/chirpy/internal/jsonDB"
	"github.com/go-chi/chi"
)

func postChirp(db *jsonDB.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Body string `json:"body"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
			return
		}

		const maxChirpLength = 140

		if len(params.Body) > maxChirpLength {
			respondWithError(w, http.StatusBadRequest, "Chirp is too long")
			return
		}

		cleanChirp := filterProfanity(params.Body)

		chirp, err := db.CreateChirp(cleanChirp)
		if err != nil {
			fmt.Printf("err with create chirp: %s", err)
			respondWithError(w, http.StatusInternalServerError, "Failed to create Chirp")
			return
		}

		respondWithJSON(w, http.StatusCreated, chirp)
	}
}

func getAllChirps(db *jsonDB.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		chirps, err := db.GetChirps()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch chirps")
			return
		}

		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Id < chirps[j].Id
		})

		respondWithJSON(w, http.StatusOK, chirps)

	}
}

func getChirp(db *jsonDB.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		chirpIdString := chi.URLParam(r, "chirpID")
		chirpID, err := strconv.Atoi(chirpIdString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
			return
		}

		chirp, err := db.GetChirp(chirpID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}

		respondWithJSON(w, http.StatusOK, chirp)

	}
}
