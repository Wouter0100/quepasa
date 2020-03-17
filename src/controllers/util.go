package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"regexp"
)

func parseJSONBody(r *http.Request) (map[string]interface{}, error) {
	var postParams map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&postParams); err != nil {
		return postParams, err
	}

	return postParams, nil
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}

func webSocketProtocol() string {
	protocol := "wss"
	if os.Getenv("APP_ENV") == "development" {
		protocol = "ws"
	}
	return protocol
}

func validateEmail(s string) bool {
	var rx = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(s) < 255 && rx.MatchString(s) {
		return true
	}

	return false
}
