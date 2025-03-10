package postgres

import (
	"database/sql"
	"fmt"
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

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	const op = "storage.postgres.SaveURL"
	stmt, err := s.DB.Prepare("INSERT INTO url (url, alias) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("%s %s %s", op, err, stmt)
	}

	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		return fmt.Errorf("%s %s %s", op, err, stmt)
	}
	fmt.Print("URL SAVED")
	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	stmt, err := s.DB.Prepare("SELECT url FROM url WHERE alias = $1")
	if err != nil {
		return "", err
	}
	var resUrl string
	err = stmt.QueryRow(alias).Scan(&resUrl)
	if err != nil {
		return "", err
	}
	return resUrl, nil

}
func (s *Storage) DeleteURL(alias string) error {
	stmt, err := s.DB.Prepare("DELETE FROM url WHERE alias = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(alias)
	if err != nil {
		return err
	}
	fmt.Print("URL DELETED")
	return nil

}
