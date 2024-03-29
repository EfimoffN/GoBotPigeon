package apitypes

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// UserRow ...
type UserRow struct {
	UserID   string `db:"userid"`
	NameUser string `db:"nameuser"`
	ChatID   string `db:"chatid"`
}

// CodeRow ...
type CodeRow struct {
	CodeID string `db:"codeid"`
	Code   string `db:"code"`
}

// RefUserCode ...
type RefUserCode struct {
	KeyID  string `db:"keyid"`
	CodeID string `db:"codeid"`
	UserID string `db:"userid"`
}

// LetterRow ...
type LetterRow struct {
	LetterID   string    `db:"letterid"`
	CodeID     string    `db:"codeid"`
	UserID     string    `db:"userid"`
	Letter     string    `db:"letter"`
	DataLetter time.Time `db:"dataletter"`
}

// StoreDB ...
type StoreDB struct {
	DB *sqlx.DB
}

// BotWork ...
type BotWork struct {
	BotWorkID   string `db:"botworkid"`
	UserID      string `db:"userid"`
	BotWorkFlag bool   `db:"botworkflag"`
}

// LastUserCommand ...
type LastUserCommand struct {
	CommandID   string    `db:"commandid"`
	UserID      string    `db:"userid"`
	Command     string    `db:"command"`
	DataCommand time.Time `db:"datacommand"`
}
