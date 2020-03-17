package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Result string `json:"result"`
}

func respondSuccess(w http.ResponseWriter, res interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func respondBadRequest(w http.ResponseWriter, err error) {
	log.Println("Bad request: ", err)

	respondError(w, err, http.StatusBadRequest)
}

func respondUnauthorized(w http.ResponseWriter, err error) {
	log.Println("Unauthorized: ", err)

	respondError(w, err, http.StatusUnauthorized)
}

func respondNotFound(w http.ResponseWriter, err error) {
	log.Println("Not found: ", err)

	respondError(w, err, http.StatusNotFound)
}

func respondServerError(w http.ResponseWriter, err error) {
	log.Println("Server error: ", err)

	respondError(w, err, http.StatusInternalServerError)
}

func respondError(w http.ResponseWriter, err error, code int) {
	res := &errorResponse{
		Result: err.Error(),
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
