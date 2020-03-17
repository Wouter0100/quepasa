package models

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store struct {
	BotID     string `db:"bot_id"`
	Data      []byte `db:"data"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func storeExists(db *sqlx.DB, botID string) (bool, error) {
	var count int
	err := db.Get(&count, "SELECT count(*) FROM sessionstore WHERE bot_id = $1", botID)
	return count > 0, err
}

func GetOrCreateStore(db *sqlx.DB, botID string) (Store, error) {
	var store Store
	exists, err := storeExists(db, botID)
	if err != nil {
		return store, err
	}

	if exists {
		store, err := GetStore(db, botID)
		if err != nil {
			return store, err
		}
	} else {
		store, err = CreateStore(db, botID)
	}
	return store, err
}

func CreateStore(db *sqlx.DB, botID string) (Store, error) {
	var user Store
	now := time.Now().Format(time.RFC3339)
	query := `INSERT INTO sessionstore
    (bot_id, created_at, updated_at)
    VALUES ($1, $2, $3)`
	if _, err := db.Exec(query, botID, now, now); err != nil {
		return user, err
	}

	return GetStore(db, botID)
}

func GetStore(db *sqlx.DB, botID string) (Store, error) {
	var store Store
	err := db.Get(&store, "SELECT * FROM sessionstore WHERE bot_id = $1", botID)
	return store, err
}

func UpdateStore(db *sqlx.DB, botID string, data []byte) ([]byte, error) {
	now := time.Now().Format(time.RFC3339)
	query := "UPDATE sessionstore SET data = ($1::bytea), updated_at = $2 WHERE bot_id = $3"
	_, err := db.Exec(query, data, now, botID)
	return data, err
}

func DeleteStore(db *sqlx.DB, botID string) error {
	query := "DELETE FROM sessionstore WHERE bot_id = $1"
	_, err := db.Exec(query, botID)
	return err
}
