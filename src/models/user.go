package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string `db:"id"`
	Email     string `db:"email"`
	Username  string `db:"username"`
	Password  string `db:"password"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func CountUsers(db *sqlx.DB) (int, error) {
	var count int
	err := db.Get(&count, "SELECT count(*) FROM users")
	return count, err
}

func FindUserByID(db *sqlx.DB, ID string) (User, error) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE id = $1", ID)
	return user, err
}

func FindUserByEmail(db *sqlx.DB, email string) (User, error) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	return user, err
}

func CheckUserExists(db *sqlx.DB, email string) (bool, error) {
	var count int
	err := db.Get(&count, "SELECT count(*) FROM users WHERE email = $1", email)
	return count > 0, err
}

func CheckUser(db *sqlx.DB, email string, password string) (User, error) {
	var user User
	var out User
	err := db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return out, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return out, err
	}

	out = user
	return out, err
}

func CreateUser(db *sqlx.DB, email string, password string) (User, error) {
	var user User
	userID := uuid.New().String()
	now := time.Now().Format(time.RFC3339)

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return user, err
	}

	query := `INSERT INTO users
    (id, email, username, password, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := db.Exec(query, userID, email, email, string(hashed), now, now); err != nil {
		return user, err
	}

	return FindUserByID(db, userID)
}
