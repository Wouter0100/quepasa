package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Bot struct {
	ID        string `db:"id" json:"id"`
	Number    string `db:"number" json:"number"`
	Verified  bool   `db:"is_verified" json:"is_verified"`
	Token     string `db:"token" json:"token"`
	UserID    string `db:"user_id" json:"user_id"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

func FindAllBots(db *sqlx.DB) ([]Bot, error) {
	bots := []Bot{}
	err := db.Select(&bots, "SELECT * FROM bots")
	return bots, err
}

func FindAllBotsForUser(db *sqlx.DB, userID string) ([]Bot, error) {
	bots := []Bot{}
	err := db.Select(&bots, "SELECT * FROM bots WHERE user_id = $1", userID)
	return bots, err
}

func FindBotByToken(db *sqlx.DB, token string) (Bot, error) {
	var bot Bot
	err := db.Get(&bot, "SELECT * FROM bots WHERE token = $1", token)
	return bot, err
}

func FindBotForUser(db *sqlx.DB, userID string, ID string) (Bot, error) {
	var bot Bot
	err := db.Get(&bot, "SELECT * FROM bots WHERE user_id = $1 AND id = $2", userID, ID)
	return bot, err
}

func FindBotByNumber(db *sqlx.DB, number string) (Bot, error) {
	var bot Bot
	err := db.Get(&bot, "SELECT * FROM bots WHERE number = $1", number)
	return bot, err
}

func CreateBot(db *sqlx.DB, userID string, number string) (Bot, error) {
	var bot Bot
	botID := uuid.New().String()
	token := uuid.New().String()
	now := time.Now().Format(time.RFC3339)
	query := `INSERT INTO bots
    (id, number, is_verified, token, user_id, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7)`
	if _, err := db.Exec(query, botID, number, false, token, userID, now, now); err != nil {
		return bot, err
	}

	return FindBotForUser(db, userID, botID)
}

func (bot *Bot) MarkVerified(db *sqlx.DB) error {
	now := time.Now().Format(time.RFC3339)
	query := "UPDATE bots SET is_verified = true, updated_at = $1 WHERE id = $2"
	_, err := db.Exec(query, now, bot.ID)
	return err
}

func (bot *Bot) CycleToken(db *sqlx.DB) error {
	token := uuid.New().String()
	now := time.Now().Format(time.RFC3339)
	query := "UPDATE bots SET token = $1, updated_at = $2 WHERE id = $3"
	_, err := db.Exec(query, token, now, bot.ID)
	return err
}

func (bot *Bot) Delete(db *sqlx.DB) error {
	query := "DELETE FROM bots WHERE id = $1"
	_, err := db.Exec(query, bot.ID)
	return err
}

func (bot *Bot) FormattedNumber() string {
	return CleanPhoneNumber(bot.Number)
}
