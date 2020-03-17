package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/digiresilience/link/quepasa/controllers"
	"gitlab.com/digiresilience/link/quepasa/models"
)

func main() {
	err := models.MigrateToLatest()
	if err != nil {
		log.Fatalf("Database migration error: %s", err.Error())
	}

	err = models.StartServer()
	if err != nil {
		log.Printf("Failed to start WhatsApp server: %s", err.Error())
	}

	go func() {
		log.Println("Starting metrics service")

		m := chi.NewRouter()
		m.Handle("/metrics", promhttp.Handler())
		host := fmt.Sprintf("%s:%s", os.Getenv("METRICS_HOST"), os.Getenv("METRICS_PORT"))
		http.ListenAndServe(host, m)
	}()

	r := newRouter()

	log.Println("Starting web server")
	http.ListenAndServe(":3000", r)
}

func newRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// web routes
	addWebRoutes(r)

	// api routes
	addAPIRoutes(r)

	// static files
	workDir, _ := os.Getwd()
	assetsDir := filepath.Join(workDir, "assets")
	fileServer(r, "/assets", http.Dir(assetsDir))

	return r
}

func addWebRoutes(r chi.Router) {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("SIGNING_SECRET")), nil)

	// authenticated web routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(authenticator)

		r.Get("/account", controllers.AccountFormHandler)
		r.Get("/bot/register", controllers.RegisterFormHandler)
		r.Post("/bot/register", controllers.RegisterHandler)
		r.Get("/bot/{botID}/verify/ws", controllers.VerifyHandler)
		r.Get("/bot/{botID}/verify", controllers.VerifyFormHandler)
		r.Post("/bot/delete", controllers.DeleteHandler)
		r.Post("/bot/cycle", controllers.CycleHandler)
		r.Get("/bot/{botID}", controllers.SendFormHandler)
		r.Get("/bot/{botID}/send", controllers.SendFormHandler)
		r.Post("/bot/{botID}/send", controllers.SendHandler)
		r.Get("/bot/{botID}/receive", controllers.ReceiveFormHandler)
	})

	// unauthenticated web routes
	r.Group(func(r chi.Router) {
		r.Get("/", controllers.IndexHandler)
		r.Get("/login", controllers.LoginFormHandler)
		r.Post("/login", controllers.LoginHandler)
		r.Get("/setup", controllers.SetupFormHandler)
		r.Post("/setup", controllers.SetupHandler)
		r.Get("/logout", controllers.LogoutHandler)
	})
}

func addAPIRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Get("/v1/bot/{token}", controllers.InfoAPIHandler)
		r.Post("/v1/bot/{token}/send", controllers.SendAPIHandler)
		r.Get("/v1/bot/{token}/receive", controllers.ReceiveAPIHandler)
	})
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"
	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if token == nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
