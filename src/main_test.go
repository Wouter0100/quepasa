package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestServer() *httptest.Server {
	router := newRouter()
	server := httptest.NewServer(router)

	return server
}

func TestReceiveAPI(t *testing.T) {
	t.Run("Successful receive, no timestamp", func(t *testing.T) {
		server := setupTestServer()
		defer server.Close()

		client := &http.Client{}
		res, err := client.Get(server.URL + "/v1/bot/zzz/receive")
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(res.StatusCode)
		fmt.Println(res)
	})
}
