package controllers

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/nbutton23/zxcvbn-go"
	"gitlab.com/digiresilience/link/quepasa/models"
)

//
// Account
//

type accountFormData struct {
	PageTitle    string
	ErrorMessage string
	Bots         []models.Bot
	User         models.User
}

// AccountFormHandler renders route GET "/account"
func AccountFormHandler(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetUser(r)
	if err != nil {
		redirectToLogin(w, r)
	}

	data := accountFormData{
		PageTitle: "Account",
		User:      user,
	}

	bots, err := models.FindAllBotsForUser(models.GetDB(), user.ID)
	if err != nil {
		data.ErrorMessage = err.Error()
	} else {
		data.Bots = bots
	}

	templates := template.Must(template.ParseFiles("views/layouts/main.tmpl", "views/account.tmpl"))
	templates.ExecuteTemplate(w, "main", data)
}

//
// Login
//

type loginFormData struct {
	PageTitle string
}

// LoginFormHandler renders route GET "/login"
func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	data := loginFormData{
		PageTitle: "Login",
	}

	templates := template.Must(template.ParseFiles("views/layouts/main.tmpl", "views/login.tmpl"))
	templates.ExecuteTemplate(w, "main", data)
}

// LoginHandler renders route POST "/login"
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	if email == "" || password == "" {
		respondUnauthorized(w, errors.New("Missing username or password"))
		return
	}

	user, err := models.CheckUser(models.GetDB(), email, password)
	if err != nil {
		respondUnauthorized(w, errors.New("Incorrect username or password"))
		return
	}

	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("SIGNING_SECRET")), nil)
	claims := jwt.MapClaims{"user_id": user.ID}
	jwtauth.SetIssuedNow(claims)
	jwtauth.SetExpiryIn(claims, 24*time.Hour)
	_, tokenString, _ := tokenAuth.Encode(claims)
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		MaxAge:   60 * 60 * 24,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/account", http.StatusFound)
}

// LogoutHandler renders route GET "/logoout"
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		MaxAge:   0,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	redirectToLogin(w, r)
}

//
// Setup
//

type setupFormData struct {
	PageTitle             string
	ErrorMessage          string
	Email                 string
	EmailError            bool
	UserExistsError       bool
	EmailInvalidError     bool
	PasswordMatchError    bool
	PasswordStrengthError bool
	PasswordCrackTime     string
}

func renderSetupForm(w http.ResponseWriter, data setupFormData) {
	templates := template.Must(template.ParseFiles("views/layouts/main.tmpl", "views/setup.tmpl"))
	templates.ExecuteTemplate(w, "main", data)
}

// SetupFormHandler renders route GET "/setup"
func SetupFormHandler(w http.ResponseWriter, r *http.Request) {
	count, err := models.CountUsers(models.GetDB())
	if count > 0 || err != nil {
		redirectToLogin(w, r)
		return
	}

	data := setupFormData{
		PageTitle: "Setup",
	}

	renderSetupForm(w, data)
}

// SetupHandler renders route POST "/setup"
func SetupHandler(w http.ResponseWriter, r *http.Request) {
	count, err := models.CountUsers(models.GetDB())
	if count > 0 || err != nil {
		redirectToLogin(w, r)
		return
	}

	data := setupFormData{
		PageTitle: "Setup",
	}

	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	passwordConfirm := r.Form.Get("passwordConfirm")

	if email == "" || password == "" {
		data.ErrorMessage = "Email and password are required"
		data.EmailError = true
		renderSetupForm(w, data)
		return
	}

	data.Email = email

	if !validateEmail(email) {
		data.ErrorMessage = "Email is invalid"
		data.EmailInvalidError = true
		renderSetupForm(w, data)
		return
	}

	if password != passwordConfirm {
		data.ErrorMessage = "Passwords don't match"
		data.PasswordMatchError = true
		renderSetupForm(w, data)
		return
	}

	res := zxcvbn.PasswordStrength(password, nil)
	if res.Score < 1 {
		data.ErrorMessage = "Password is too weak"
		data.PasswordStrengthError = true
		data.PasswordCrackTime = res.CrackTimeDisplay
		renderSetupForm(w, data)
		return
	}

	exists, err := models.CheckUserExists(models.GetDB(), email)
	if err != nil {
		data.ErrorMessage = err.Error()
		renderSetupForm(w, data)
		return
	}

	if exists {
		data.UserExistsError = true
		renderSetupForm(w, data)
		return
	}

	_, err = models.CreateUser(models.GetDB(), email, password)
	if err != nil {
		data.ErrorMessage = err.Error()
		renderSetupForm(w, data)
		return
	}

	redirectToLogin(w, r)
}
