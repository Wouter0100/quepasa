package controllers

import (
	"gitlab.com/digiresilience/link/quepasa/models"
	"html/template"
	"net/http"
)

type indexData struct {
	PageTitle string
}

// IndexHandler renders route GET "/"
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	_, err := models.GetUser(r)
	if err != nil {
		redirectToLogin(w, r)
		return
	}

	data := indexData{
		PageTitle: "Home",
	}

	templates := template.Must(template.ParseFiles("views/layouts/main.tmpl", "views/index.tmpl"))
	templates.ExecuteTemplate(w, "main", data)
}
