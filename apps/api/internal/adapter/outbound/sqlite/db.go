// Package sqlite est un adaptateur outbound : persistance via SQLite
// (driver pur-Go modernc.org/sqlite, sans CGO).
package sqlite

import (
	"database/sql"

	// Driver SQLite pur-Go (enregistré sous le nom "sqlite").
	_ "modernc.org/sqlite"
)

// Open ouvre la base SQLite au chemin donné et applique les migrations.
// path peut être un fichier (ex. /data/graefik.db) ou ":memory:" (tests).
func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	// SQLite gère un seul écrivain : on limite à une connexion pour éviter
	// les erreurs "database is locked" (suffisant pour une app mono-instance).
	db.SetMaxOpenConns(1)

	if err := migrate(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	const schema = `
CREATE TABLE IF NOT EXISTS users (
	id            TEXT PRIMARY KEY,
	username      TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	created_at    INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS sessions (
	token      TEXT PRIMARY KEY,
	user_id    TEXT NOT NULL,
	created_at INTEGER NOT NULL,
	expires_at INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
`
	_, err := db.Exec(schema)
	return err
}
