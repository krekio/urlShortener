package postgres

import (
	"database/sql"
	"github.com/krekio/urlShortener.git/internal/config"
	_ "github.com/lib/pq"
	"log"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(cfg config.Config) (*Storage, error) {

	db, err := sql.Open("postgres", cfg.StorageDsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS url (
            id SERIAL PRIMARY KEY,
            alias TEXT UNIQUE NOT NULL,
            url TEXT NOT NULL
        );
    `)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_alias ON url (alias);
    `)
	if err != nil {
		return nil, err
	}

	log.Println("Database connection established and tables are ready")
	return &Storage{DB: db}, nil
}
